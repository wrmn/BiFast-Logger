package main

import (
	"database/sql"
	"fmt"
	"log"
)

func dbInit(config DbConfig) (*sql.DB, error) {
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := sql.Open("mysql", dbConfig)

	return db, err
}

func (dBase DatabaseConf) Count(traceNum string) (int64, error) {
	query := fmt.Sprintf("select count(*) as count from logging where tracenum = %s", traceNum)
	rows, err := dBase.Db.Query(query)
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

func (dBase DatabaseConf) insertData(commonValue CommonValue, status string) {
	query := fmt.Sprintf("INSERT INTO logging VALUES ( '%s', '%s', now(), '%s', '%s', null )",
		commonValue.TraceNum, commonValue.CreationDate, status, commonValue.MsgDefId)
	dBase.runQuery(query)
}

func (dBase DatabaseConf) updateData(commonValue CommonValue, status string) {
	query := fmt.Sprintf("UPDATE `logging` SET `response_message_identifier` = '%s', `status`= '%s', `date_updated`=now() WHERE `logging`.`tracenum` = %s",
		commonValue.MsgDefId, status, commonValue.TraceNum)
	dBase.runQuery(query)

}

func (dBase DatabaseConf) runQuery(query string) {
	queryRun, err := dBase.Db.Query(query)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("success run query `%s`", query)
	defer queryRun.Close()
}
