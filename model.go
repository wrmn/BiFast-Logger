package main

import "database/sql"

type DatabaseConf struct {
	Db *sql.DB
}
type DbConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Port     string `json:"port"`
}

type KafkaConfig struct {
	Broker         string   `json:"broker"`
	ConsumerTopics []string `json:"consumer_topics"`
	Group          string   `json:"group"`
}

type CommonValue struct {
	MsgDefId     string `json:"messageDefinitionId"`
	TraceNum     string `json:"traceNum"`
	CreationDate string `json:"creationDate"`
}

type Config struct {
	Kafka KafkaConfig
	Db    DbConfig
}
