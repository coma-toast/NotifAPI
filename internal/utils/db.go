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
	PubID       string `db:"pub_id" json:"pub_id"`
	Date        string `db:"date" json:"date"`
	Source      string `db:"source" json:"source"`
	Destination string `db:"destination" json:"destination"`
	Interests   string `db:"interests" json:"interests"`
	Title       string `db:"title" json:"title"`
	Message     string `db:"message" json:"message"`
	Metadata    string `db:"metadata" json:"metadata"`
}

type User struct {
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

type InterestRow struct {
	Id           string `db:"id" json:"id"`
	Date_added   string `db:"date_added" json:"date_added"`
	Date_updated string `db:"date_updated" json:"date_updated"`
	UserID       string `db:"userid" json:"userid"`
	Interest     string `db:"interest" json:"interest"`
	Webhook      string `db:"webhook" json:"webhook"`
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
		destination TEXT,
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

func (d *DataModel) AddNotification(pubID, source, destination, title, message string, interests []string, metadata map[string]interface{}) (sql.Result, error) {
	insert := `INSERT INTO notifications 
	(
		pub_id,
		source,
		destination,
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

	return d.DB.Exec(insert, pubID, source, destination, jsonInterests, title, message, jsonMetadata)
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
	// test := date.Format("2006-01-02 15:04:05")
	// fmt.Println(test)
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

func (d *DataModel) GetInterestByName(name string) (InterestRow, error) {
	var returnData InterestRow
	statement := `SELECT * FROM interests WHERE interest = ? ORDER BY date_updated`
	row := d.DB.QueryRowx(statement, name)
	err := row.StructScan(&returnData)

	return returnData, err
}

func (d *DataModel) GetInterestsByUserAndName(userId, name string) ([]InterestRow, error) {
	returnData := make([]InterestRow, 0)
	statement := `SELECT * FROM interests WHERE interest = ? AND userid = ? ORDER BY date_updated`
	rows, err := d.DB.Queryx(statement, name, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row InterestRow
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		returnData = append(returnData, row)
	}
	return returnData, err
}

func (d *DataModel) GetInterestsByUser(userId string) ([]InterestRow, error) {
	returnData := make([]InterestRow, 0)
	statement := `SELECT * FROM interests WHERE userid = ? ORDER BY date_updated`
	rows, err := d.DB.Queryx(statement, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var row InterestRow
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		returnData = append(returnData, row)
	}

	return returnData, err
}

func (d *DataModel) InsertInterest(name, webhook, userId string) (sql.Result, error) {
	insert := `INSERT INTO interests 
	(
		interest,
		webhook,
		userid
	)
	VALUES 
	(
		?,
		?,
		?
	);`

	return d.DB.Exec(insert, name, webhook, userId)
}

func (d *DataModel) AddUser(user User) (sql.Result, error) {
	user.Password = HashPassword(user.Password)
	insert := `INSERT INTO users 
	(
		username,
		password,
		first_name,
		last_name,
		email
	)
	VALUES 
	(
		?,
		?,
		?,
		?,
		?
	);`

	return d.DB.Exec(insert, user.Username, user.Password, user.First_name, user.Last_name, user.Email)
}
