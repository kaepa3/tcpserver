package main

import (
	"testing"
	"time"
)

func TestExchangeCode(t *testing.T) {
	code, err := exchangeCode("e")
	if err == nil {
		t.Error("this is sting")
	}
	code, err = exchangeCode("1")
	if err != nil || code != 1 {
		t.Error("this is sting", code, ":", err)
	}
	code, err = exchangeCode("0x")
	if err == nil {
		t.Error("this is sting", code, ":", err)
	}
	code, err = exchangeCode("0x3")
	if err != nil || code != 3 {
		t.Error("this is sting", code, ":", err)
	}
	code, err = exchangeCode("0xff")
	if err != nil || code != 255 {
		t.Error("this is sting", code, ":", err)
	}
}

func TestReadConfigFile(t *testing.T) {
	config = settingGet("config.json")
	if config.Port != ":54000" {
		t.Error("cant Read by json")
	}
	if config.Health.Time == 0 {
		t.Error("cant Read health config")
	}
	if config.Health.File == "" {
		t.Error("cant Read health config")
	}
	// insertFile chk
	insertFile(4)
	if sendQue.Len() == 0 {
		t.Error("can't ADD!!")
	}
}

func TestGrantTime(t *testing.T) {
	buffer := make([]byte, 20)
	grantTime(buffer)
	now := time.Now()
	if buffer[2] != byte(now.Month()) {
		t.Error("not equal month")
	} else if buffer[3] != byte(now.Day()) {
		t.Error("not equal month")
	}
	/* second*/
	naget := buffer[8:19]
	grantTime(naget)
	now = time.Now()
	if buffer[10] != byte(now.Month()) {
		t.Error("not equal month")
	} else if buffer[11] != byte(now.Day()) {
		t.Error("not equal month")
	}

}
