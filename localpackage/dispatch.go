package dispatch

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const idStart int = 8
const idEnd int = 12

func GetCode(rev []byte) (uint32, error) {
	if idEnd > len(rev) {
		return 0, fmt.Errorf("length error")
	}
	var val uint32
	dataAry := rev[idStart:idEnd]
	err2 := binary.Read(bytes.NewBuffer(dataAry), binary.BigEndian, &val)
	if err2 != nil {
		return 0, fmt.Errorf("exchange error")
	}
	return val, nil
}
