package utils

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DataModel struct {
	DB *sqlx.DB
}

func (d *DataModel) Init(location string) {
	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		log.Fatal("Unable to create database directory", err)
	}

	d.DB = sqlx.MustConnect("sqlite3", location+"/data.db")
}
