package dispatch

import "testing"

func TestGetCodeEase(t *testing.T) {
	testData := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	code, err := getCode(testData)
	if err != nil {
		t.Error("cant exchange!!! -> some err")
	}
	if code != 0 {
		t.Error("cant exchagotnge!! -> val err")
	}
}

func TestErr(t *testing.T) {
	testData := []byte{0, 0, 0, 0, 0}
	_, err := getCode(testData)
	if err == nil {
		t.Error("can't catch error")
	}
}

func TestGetCode(t *testing.T) {
	testData := []byte{0, 0, 1, 0, 4, 0, 0, 0}
	code, err := getCode(testData)
	if err != nil {
		t.Error("can't catch error")
	}
	if code != 4 {
		t.Error("hope:4 comming ->", code)
	}
}
