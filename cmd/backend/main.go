package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var authKey string

type timeResponse struct {
	UnixDate string
	Host     string
}

func main() {

	key, err := AuthKey()

	if err != nil {
		panic(err)
	} else {
		authKey = key
	}

	r := mux.NewRouter()
	r.HandleFunc("/time", timeHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}

	fmt.Println("Backend web server started on port 8080")
	log.Fatal(srv.ListenAndServe())
}

func timeHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s\t%s", r.Method, r.URL)

	if r.URL.Query().Get("auth") != authKey {
		log.Println("Invalid auth key")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	now := time.Now()
	host, _ := os.Hostname()
	resp := timeResponse{
		UnixDate: now.Format(time.UnixDate),
		Host:     host,
	}

	json, err := json.Marshal(resp)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

// AuthKey reads the environment variable AUTH_KEY or returns an error.
func AuthKey() (string, error) {
	key := os.Getenv("AUTH_KEY")

	if key != "" {
		return key, nil
	}

	return "", errors.New("Environment variable AUTH_KEY does not exist. Set this to whatever value you want, and make sure that the frontend use the same value when calling the backend service.")
}
