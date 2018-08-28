package main

import (
	"testing"
)

func TestExistsFile(t *testing.T) {
	result := Exists("./util_test.go")
	if !result {
		t.Fatalf("failed test")
	}
}
func TestNotExistsFile(t *testing.T) {
	result := Exists("./util_testXXXXXX.go")
	if result {
		t.Fatalf("failed test")
	}
}
