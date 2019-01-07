package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	hashids "github.com/speps/go-hashids"
)

var db *sql.DB
var hasher *hashids.HashID

func initDb(connection string) {
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	hasher, _ = hashids.NewWithData(hd)

}

func retrieveLongURL(token string) (string, error) {
	// unhash primary key
	// Get longurl from db that corresponds to key

	return "https://www.quora.com/What-are-some-examples-of-very-clever-long-domain-names", nil

}

func submitLongURL(longURL string) (string, error) {

	query := `INSERT INTO URL(longurl) VALUES ($1) RETURNING id`

	statement, err := db.Prepare(query)
	defer statement.Close()

	if err != nil {
		log.Println(err)
	}

	var id int64
	err = statement.QueryRow("longURL").Scan(&id)

	return hash(id), nil

}

func hash(id int64) string {
	e, _ := hasher.EncodeInt64([]int64{id})
	return e
}

func dehash(token string) int64 {
	key, _ := hasher.DecodeInt64WithError(token)
	return key[0]
}
