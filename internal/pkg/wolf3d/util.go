package wolf3d

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

func readUint16(data []byte, offset int) uint16 {
	var value uint16
	if err := binary.Read(bytes.NewBuffer(data[offset:offset+2]), binary.LittleEndian, &value); err != nil {
		panic(1)
	}

	return value
}

func writeToFile(data []byte, filename string) {
	file, _ := os.Create(filename)
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	_, _ = writer.Write(data)
}

func readFile(filename string) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
