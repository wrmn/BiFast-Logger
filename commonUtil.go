package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func getConfig(fileName string) (config Config, err error) {
	// Return config for setting up Kafka Producer and Consumer
	log.Printf("Get config database")
	file, err := os.Open(fileName)
	if err != nil {
		return config, err
	}

	defer file.Close()

	// Read config file
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return config, err
	}

	// Unmarshal config to variable
	err = json.Unmarshal(b, &config.Kafka)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(b, &config.Db)
	if err != nil {
		return config, err
	}

	log.Printf("Host: `%v`, username: `%v`, database: `%v`, port: `%v`",
		config.Db.Host, config.Db.Username, config.Db.Database, config.Db.Port)
	log.Printf("Kafka Config -> Broker: `%v`, Consumer Topics: `%v`, Group: `%v`",
		config.Kafka.Broker, config.Kafka.ConsumerTopics, config.Kafka.Group)

	return config, nil
}

func getCommonValue(content []byte) (res CommonValue) {
	json.Unmarshal(content, &res)
	return res
}
