package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection = db().Database("goTest").Collection("users")

func what(w http.ResponseWriter, body io.Reader, toPopulate interface{}) error {
	w.Header().Set("Content-Type", "application/json") //for adding content-type
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return json.NewDecoder(body).Decode(toPopulate) // storing in person, variable of type user
}

type user struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

//create profile
func createProfile(w http.ResponseWriter, r *http.Request) {
	var person user
	err := what(w, r.Body, &person)

	if err != nil {
		fmt.Print(err)
		//TODO: Return error message to response writer
		return
	}

	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		fmt.Print(err)
		//TODO: Return error message to response writer
		return
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID) //return the id of the mongodb document
}

//get user profile
func getUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var body user

	e := json.NewDecoder(r.Body).Decode(&body)

	if e != nil {
		fmt.Print(e)
	}

	var result primitive.M //an unordered representation of BSON document which is a map

	err := userCollection.FindOne(context.TODO(), bson.D{{"name", body.Name}}).Decode(&result)

	if err != nil {
		fmt.Print(err)
	}

	json.NewEncoder(w).Encode(result) //returns a map containing mongodb document
}

//update profile
func updateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	type updateBody struct {
		Name string `json:"name"` //value that has to be matched
		City string `json:"city"` //value that has to be modified
	}

	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
	}

	filter := bson.D{{"name", body.Name}} //converting value to bson
	after := options.After                //for returning updated document

	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.D{{"set", bson.D{{"city", body.City}}}}

	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M

	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
}

//deleteProfile
func deleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)["id"] //get Parameter value as string

	_id, err := primitive.ObjectIDFromHex(params) //convert params to mongodb hex

	if err != nil {
		fmt.Print(err.Error())
	}

	opts := options.Delete().SetCollation(&options.Collation{}) //specify language rules for string comparison

	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", _id}}, opts)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(res.DeletedCount) //return number of deleted docs
}

//get all users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var results []primitive.M                                   //slice for multiple documents
	cur, err := userCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	json.NewEncoder(w).Encode(results)
}
