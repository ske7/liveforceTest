package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("start test server")

	// make DB connection and create needed table
	_, err := SetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// define router
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginHandler)           // post login data, example {"Email": "user@email.com", "Password": "password"}
	router.HandleFunc("/logout", LogoutHandler)         // clear stored session
	router.HandleFunc("/adduser", AddUserHandler)       // add user, post data {"Name": "name", "Email": "user@email.com", "Password": "password", "Role": "user"}
	router.HandleFunc("/deleteuser", DeleteUserHandler) // delete user, post data {"Email": "user@email.com"}
	router.HandleFunc("/getusers", GetUsersHandler)     // get users table
	http.Handle("/", router)

	// server job start
	listenAddr := "127.0.0.1:8000"
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
