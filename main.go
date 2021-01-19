package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Print("Starting server . . .")
	route := mux.NewRouter()

	s := route.PathPrefix("/api").Subrouter() //base path.

	//Routes
	s.HandleFunc("/createProfile", createProfile).Methods("POST")
	s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	s.HandleFunc("/getUserProfile", getUserProfile).Methods("POST")
	s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")

	route.Use(mux.CORSMethodMiddleware(route))
	port := 8000
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), s)) //run server
}
