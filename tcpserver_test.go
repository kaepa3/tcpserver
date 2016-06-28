package main

import (
	"testing"
)

func ReadConfigFile_test(t *testing.T) {
	setting := settingGet("config.json")
	if setting.Port != ":50000" {
		t.Error("cant Read by json")
	}
	if setting.Health.Time == 0 {
		t.Error("cant Read health config")
	}
	if setting.Health.File == "" {
		t.Error("cant Read health config")
	}
}
