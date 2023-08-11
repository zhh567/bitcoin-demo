package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func UintToByte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.LittleEndian, num)
	if err != nil {
		fmt.Println("convert uint to byte fail: ", err)
		return nil
	}

	return buffer.Bytes()
}

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
