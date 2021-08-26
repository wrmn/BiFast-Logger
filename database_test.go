package main

import "testing"

func TestCounting(t *testing.T) {
	config, _ := getConfig("./config.json")
	db, _ := dbInit(config.Db)

	loggingDatabase := DatabaseConf{
		Db: db,
	}

	res, err := loggingDatabase.Count("123456")

	if err != nil {
		t.Log("error should be nil", err)
		t.Fail()
	}

	if res > 1 {
		t.Log("Data should get 1 count at max")
		t.Fail()
	}
}
