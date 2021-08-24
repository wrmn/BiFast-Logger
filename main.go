package main

import (
	"database/sql"
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
		log.Fatal("Found error in log ", err)
	}
	log.SetOutput(file)
	log.Println("Service started")

	configDatabase, configKafka := getConfig()
	db := dbInit(configDatabase)
	consumeFromKafka(configKafka, db)

	defer db.Close()
}

func dbInit(config DbConfig) *sql.DB {
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := sql.Open("mysql", dbConfig)

	if err != nil {
		panic(err.Error())
	}

	return db
}

func consumeFromKafka(conf KafkaConfig, db *sql.DB) {
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

			productModel := ProductModel{
				Db: db,
			}

			countProduct, err2 := productModel.Count(commonValue.TraceNum)

			if err2 != nil {
				fmt.Println(err2)
			} else {
				if countProduct == 0 {
					query := fmt.Sprintf("INSERT INTO logging VALUES ( '%s', now(), now(), '%s', '%s', null )",
						commonValue.TraceNum, status, commonValue.MsgDefId)

					insert, err := db.Query(query)

					// if there is an error inserting, handle it
					if err != nil {
						log.Println(err.Error())
					}
					// be careful deferring Queries if you are using transactions
					defer insert.Close()
				} else {
					query := fmt.Sprintf("UPDATE `logging` SET `response_message_identifier` = '%s', `status`= '%s', `date_updated`=now() WHERE `logging`.`tracenum` = %s",
						commonValue.MsgDefId, status, commonValue.TraceNum)

					update, err := db.Query(query)

					// if there is an error inserting, handle it
					if err != nil {
						log.Println(err.Error())
					}
					// be careful deferring Queries if you are using transactions
					defer update.Close()
				}
			}

		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	c.Close()
}
