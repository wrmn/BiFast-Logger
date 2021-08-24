package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func getConfig() (configDatabase DbConfig, configKafka KafkaConfig) {
	// Return config for setting up Kafka Producer and Consumer
	log.Printf("Get config database")
	file, _ := os.Open("./config.json")
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

func (productModel ProductModel) Count(traceNum string) (int64, error) {
	query := fmt.Sprintf("select count(*) as count from logging where tracenum = %s", traceNum)
	rows, err := productModel.Db.Query(query)
	if err != nil {
		return 0, err
	} else {
		var count_product int64
		for rows.Next() {
			rows.Scan(&count_product)
		}
		return count_product, nil
	}
}
