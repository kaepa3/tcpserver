package main

import (
	"testing"
)

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
