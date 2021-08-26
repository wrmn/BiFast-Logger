package main

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	// check if file doesnt exist
	file := "./sample.here"
	_, err := getConfig(file)
	if err == nil {
		t.Log("file should not found")
		t.Fail()
	}
}

func TestGetConfigFail(t *testing.T) {
	file := "./example.config.json"
	config, _ := getConfig(file)
	_, err := config.Db.dbInit()
	if err == nil {
		t.Log("should be error on unknown database")
		t.Fail()
	}
}
