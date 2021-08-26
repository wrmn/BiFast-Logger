package main

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
)

func main() {
	// initialize log
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Found error in log : ", err)
	}
	log.SetOutput(file)
	log.Println("Service started")

	config, err := getConfig("./config.json")
	if err != nil {
		log.Fatal("Error Get Config : ", err)
	}
	db, err := config.Db.dbInit()

	if err != nil {
		log.Fatal("Error initializing database : ", err)
	}

	loggingDatabase := DatabaseConf{
		Db: db,
	}
	config.Kafka.consumeFromKafka(loggingDatabase)

	defer db.Close()
}

func (conf KafkaConfig) consumeFromKafka(loggingDatabase DatabaseConf) {
	log.Println("Consumer (Kafka) started!")

	// Setting up Consumer (Kafka) config
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Broker,
		"group.id":          conf.Group,
		"auto.offset.reset": "latest",
	})

	if err != nil {
		log.Println(err.Error())
	}

	// Subscribe to topics
	c.SubscribeTopics(conf.ConsumerTopics, nil)

	// Reading message (run in loop)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Println("New Request from Kafka")
			log.Printf("Message consumed from %s: %s, Header: %s\n", msg.TopicPartition, string(msg.Value), msg.Headers[0].Value)

			commonValue := getCommonValue(msg.Value)
			status := fmt.Sprintf("received from %s", msg.TopicPartition)

			countProduct, err2 := loggingDatabase.Count(commonValue.TraceNum)

			if err2 != nil {
				log.Println(err2)
			} else {
				if countProduct == 0 {
					loggingDatabase.insertData(commonValue, status)
				} else {
					loggingDatabase.updateData(commonValue, status)
				}
			}

		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	c.Close()
}
