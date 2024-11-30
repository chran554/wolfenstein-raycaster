package wolf3d

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
)

// https://archive.org/details/wolf3dsw
// https://moddingwiki.shikadi.net/wiki/GameMaps_Format

//go:embed resources/MAPHEAD.WL1
var wolfenstein3DMapHeaderData []byte

//go:embed resources/GAMEMAPS.WL1
var wolfenstein3DMapData []byte

type mapHeader struct {
	magic       uint16
	levelOffset []int32
}

type levelHeader struct {
	offPlane0 int32  // Offset 0 - Offset in GAMEMAPS to the beginning of compressed plane 0 data (or <= 0 if plane is not present)
	offPlane1 int32  // Offset 4 - Offset in GAMEMAPS to the beginning of compressed plane 1 data (or <= 0 if plane is not present)
	offPlane2 int32  // Offset 8 - Offset in GAMEMAPS to the beginning of compressed plane 2 data (or <= 0 if plane is not present)
	lenPlane0 uint16 // Offset 12 - Length of compressed plane 0 data (in bytes)
	lenPlane1 uint16 // Offset 14 - Length of compressed plane 1 data (in bytes)
	lenPlane2 uint16 // Offset 16 - Length of compressed plane 2 data (in bytes)
	width     uint16 // Offset 18 - Width of level (in tiles)
	height    uint16 // Offset 20 - Height of level (in tiles)
	name      string // Offset 22 - Internal name for level (used only by editor, not displayed in-game. 16 byte array, null-terminated).
}

type LevelMap struct {
	Name   string
	Width  int
	Height int
	plane0 []uint16 // Wolfenstein 3D: Walls and doors
	plane1 []uint16 // Wolfenstein 3D: Enemy initial positions and objects
	plane2 []uint16 // Wolfenstein 3D: empty
}

// Value gives the map value for coordinate x and y.
// Valid values are 0 <= x < width, and 0 <= y < height.
// -1 is returned if x or y is out of bounds.
func (lm *LevelMap) Value(plane, x, y int) int {
	if x < 0 || y < 0 || x >= lm.Width || y >= lm.Height {
		return -1
	}

	mapPlaneData := lm.plane(plane)
	if mapPlaneData == nil {
		return -1
	} else {
		return int(mapPlaneData[y*lm.Width+x])
	}
}

func (lm *LevelMap) plane(plane int) []uint16 {
	if plane == 0 {
		return lm.plane0
	} else if plane == 1 {
		return lm.plane1
	} else if plane == 2 {
		return lm.plane2
	}

	return nil
}

func (lm *LevelMap) planeExist(plane int) bool {
	if plane == 0 {
		return lm.plane0 != nil
	} else if plane == 1 {
		return lm.plane1 != nil
	} else if plane == 2 {
		return lm.plane2 != nil
	}

	return false
}

func (h *mapHeader) LevelCount() int {
	var levelCount int
	for levelCount = 0; levelCount < len(h.levelOffset) && (h.levelOffset[levelCount] > 0); levelCount++ {
	}
	return levelCount
}

func (h *mapHeader) LevelOffset(levelIndex int) int32 {
	if levelIndex >= h.LevelCount() {
		return -1
	}

	return h.levelOffset[levelIndex]
}

func Wolfenstein3DMap() ([]LevelMap, error) {
	var levelMaps []LevelMap

	mh, err := readMapHeader(wolfenstein3DMapHeaderData)
	if err != nil {
		return nil, err
	}

	lhs, err := readLevelHeaders(mh, wolfenstein3DMapData)
	if err != nil {
		return nil, err
	}

	for levelIndex, lh := range lhs {
		fmt.Printf("Reading level %d: %s\n", levelIndex, lh.name)
		levelMap := LevelMap{Name: lh.name, Width: int(lh.width), Height: int(lh.height)}

		expectedMapByteSize := lh.width * lh.height * 2

		plane0Exist := lh.offPlane0 > 0
		plane1Exist := lh.offPlane1 > 0
		plane2Exist := lh.offPlane2 > 0

		if plane0Exist {
			plane0MapData, err := readPlaneData(lh.offPlane0, lh.offPlane0+int32(lh.lenPlane0), expectedMapByteSize, mh.magic)
			if err != nil {
				return levelMaps, err
			}
			levelMap.plane0 = plane0MapData
		}

		if plane1Exist {
			plane1MapData, err := readPlaneData(lh.offPlane1, lh.offPlane1+int32(lh.lenPlane1), expectedMapByteSize, mh.magic)
			if err != nil {
				return levelMaps, err
			}
			levelMap.plane1 = plane1MapData
		}

		if plane2Exist {
			plane2MapData, err := readPlaneData(lh.offPlane2, lh.offPlane2+int32(lh.lenPlane2), expectedMapByteSize, mh.magic)
			if err != nil {
				return levelMaps, err
			}
			levelMap.plane2 = plane2MapData
		}

		levelMaps = append(levelMaps, levelMap)
	}

	return levelMaps, nil
}

func readPlaneData(startOffset int32, endOffset int32, expectedMapByteSize uint16, rleFlag uint16) ([]uint16, error) {
	var levelMapData []uint16
	compressedLevelData := wolfenstein3DMapData[startOffset:endOffset]

	rlew, carmackAndRLEW, err := checkCompressionMethods(compressedLevelData, expectedMapByteSize)
	if err != nil {
		return nil, err
	}

	if rlew {
		_, levelData := RLEWDecodeWithLengthPrefixAndRLEFlag(compressedLevelData, rleFlag)
		levelMapData = toUint16(levelData)
	} else if carmackAndRLEW {
		_, compressedLevelData = CarmackDecodeWithLengthPrefix(compressedLevelData)
		_, levelData := RLEWDecodeWithLengthPrefixAndRLEFlag(compressedLevelData, rleFlag)
		levelMapData = toUint16(levelData)
	}

	return levelMapData, nil
}

func toUint16(data []byte) []uint16 {
	uints := make([]uint16, len(data)/2)

	var value uint16
	for i := 0; i < len(uints); i++ {
		_ = binary.Read(bytes.NewBuffer(data[i*2:i*2+2]), binary.LittleEndian, &value)
		uints[i] = value
	}

	return uints
}

func checkCompressionMethods(compressedLevelData []byte, expectedMapByteSize uint16) (compressionRLEW bool, compressionCarmackAndRLEW bool, err error) {
	var mapSize1 uint16
	var mapSize2 uint16

	mapSizeReader := bytes.NewBuffer(compressedLevelData)
	if err = binary.Read(mapSizeReader, binary.LittleEndian, &mapSize1); err != nil {
		return false, false, err
	}
	if err = binary.Read(mapSizeReader, binary.LittleEndian, &mapSize2); err != nil {
		return false, false, err
	}

	compressionRLEW = mapSize1 == expectedMapByteSize
	compressionCarmackAndRLEW = mapSize2 == expectedMapByteSize

	return compressionRLEW, compressionCarmackAndRLEW, nil
}

func readMapHeader(mapHeaderData []byte) (*mapHeader, error) {
	mapHeaderDataBuffer := bytes.NewBuffer(mapHeaderData)

	var magic uint16
	err := binary.Read(mapHeaderDataBuffer, binary.LittleEndian, &magic)
	if err != nil {
		return nil, err
	}

	levelPtr := make([]int32, 100)
	for i := 0; i < 100; i++ {
		err = binary.Read(mapHeaderDataBuffer, binary.LittleEndian, &levelPtr[i])
		if err != nil {
			return nil, err
		}
	}

	return &mapHeader{magic: magic, levelOffset: levelPtr}, nil
}

func readLevelHeaders(mh *mapHeader, mapData []byte) ([]levelHeader, error) {
	var levelHeaders []levelHeader

	for levelIndex := 0; levelIndex < mh.LevelCount(); levelIndex++ {
		levelOffset := mh.LevelOffset(levelIndex)
		if levelOffset > 0 {
			levelDataDataBuffer := bytes.NewBuffer(mapData[levelOffset:])

			lh := levelHeader{}

			if err := binary.Read(levelDataDataBuffer, binary.LittleEndian, &lh.offPlane0); err != nil {
				return nil, err
			}
			if err := binary.Read(levelDataDataBuffer, binary.LittleEndian, &lh.offPlane1); err != nil {
				return nil, err
			}
			if err := binary.Read(levelDataDataBuffer, binary.LittleEndian, &lh.offPlane2); err != nil {
				return nil, err
			}

			if err := binary.Read(levelDataDataBuffer, binary.LittleEndian, &lh.lenPlane0); err != nil {
				return nil, err
			}
			if err := binary.Read(levelDataDataBuffer, binary.LittleEndian, &lh.lenPlane1); err != nil {
				return nil, err
			}
			if err := binary.Read(levelDataDataBuffer, binary.LittleEndian, &lh.lenPlane2); err != nil {
				return nil, err
			}
			if err := binary.Read(levelDataDataBuffer, binary.LittleEndian, &lh.width); err != nil {
				return nil, err
			}
			if err := binary.Read(levelDataDataBuffer, binary.LittleEndian, &lh.height); err != nil {
				return nil, err
			}
			nameBuffer := make([]byte, 16)
			if _, err := levelDataDataBuffer.Read(nameBuffer); err != nil {
				return nil, err
			}
			lh.name = nullTerminatedBytesToString(nameBuffer)

			levelHeaders = append(levelHeaders, lh)
		}
	}
	return levelHeaders, nil
}

func nullTerminatedBytesToString(s []byte) string {
	n := bytes.IndexByte(s, 0x00)
	if n >= 0 {
		s = s[:n]
	}
	return string(s)
}
