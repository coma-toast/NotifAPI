package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/coma-toast/notifapi/backend/pkg/notification"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DataModel struct {
	DB *sqlx.DB
}
type NotificationRow struct {
	PubID       string `db:"pub_id" json:"pub_id"`
	Date        string `db:"date" json:"date"`
	Source      string `db:"source" json:"source"`
	Destination string `db:"destination" json:"destination"`
	Buckets     []int  `db:"buckets" json:"buckets"` // Changed from Interests to Buckets
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

type BucketRow struct {
	Id           string `db:"id" json:"id"`
	Date_added   string `db:"date_added" json:"date_added"`
	Date_updated string `db:"date_updated" json:"date_updated"`
	UserID       string `db:"userid" json:"userid"`
	Bucket       string `db:"bucket" json:"bucket"`
	Webhook      string `db:"webhook" json:"webhook"`
}

func (d *DataModel) Init(config *Config) {
	// PostgreSQL connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUser,
		config.PostgresPass,
		config.PostgresDB,
	)

	// Connect to PostgreSQL
	var err error
	d.DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v", err)
	}

	// Create tables if they don't exist
	notifications := `CREATE TABLE IF NOT EXISTS notifications (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		server TEXT,
		buckets INTEGER[] NOT NULL,
		title TEXT,
		body TEXT,
		link TEXT,
		request_data JSONB,
		metadata JSONB
	);`

	d.DB.MustExec(notifications)

	users := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		date_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		username TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT false,
		password TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL,
		account_confirmed BOOLEAN DEFAULT false
	);`

	d.DB.MustExec(users)
	fmt.Println("DB Initialized: users")

	buckets := `CREATE TABLE IF NOT EXISTS buckets (
		id SERIAL PRIMARY KEY,
		date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		date_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		userid TEXT NOT NULL,
		bucket TEXT NOT NULL,
		webhook TEXT NOT NULL
	);`

	d.DB.MustExec(buckets)
	fmt.Println("DB Initialized: buckets")
}

func (d *DataModel) AddNotification(payload notification.Message) (sql.Result, error) {
	bucketIDs := make([]int, len(payload.Buckets))
	for i, bucket := range payload.Buckets {
		var err error
		bucketRow, err := d.GetBucketByName(bucket)
		if err != nil {
			if err == sql.ErrNoRows {
				_, err := d.AddBucket(BucketRow{
					Bucket:  bucket,
					Webhook: payload.Link,
					UserID:  payload.Server,
				})
				if err != nil {
					return nil, fmt.Errorf("error adding bucket '%s': %v", bucket, err)
				}

				// Use RETURNING to get the ID of the inserted bucket
				var newBucketID int
				insertQuery := `INSERT INTO buckets (bucket, webhook, userid) VALUES ($1, $2, $3) RETURNING id`
				err = d.DB.QueryRow(insertQuery, bucket, payload.Link, payload.Server).Scan(&newBucketID)
				if err != nil {
					return nil, fmt.Errorf("error retrieving new bucket ID for '%s': %v", bucket, err)
				}
				bucketIDs[i] = newBucketID
			} else {
				return nil, fmt.Errorf("error retrieving bucket '%s': %v", bucket, err)
			}
		} else {
			bucketIDs[i], err = strconv.Atoi(bucketRow.Id)
			if err != nil {
				return nil, fmt.Errorf("invalid bucket ID '%s': %v", bucketRow.Id, err)
			}
		}
	}

	// Serialize the ipinfo.Core struct to JSON
	metadata, err := json.Marshal(payload.Metadata) // Assuming payload.Metadata is of type ipinfo.Core
	if err != nil {
		return nil, fmt.Errorf("failed to serialize metadata: %v", err)
	}
	requestData, err := json.Marshal(payload.RequestData) // Assuming payload.RequestData is of type ipinfo.Core
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request data: %v", err)
	}

	// Use the serialized JSON in the query
	insert := `INSERT INTO notifications
	(
		server,
		buckets,
		title,
		body,
		link,
		request_data,
		metadata
	)
	VALUES
	(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	);`

	// Pass the serialized JSON (metadata) as the $6 parameter
	return d.DB.Exec(insert, payload.Server, pq.Array(bucketIDs), payload.Title, payload.Body, payload.Link, string(metadata), requestData)
}

func (d *DataModel) GetRecentNotifications(limit int) ([]NotificationRow, error) {
	notifications := make([]NotificationRow, 0)
	statement := `SELECT * FROM notifications ORDER BY date DESC LIMIT $1`
	err := d.DB.Select(&notifications, statement, limit)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (d *DataModel) GetHistory(date time.Time) ([]NotificationRow, error) {
	notifications := make([]NotificationRow, 0)
	statement := `SELECT * FROM notifications WHERE date > $1 ORDER BY date`
	err := d.DB.Select(&notifications, statement, date)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (d *DataModel) GetBucketByName(name string) (BucketRow, error) {
	var returnData BucketRow
	statement := `SELECT * FROM buckets WHERE bucket = $1 ORDER BY date_updated`
	err := d.DB.Get(&returnData, statement, name)

	return returnData, err
}

func (d *DataModel) GetBucketByID(id string) (BucketRow, error) {
	var returnData BucketRow
	statement := `SELECT * FROM buckets WHERE id = $1`
	err := d.DB.Get(&returnData, statement, id)
	if err != nil {
		return BucketRow{}, err
	}
	return returnData, nil
}

func (d *DataModel) AddBucket(bucket BucketRow) (sql.Result, error) {
	insert := `INSERT INTO buckets 
	(
		bucket,
		webhook,
		userid
	)
	VALUES 
	(
		$1,
		$2,
		$3
	);`

	return d.DB.Exec(insert, bucket.Bucket, bucket.Webhook, bucket.UserID)
}

func (d *DataModel) UpdateBucket(bucket BucketRow) (sql.Result, error) {
	update := `UPDATE buckets
	SET
		bucket = $1,
		webhook = $2,
		userid = $3,
		date_updated = CURRENT_TIMESTAMP
	WHERE id = $4;`
	return d.DB.Exec(update, bucket.Bucket, bucket.Webhook, bucket.UserID, bucket.Id)
}

func (d *DataModel) DeleteBucket(id string) (sql.Result, error) {
	delete := `DELETE FROM buckets WHERE id = $1;`
	return d.DB.Exec(delete, id)
}

func (d *DataModel) GetBucketsByUserAndName(userId, name string) ([]BucketRow, error) {
	returnData := make([]BucketRow, 0)
	statement := `SELECT * FROM buckets WHERE bucket = $1 AND userid = $2 ORDER BY date_updated`
	err := d.DB.Select(&returnData, statement, name, userId)
	if err != nil {
		return nil, err
	}

	return returnData, nil
}

func (d *DataModel) GetBucketsByUser(userId string) ([]BucketRow, error) {
	returnData := make([]BucketRow, 0)
	statement := `SELECT * FROM buckets WHERE userid = $1 ORDER BY date_updated`
	err := d.DB.Select(&returnData, statement, userId)
	if err != nil {
		return nil, err
	}

	return returnData, nil
}

func (d *DataModel) InsertBucket(name, webhook, userId string) (sql.Result, error) {
	insert := `INSERT INTO buckets 
	(
		bucket,
		webhook,
		userid
	)
	VALUES 
	(
		$1,
		$2,
		$3
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
		$1,
		$2,
		$3,
		$4,
		$5
	);`

	return d.DB.Exec(insert, user.Username, user.Password, user.First_name, user.Last_name, user.Email)
}
