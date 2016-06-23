package main

import (
	"testing"
)

func TestConstructorWithFullname(t *testing.T) {
	setting := settingGet("config.json")
	if setting.Port != ":50000" {
		t.Error("cant Read by json")
	}
}
