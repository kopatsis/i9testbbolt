package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	client, database, err := ConnectDB()
	if err != nil {
		log.Fatalf("Error while connecting to mongoDB: %s.\nExiting.", err)
	}
	defer DisConnectDB(client)

	db := initializeDB()
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter 1 to run bbolt mongo test: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}
		if number == 1 {
			ret, err := GetExersHelper(database, db)
			if err != nil {
				fmt.Println("Error from bbolt or mongo.")
			} else {
				fmt.Println(ret[int(rand.Float32()*float32(len(ret)))].Name)
			}
		}
	}
}
