package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"
)

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func FileExists(path string) bool {
	_, err := os.Stat(path) //Stat returns metadata for said file but does not open it
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	Handle(err)
	return false
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
