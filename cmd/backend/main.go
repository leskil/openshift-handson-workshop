package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type timeResponse struct {
	UnixDate string
	Host     string
}

func main() {
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
