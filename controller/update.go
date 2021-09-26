package controller

import (
	"context"
	"covidApi/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func fetch() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://yogesh9643:Chauhan123@cluster0.kkure.mongodb.net/coviddb?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://data.covid19india.org/data.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var statelist models.Statelist
	json.Unmarshal(body, &statelist)

	db := client.Database("covid")
	stateData := db.Collection("StateData")

	for i := 0; i < len(statelist.Statelist); i++ {
		Result, err := stateData.InsertOne(ctx, bson.D{
			{Key: "state", Value: statelist.Statelist[i].State},
			{Key: "active", Value: statelist.Statelist[i].Active},
			{Key: "confirmed", Value: statelist.Statelist[i].Confirmed},
			{Key: "lastupdatedtime", Value: statelist.Statelist[i].Lastupdatedtime},
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Inserted %v documents into State collection!\n", Result.InsertedID)
	}
}

func Update() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://yogesh9643:Chauhan123@cluster0.kkure.mongodb.net/coviddb?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://data.covid19india.org/data.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var statelist models.Statelist
	json.Unmarshal(body, &statelist)

	db := client.Database("covid")
	stateData := db.Collection("StateData") //StateData
	fmt.Println("Length")

	var data models.State
	filter := bson.D{{"state", "Delhi"}}
	if err = stateData.FindOne(ctx, filter).Decode(&data); err != nil {
		fmt.Print(err)
		fetch()

	} else {
		for i := 0; i < len(statelist.Statelist); i++ {
			filter := bson.D{{"state", statelist.Statelist[i].State}}

			update := bson.D{
				{"$set", bson.D{
					{"confirmed", statelist.Statelist[i].Confirmed},
					{"active", statelist.Statelist[i].Active},
					{"lastupdatedtime", statelist.Statelist[i].Lastupdatedtime},
				}},
			}
			stateData.UpdateOne(ctx, filter, update)

			//fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

		}
	}
	fmt.Print(data)

}
