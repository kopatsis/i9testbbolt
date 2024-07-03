package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const bucketName = "CacheBucket"

func GetExersHelper(database *mongo.Database, boltDB *bolt.DB) ([]Exercise, error) {
	var exercises []Exercise

	err := boltDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		v := b.Get([]byte("Exercise"))
		if v != nil {
			err := json.Unmarshal(v, &exercises)
			if err == nil {
				fmt.Println("Successfully bbolted")
				return nil
			}
			fmt.Println("Error bbolted, trying mongo")
			log.Printf("Failed to unmarshal from bbolt: %v, fetching from MongoDB", err)
		}

		findOptions := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
		cursor, err := database.Collection("exercise").Find(context.Background(), bson.D{}, findOptions)
		if err != nil {
			fmt.Println("Error with mongo: get")
			return err
		}
		defer cursor.Close(context.Background())

		if err = cursor.All(context.Background(), &exercises); err != nil {
			fmt.Println("Error with mongo: cursor")
			return err
		}

		data, err := json.Marshal(exercises)
		if err != nil {
			fmt.Println("Error with marshalling from mongo")
			return err
		}

		fmt.Println("Marshalled from mongo")
		err = b.Put([]byte("Exercise"), data)
		if err != nil {
			fmt.Println("Error saving to bbolt")
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	fmt.Println("Was a success")
	return exercises, nil
}

func DisConnectDB(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() (*mongo.Client, *mongo.Database, error) {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	connectStr := os.Getenv("MONGOSTRING")
	clientOptions := options.Client().ApplyURI(connectStr).SetServerAPIOptions(serverAPI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	// Specify the database and collection
	database := client.Database("i9")

	return client, database, nil
}

func initializeDB() *bolt.DB {
	db, err := bolt.Open("cache.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database initialized!")
	return db
}
