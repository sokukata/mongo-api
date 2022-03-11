package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sokukata/mongo-api/mongoclient"
)

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", mongoclient.Usage)

	// Login
	myRouter.HandleFunc("/login", mongoclient.Login).Methods(http.MethodPost)
	// Create
	myRouter.HandleFunc("/add/users", mongoclient.AuthenticateMiddleware(mongoclient.AddUsers)).Methods(http.MethodPost)

	// Read
	myRouter.HandleFunc("/users/list", mongoclient.AuthenticateMiddleware(mongoclient.GetUsersList)).Methods(http.MethodGet)
	myRouter.HandleFunc("/user/{id}", mongoclient.AuthenticateMiddleware(mongoclient.GetUser)).Methods(http.MethodGet)

	// Update
	myRouter.HandleFunc("/user/{id}", mongoclient.AuthenticateMiddleware(mongoclient.UpdateUser)).Methods(http.MethodPut)

	// Delete
	myRouter.HandleFunc("/delete/user/{id}", mongoclient.AuthenticateMiddleware(mongoclient.DeleteUser)).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	servPtr := flag.String("serv", mongoclient.Server, "mongodb server. default: "+mongoclient.Server)
	dbPtr := flag.String("db", mongoclient.DatabaseName, "mongodb Database. default: "+mongoclient.DatabaseName)
	collPtr := flag.String("coll", mongoclient.Collection, "mongodb Colloection. default: "+mongoclient.Collection)

	flag.Parse()

	mongoclient.Server = *servPtr
	mongoclient.DatabaseName = *dbPtr
	mongoclient.Collection = *collPtr

	//create output directory if not exist
	_ = os.MkdirAll(mongoclient.DIRECTORY, os.ModePerm)

	handleRequests()
}
