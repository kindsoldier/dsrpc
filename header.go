/*
 *
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 *
 */

package dsrpc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
)

const (
	headerSize  int64 = 16 * 2
	sizeOfInt64 int   = 8
	magicCodeA  int64 = 0xEE00ABBA
	magicCodeB  int64 = 0xEE44ABBA
)

type Header struct {
	magicCodeA int64 `json:"magicCodeA"`
	rpcSize    int64 `json:"rpcSize"`
	binSize    int64 `json:"binSize"`
	magicCodeB int64 `json:"magicCodeB"`
}

func NewEmptyHeader() *Header {
	return &Header{
		magicCodeA: magicCodeA,
		magicCodeB: magicCodeB,
	}
}

func (hdr *Header) ToJson() []byte {
	jBytes, _ := json.Marshal(hdr)
	return jBytes
}

func (hdr *Header) Pack() ([]byte, error) {
	var err error
	headerBytes := make([]byte, 0, headerSize)
	headerBuffer := bytes.NewBuffer(headerBytes)

	magicCodeABytes := EncoderI64(hdr.magicCodeA)
	headerBuffer.Write(magicCodeABytes)

	rpcSizeBytes := EncoderI64(hdr.rpcSize)
	headerBuffer.Write(rpcSizeBytes)

	binSizeBytes := EncoderI64(hdr.binSize)
	headerBuffer.Write(binSizeBytes)

	magicCodeBBytes := EncoderI64(hdr.magicCodeB)
	headerBuffer.Write(magicCodeBBytes)

	return headerBuffer.Bytes(), err
}

func UnpackHeader(headerBytes []byte) (*Header, error) {
	var err error

	headerReader := bytes.NewReader(headerBytes)

	magicCodeABytes := make([]byte, sizeOfInt64)
	headerReader.Read(magicCodeABytes)

	rpcSizeBytes := make([]byte, sizeOfInt64)
	headerReader.Read(rpcSizeBytes)

	binSizeBytes := make([]byte, sizeOfInt64)
	headerReader.Read(binSizeBytes)

	magicCodeBBytes := make([]byte, sizeOfInt64)
	headerReader.Read(magicCodeBBytes)

	header := &Header{
		magicCodeA: DecoderI64(magicCodeABytes),
		rpcSize:    DecoderI64(rpcSizeBytes),
		binSize:    DecoderI64(binSizeBytes),
		magicCodeB: DecoderI64(magicCodeBBytes),
	}

	if header.magicCodeA != magicCodeA || header.magicCodeB != magicCodeB {
		err = errors.New("Wrong protocol magic code")
		return header, err
	}
	return header, err
}

func EncoderI64(i int64) []byte {
	buffer := make([]byte, sizeOfInt64)
	binary.BigEndian.PutUint64(buffer, uint64(i))
	return buffer
}

func DecoderI64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}
