package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/leskil/openshift-handson-workshop/pkg/config"
)

type backendResponse struct {
	UnixDate string
	Host     string
}

type viewModel struct {
	Title        string
	UnixDate     string
	FrontendHost string
	BackendHost  string
	AuthKey      string
}

func main() {

	log.Println("Frontend web server running...")

	http.HandleFunc("/", renderTemplate)
	http.ListenAndServe(":8081", nil)
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.URL)

	data, err := callBackend()

	if err != nil {
		log.Print(err)
		return
	}

	tmpl := template.Must(template.ParseFiles("layout.html"))
	host, _ := os.Hostname()
	vm := viewModel{
		Title:        "Testing layout",
		UnixDate:     data.UnixDate,
		FrontendHost: host,
		BackendHost:  data.Host,
		AuthKey:      os.Getenv("AUTH_KEY"),
	}

	tmpl.Execute(w, vm)
}

func callBackend() (*backendResponse, error) {
	key, err := config.AuthKey()

	if err != nil {
		return nil, err
	}

	url := backendEndpoint() + "/time" + "?auth=" + key
	resp, err := http.Get(url)

	log.Printf("Response from %s: %s", url, resp.Status)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	r := &backendResponse{}
	err = json.Unmarshal(body, r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func backendEndpoint() string {
	endp := os.Getenv("BACKEND_ENDPOINT")

	if endp == "" {
		panic(errors.New("Missing environment variable: BACKEND_ENDPOINT"))
	}

	return endp
}
