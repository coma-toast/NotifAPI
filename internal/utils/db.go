package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DataModel struct {
	DB *sqlx.DB
}

type NotificationRow struct {
	PubID     string `db:"pub_id" json:"pub_id"`
	Date      string `db:"date" json:"date"`
	Source    string `db:"source" json:"source"`
	Interests string `db:"interests" json:"interests"`
	Title     string `db:"title" json:"title"`
	Message   string `db:"message" json:"message"`
	Metadata  string `db:"metadata" json:"metadata"`
}

type UserRow struct {
	Id                string `db:"id" json:"id"`
	Date_added        string `db:"date_added" json:"date_added"`
	Date_updated      string `db:"date_updated" json:"date_updated"`
	Username          string `db:"username" json:"username"`
	Is_admin          string `db:"is_admin" json:"is_admin"`
	Password          string `db:"password" json:"password"`
	First_name        string `db:"first_name" json:"first_name"`
	Last_name         string `db:"last_name" json:"last_name"`
	Email             string `db:"email" json:"email"`
	Account_confirmed string `db:"account_confirmed" json:"account_confirmed"`
}

func (d *DataModel) Init(location string) {
	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		log.Fatal("Unable to create database directory", err)
	}

	d.DB = sqlx.MustConnect("sqlite3", location+"/data.db")

	notifications := `CREATE TABLE IF NOT EXISTS notifications (
		pub_id text PRIMARY KEY,
		date TEXT DEFAULT CURRENT_TIMESTAMP,
		source TEXT,
		interests TEXT,
		title TEXT,
		message TEXT,
		metadata TEXT
	);`

	d.DB.MustExec(notifications)
	fmt.Println("DB Initialized: notifications")

	users := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		date_added TEXT DEFAULT CURRENT_TIMESTAMP,
		date_updated TEXT DEFAULT CURRENT_TIMESTAMP,
		username TEXT NOT NULL,
		is_admin BOOL DEFAULT false,
		password TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL,
		account_confirmed TEXT DEFAULT ""
		)`

	d.DB.MustExec(users)
	fmt.Println("DB Initialized: users")

	interests := `CREATE TABLE IF NOT EXISTS interests (
		id INTEGER PRIMARY KEY,
		date_added TEXT DEFAULT CURRENT_TIMESTAMP,
		date_updated TEXT DEFAULT CURRENT_TIMESTAMP,
		userid TEXT NOT NULL,
		interest TEXT NOT NULL,
		webhook TEXT NOT NULL
		)`

	d.DB.MustExec(interests)
	fmt.Println("DB Initialized: interests")
}

func (d *DataModel) AddNotification(pubID, source, title, message string, interests []string, metadata map[string]interface{}) (sql.Result, error) {
	insert := `INSERT INTO notifications 
	(
		pub_id,
		source,
		interests,
		title,
		message,
		metadata
	)
	VALUES 
	(
		?,
		?,
		?,
		?,
		?,
		?
	);`

	jsonInterests, err := json.Marshal(interests)
	if err != nil {
		return nil, err
	}

	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return d.DB.Exec(insert, pubID, source, jsonInterests, title, message, jsonMetadata)
}

func (d *DataModel) GetRecentNotifications(limit int) ([]NotificationRow, error) {
	notifications := make([]NotificationRow, 0)
	statement := `SELECT * FROM notifications ORDER BY date DESC LIMIT ?`
	rows, err := d.DB.Queryx(statement, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row NotificationRow
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, row)
	}

	return notifications, nil
}

func (d *DataModel) GetHistory(date time.Time) ([]NotificationRow, error) {
	notifications := make([]NotificationRow, 0)
	statement := `SELECT * FROM notifications WHERE date > ? ORDER BY date`
	test := date.Format("2006-01-02 15:04:05")
	fmt.Println(test)
	rows, err := d.DB.Queryx(statement, date.String())
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row NotificationRow
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, row)
	}

	return notifications, nil
}
