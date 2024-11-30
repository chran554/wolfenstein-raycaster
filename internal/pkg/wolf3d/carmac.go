package wolf3d

import (
	"bytes"
	"encoding/binary"
)

// https://moddingwiki.shikadi.net/wiki/Carmack_compression

// CarmackDecodeWithLengthPrefix decodes (decompresses) a Carmack compressed uint16 stream.
// This function should be used if the compressed stream is prefixed with uint16, little endian (two bytes) with the expected decompressed length.
// This function does not perform any consistency checks that the decompressed is of expected length.
func CarmackDecodeWithLengthPrefix(source []byte) (expectedSize int, data []byte) {
	size := readUint16(source, 0)
	return int(size), CarmackDecode(source[2:])
}

func CarmackDecode(source []byte) []byte {
	const nearPointerMarker = byte(0xA7)
	const farPointerMarker = byte(0xA8)

	destBuffer := bytes.Buffer{}
	sourceOffset := 0

	for sourceOffset < len(source)-1 {

		possibleNearPointer := false
		possibleFarPointer := false

		var sourceLookahead0 = source[sourceOffset]
		var sourceLookahead1 = source[sourceOffset+1]

		possibleNearPointer = sourceLookahead1 == nearPointerMarker
		possibleFarPointer = sourceLookahead1 == farPointerMarker

		nearPointer := sourceLookahead0 != 0x00 && possibleNearPointer
		farPointer := sourceLookahead0 != 0x00 && possibleFarPointer

		if (possibleNearPointer || possibleFarPointer) && !(nearPointer || farPointer) {
			destBuffer.WriteByte(source[sourceOffset+2])
			destBuffer.WriteByte(source[sourceOffset+1])
			sourceOffset += 3
		}

		if nearPointer {
			var pointerOffset = 2 * int(source[sourceOffset+2])
			destBytes := destBuffer.Bytes()
			destReadPos := len(destBytes) - pointerOffset
			for i := 0; i < int(sourceLookahead0); i++ {
				destBuffer.WriteByte(destBytes[destReadPos+i*2])
				destBuffer.WriteByte(destBytes[destReadPos+i*2+1])
			}

			sourceOffset += 3
		}

		if farPointer {
			var wordOffset uint16
			_ = binary.Read(bytes.NewBuffer(source[sourceOffset+2:sourceOffset+4]), binary.LittleEndian, &wordOffset)
			var pointerOffset = 2 * int(wordOffset)
			destBytes := destBuffer.Bytes()
			for i := 0; i < int(sourceLookahead0); i++ {
				destBuffer.WriteByte(destBytes[pointerOffset+2*i])
				destBuffer.WriteByte(destBytes[pointerOffset+2*i+1])
			}

			sourceOffset += 4
		}

		if !(possibleNearPointer || possibleFarPointer) {
			destBuffer.WriteByte(source[sourceOffset])
			destBuffer.WriteByte(source[sourceOffset+1])
			sourceOffset += 2
		}
	}

	return destBuffer.Bytes()
}
