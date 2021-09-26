package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"covidApi/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Statecases(state string) models.Responsejson {
	clientOptions := options.Client().ApplyURI("mongodb+srv://yogesh9643:Chauhan123@cluster0.kkure.mongodb.net/coviddb?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	db := client.Database("covid")
	stateData := db.Collection("StateData")
	if err != nil {
		log.Fatal(err)
	}
	var data models.State
	filter := bson.D{{"state", state}}
	if err = stateData.FindOne(ctx, filter).Decode(&data); err != nil {
		log.Fatal(err)
	}

	var total models.State
	filter = bson.D{{"state", "Total"}}
	if err = stateData.FindOne(ctx, filter).Decode(&total); err != nil {
		log.Fatal(err)
	}
	var indiatotal models.Responsejson
	indiatotal.Active = data.Active
	indiatotal.Confirmed = data.Confirmed
	indiatotal.Lastupdatedtime = data.Lastupdatedtime
	indiatotal.State = data.State
	indiatotal.Totalincountry = total.Confirmed

	fmt.Print(indiatotal)
	return indiatotal
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

	var statename models.StateMapBox
	json.Unmarshal([]byte(body), &statename)
	return statename.StateMapBox[0].StateName
}
