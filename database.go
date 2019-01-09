package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	hashids "github.com/speps/go-hashids"
)

var db *sql.DB
var hasher *hashids.HashID

// InitDb attempts to connect to postgres database with supplied connection string
// It also initializes the hash function used for the server
func InitDb(connection string) {
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}
}

// InitHash initializes hash function
func InitHash(salt string) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength

	var err error
	hasher, err = hashids.NewWithData(hd)
	if err != nil {
		log.Fatal(err)
	}
}

// RetrieveLongURL dehashes a given token to a Primary key in the database
// Error messages are logged and abstracted error messages for end users are returned.
func RetrieveLongURL(id int64) (string, error) {
	var longurl string

	query := `SELECT longurl from url WHERE ID = ($1)`

	dberr := db.QueryRow(query, id).Scan(&longurl)

	if dberr != nil {
		log.Println(dberr)
		return "", errors.New("Entry not found")
	}

	return longurl, nil

}

// SubmitLongURL adds a new entry to the database, and returns the hash of its primary key.
// Errors are logged, and returned in an abstracted form.
func SubmitLongURL(longURL string) (string, error) {

	query := `INSERT INTO url(longurl) VALUES ($1) RETURNING id`

	dbErr := errors.New("Database error")

	statement, qerr := db.Prepare(query)
	defer statement.Close()

	if qerr != nil {
		log.Println(qerr)
		return "", dbErr
	}

	var id int64
	iderr := statement.QueryRow(longURL).Scan(&id)

	if iderr != nil {
		log.Println(iderr)
		return "", dbErr
	}

	return hash(id), nil

}

func hash(id int64) string {
	e, _ := hasher.EncodeInt64([]int64{id})
	return e
}

// Dehash converts a token to a Primary key
func Dehash(token string) (int64, error) {
	key, err := hasher.DecodeInt64WithError(token)
	if err != nil {
		log.Println(err)
		return int64(0), errors.New("Invalid token")
	}
	return key[0], nil
}
