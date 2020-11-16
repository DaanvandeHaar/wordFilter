package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

func main() {
	// declare slice of words
	var words []string

	//open words.txt file
	file, err := os.Open("words.txt")
	//check for errors
	if err != nil {
		log.Fatal(err)
	}
	//keep file open
	defer file.Close()

	//open new scanner for words.txt
	scanner := bufio.NewScanner(file)

	//check for non alpha chars and select words with 5-7 char len words
	for scanner.Scan() {
		if checkIfAlpha(scanner.Text()) {
			if len(scanner.Text()) >= 5 && len(scanner.Text()) <= 7 {
				words = append(words, scanner.Text())
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	setWords(words)
}

//check if all chars in word are letters
func checkIfAlpha(s string) bool {
	const alpha = "abcdefghijklmnopqrstuvwxyz"
	for _, char := range s {
		if !strings.Contains(alpha, string(char)) {
			return false
		}
	}
	return true
}
func setWords(words []string) {
	//declare db string
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "admin"
		dbname   = "lingo_db"
	)
	//make connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// check for err and open connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// keep db from closing conn
	defer db.Close()

	//prepare stm

	stm := "INSERT INTO words (word) VALUES ($1)"

	//for loop to add valid words to db
	for _, word := range words {
		_, err := db.Exec(stm, word)
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Println(word)
	}
	db.Close()
}
