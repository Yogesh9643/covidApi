package controller

import (
	"context"
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

type StateMapBox struct {
	StateMapBox []StateCoor `json:"features"`
}

type StateCoor struct {
	StateName string `json:"text"`
}

func Statecases(state string) State {
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

func CordToState(longitude string, latitude string) string {
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
