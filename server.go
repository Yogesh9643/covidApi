package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"io/ioutil"

	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Coordinate is a [longitude, latitude]
type Users struct {
	Users []User `json:"statewise"`
}

type User struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	State           string `json:"state"`
}

type Coordinate [2]float64

// Connection URI
const uri = "mongodb+srv://yogesh9643:Chauhan123@cluster0.kkure.mongodb.net/coviddb?retryWrites=true&w=majority"

type StateMapBox struct {
	StateMapBox []StateCoor `json:"features"`
}

type StateCoor struct {
	StateName string `json:"text"`
}

func GetStateData(longitude string, latitude string) string {
	url := "https://api.mapbox.com/geocoding/v5/mapbox.places/" + longitude + "," + latitude + ".json?types=region&access_token=pk.eyJ1IjoieW9nZXNoOTY0MyIsImEiOiJja3A3NjBlaXIwNmpvMnZtcnZ0ZDJ2dGxmIn0.yw2cd7raPthzGXL6BW3Tnw"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print("error occures")
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var statename StateMapBox
	json.Unmarshal([]byte(body), &statename)
	return statename.StateMapBox[0].StateName
}

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

	var users Users
	json.Unmarshal(body, &users)

	db := client.Database("covid")
	stateData := db.Collection("StateData")

	for i := 0; i < len(users.Users); i++ {
		Result, err := stateData.InsertOne(ctx, bson.D{
			{Key: "state", Value: users.Users[i].State},
			{Key: "active", Value: users.Users[i].Active},
			{Key: "confirmed", Value: users.Users[i].Confirmed},
			{Key: "lastupdatedtime", Value: users.Users[i].Lastupdatedtime},
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Inserted %v documents into State collection!\n", Result.InsertedID)
	}
}
func stateCases(state string) User {
	clientOptions := options.Client().ApplyURI("mongodb+srv://yogesh9643:Chauhan123@cluster0.kkure.mongodb.net/coviddb?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	db := client.Database("covid")
	stateData := db.Collection("StateData")
	if err != nil {
		log.Fatal(err)
	}
	var data User
	filter := bson.D{{"state", state}}
	if err = stateData.FindOne(ctx, filter).Decode(&data); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Confirmed")
	fmt.Print(data)
	return data
}
func update() {
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

	var users Users
	json.Unmarshal(body, &users)

	db := client.Database("covid")
	stateData := db.Collection("StateData")

	for i := 0; i < len(users.Users); i++ {
		filter := bson.D{{"state", users.Users[i].State}}

		update := bson.D{
			{"$set", bson.D{
				{"confirmed", users.Users[i].Confirmed},
				{"active", users.Users[i].Active},
				{"lastupdatedtime", users.Users[i].Lastupdatedtime},
			}},
		}
		stateData.UpdateOne(ctx, filter, update)

		//fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	}
}
func fetchtoDB(c echo.Context) error {
	update()
	fmt.Print("updated")
	return c.String(http.StatusOK, "Updated")
}
func stateCovid(c echo.Context) error {
	// Get team and member from the query string
	var x string = c.QueryParam("longitude")
	var y string = c.QueryParam("latitude")

	//longitude, err := strconv.ParseFloat(x, 8)
	//latitude, err := strconv.ParseFloat(y, 8)

	//log.Printf(longitude, latitude, err)
	fmt.Print(x, y)
	var cases User
	state := GetStateData(x, y)
	fmt.Print(state)
	cases = stateCases(state)
	fmt.Print(cases)
	casesMarshal, err := json.Marshal(cases)
	fmt.Print(err)
	return c.String(http.StatusOK, string(casesMarshal))
}
func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/state", stateCovid)
	e.GET("/fetchtodb", fetchtoDB)
	e.Logger.Fatal(e.Start(":1323"))
}
