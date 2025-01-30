package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/pusher/pusher-http-go"
)

type loggingTransport struct {
	Transport http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Println("Error dumping request:", err)
	} else {
		log.Println("REQUEST:")
		log.Println(string(requestDump))
	}
	resp, err := t.Transport.RoundTrip(req)
	return resp, err
}

var client = pusher.Client{
	AppId:  "local",
	Secret: "3ueVhNtuicMeJpYq",
	Key:    "d6kAd89bMqDrLrFh",
	Secure: false,
	Host:   "127.0.0.1:8081",
	HttpClient: &http.Client{
		Transport: &loggingTransport{
			Transport: http.DefaultTransport,
		},
	},
}

type Body struct {
	Message string `json:"message"`
}

func handelFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	var body Body
	json.NewDecoder(r.Body).Decode(&body)
	
	err := client.Trigger("test", "publish", body.Message)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Message sent successfully!")
	}
	w.Write([]byte("h"))
}

func main() {

	http.HandleFunc("/send", handelFunc)
	http.ListenAndServe(":8000", nil)

}
