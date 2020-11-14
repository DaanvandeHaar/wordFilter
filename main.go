package main

import (
	"bufio"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	//var text []string
	var words []string
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if checkIfAlpha(scanner.Text()){
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
func checkIfAlpha(s string) bool {
	const alpha = "abcdefghijklmnopqrstuvwxyz"
	for _, char := range s {
		if !strings.Contains(alpha, string(char)) {
			return false
		}
	}
	return true
}
func setWords(words []string){
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/lingoDB"))

	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("lingoDB")
	col := db.Collection("words")
	for _, word := range words {
		var _, _ = col.InsertOne(ctx, bson.D{
			{"word", word},
		})
		fmt.Println(word)
	}
}
