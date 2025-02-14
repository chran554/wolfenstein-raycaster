package wolf3d

import (
	"bytes"
)

// https://moddingwiki.shikadi.net/wiki/Id_Software_RLEW_compression

// RLEWDecode decodes an RLE
func RLEWDecode(source []byte) []byte {
	return RLEWDecodeWithRLEFlag(source, 0xFEFE)
}

// RLEWDecodeWithLengthPrefix decodes an RLE
func RLEWDecodeWithLengthPrefix(source []byte) (expectedSize int, data []byte) {
	size := readUint16(source, 0)
	return int(size), RLEWDecode(source[2:])
}

// RLEWDecodeWithLengthPrefixAndRLEFlag decodes (decompresses) an RLE compressed uint16 stream.
// This function should be used if the compressed stream is prefixed with uint16, little endian (two bytes) with the expected decompressed length.
// This function does not perform any consistency checks that the decompressed is of expected length.
func RLEWDecodeWithLengthPrefixAndRLEFlag(source []byte, rleFlag uint16) (expectedSize int, data []byte) {
	size := readUint16(source, 0)
	return int(size), RLEWDecodeWithRLEFlag(source[2:], rleFlag)
}

// RLEWDecodeWithRLEFlag decodes (decompresses) an RLE compressed uint16 stream.
// This function should be used if the compressed stream is prefixed with uint16, little endian (two bytes) with the expected decompressed length.
// This function does not perform any consistency checks that the decompressed is of expected length.
func RLEWDecodeWithRLEFlag(source []byte, rleFlag uint16) []byte {
	outputBuffer := bytes.Buffer{}

	var rleFlagBytes = []byte{byte(rleFlag & 0xff), byte(rleFlag >> 8)}

	var inOffset = 0

	for inOffset < len(source) {
		var word = source[inOffset : inOffset+2]
		inOffset += 2
		if bytes.Equal(word, rleFlagBytes) {
			var length = int(source[inOffset]) | int(source[inOffset+1])<<8
			var value = source[inOffset+2 : inOffset+2+2]
			inOffset += 4
			for index := 0; index < length; index++ {
				outputBuffer.Write(value)
			}
		} else {
			outputBuffer.Write(word)
		}
	}

	return outputBuffer.Bytes()
}
