package main

import (
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	_, err := dbInit(config.Db)
	if err == nil {
		t.Log("should be error on unknown database")
		t.Fail()
	}

	_, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.Broker,
		"group.id":          config.Kafka.Group,
		"auto.offset.reset": "latest",
	})
	if err == nil {
		t.Log("should be error on unknown broker")
		t.Fail()
	}

	c, _ := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "testLog",
		"auto.offset.reset": "latest",
	})
	err = c.SubscribeTopics(config.Kafka.ConsumerTopics, nil)
	if err == nil {
		t.Log("should be error on unknown topic")
		t.Fail()
	}
}
