package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

type backendResponse struct {
	UnixDate string
	Host     string
}

type viewModel struct {
	Title           string
	UnixDate        string
	FrontendHost    string
	BackendHost     string
	AuthKey         string
	BackendEndpoint string
}

func main() {

	log.Println("Frontend web server running on port 8081...")

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

	//tmpl := template.Must(template.ParseFiles("layout.html"))
	tmpl := template.Must(template.New("layout").Parse(html()))
	host, _ := os.Hostname()
	vm := viewModel{
		Title:           "Testing layout",
		UnixDate:        data.UnixDate,
		FrontendHost:    host,
		BackendHost:     data.Host,
		AuthKey:         os.Getenv("AUTH_KEY"),
		BackendEndpoint: os.Getenv("BACKEND_ENDPOINT"),
	}

	tmpl.Execute(w, vm)
}

func callBackend() (*backendResponse, error) {
	key, err := AuthKey()

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

// AuthKey reads the environment variable AUTH_KEY or returns an error.
func AuthKey() (string, error) {
	key := os.Getenv("AUTH_KEY")

	if key != "" {
		return key, nil
	}

	return "", errors.New("Environment variable AUTH_KEY does not exist. Make sure it's using the same value as the backend service.")
}

func html() string {
	return `<!DOCTYPE html>
	<html>
	
	<head>
		<title>OpenShift hands on</title>
		<meta charset="UTF-8">
		<link rel="stylesheet" type="text/css"
			href="https://cdnjs.cloudflare.com/ajax/libs/patternfly/3.24.0/css/patternfly.min.css">
		<link rel="stylesheet" type="text/css"
			href="https://cdnjs.cloudflare.com/ajax/libs/patternfly/3.24.0/css/patternfly-additions.min.css">
	
	</head>
	
	<body style="background-color: rgb(245, 245, 245); padding: 20px;">
		<div class="container">
			<div class="row row-cards-pf">
				<h1 style="text-align: center; margin-bottom: 20px">OpenShift hands-on frontend</h1>
				<div class="card-pf">
					<h1 class="card-pf-title">Configuration</h1>
					<div class="card-pf-body">
						<p>
						<dl>
							<dt>Frontend host/pod:</dt>
							<dd>{{.FrontendHost}}</dd>
							<dt>Auth key:</dt>
							<dd>{{.AuthKey}}</dd>
							<dt>Endpoint:</dt>
							<dd>{{.BackendEndpoint}}</dd>                        
						</dl>
						</p>
					</div>
				</div>
	
				<div class="card-pf">
					<h1 class="card-pf-title">Backend service results</h1>
					<div class="card-pf-body">
						<p>
						<dl>
							<dt>UnixTime:</dt>
							<dd>{{.UnixDate}}</dd>
							<dt>Backend host/pod:</dt>
							<dd>{{.BackendHost}}</dd>
						</dl>
						</p>
					</div>
				</div>
	
			</div>
		</div>
		</div>
	</body>
	
	</html>

`
}
