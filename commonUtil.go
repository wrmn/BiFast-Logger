package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func getConfig(fileName string) (configDatabase DbConfig, configKafka KafkaConfig) {
	// Return config for setting up Kafka Producer and Consumer
	log.Printf("Get config database")
	file, _ := os.Open(fileName)
	defer file.Close()

	// Read config file
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal config to variable
	json.Unmarshal(b, &configDatabase)
	json.Unmarshal(b, &configKafka)

	log.Printf("Host: `%v`, username: `%v`, database: `%v`, port: `%v`",
		configDatabase.Host, configDatabase.Username, configDatabase.Database, configDatabase.Port)
	log.Printf("Kafka Config -> Broker: `%v`, Consumer Topics: `%v`, Group: `%v`",
		configKafka.Broker, configKafka.ConsumerTopics, configKafka.Group)

	return configDatabase, configKafka
}

func getCommonValue(content []byte) (res CommonValue) {
	json.Unmarshal(content, &res)
	return res
}
