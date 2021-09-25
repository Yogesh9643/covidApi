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
type Statelist struct {
	Statelist []State `json:"statewise"`
}

type State struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	State           string `json:"state"`
}

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

	var statelist Statelist
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
func stateCases(state string) State {
	clientOptions := options.Client().ApplyURI("mongodb+srv://yogesh9643:Chauhan123@cluster0.kkure.mongodb.net/coviddb?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	db := client.Database("covid")
	stateData := db.Collection("StateData")
	if err != nil {
		log.Fatal(err)
	}
	var data State
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

	var statelist Statelist
	json.Unmarshal(body, &statelist)

	db := client.Database("covid")
	stateData := db.Collection("StateData")

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
	var cases State
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
