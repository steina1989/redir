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
	log.Println("Connection to database successful")

	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	hasher, _ = hashids.NewWithData(hd)

}

func retrieveLongURL(token string) (string, error) {
	// unhash primary key
	// Get longurl from db that corresponds to key
	var longurl string
	id := dehash(token)
	query := `SELECT longurl from url WHERE ID = ($1)`

	err := db.QueryRow(query, id).Scan(&longurl)

	if err != nil {
		log.Println(err)
	}

	return longurl, nil

}

func submitLongURL(longURL string) (string, error) {

	query := `INSERT INTO url(longurl) VALUES ($1) RETURNING id`

	statement, qerr := db.Prepare(query)
	defer statement.Close()

	if qerr != nil {
		log.Println(qerr)
	}

	var id int64
	iderr := statement.QueryRow(longURL).Scan(&id)

	if iderr != nil {
		log.Println(iderr)
	}

	return hash(id), nil

}

func hash(id int64) string {
	e, err := hasher.EncodeInt64([]int64{id})
	if err != nil {
		log.Println(err)
	}
	return e
}

func dehash(token string) int64 {
	key, err := hasher.DecodeInt64WithError(token)
	if err != nil {
		log.Println(err)
	}
	return key[0]
}
