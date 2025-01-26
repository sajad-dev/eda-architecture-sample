package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type MsgStruct struct {
	Message string `json:"message"`
}

func handler(msg chan string, handlerFunc func(http.ResponseWriter, *http.Request, chan string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc(w, r, msg)
	}
}

func controller(w http.ResponseWriter, r *http.Request, msg chan string) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    var message MsgStruct
    if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    msg <- message.Message
    w.Write([]byte("ok"))
}

func main() {
	godotenv.Load(".env")
	
	url := fmt.Sprintf("ws://%s:%s/%s",os.Getenv("IP"),os.Getenv("WEBSOCKET_PORT"),os.Getenv("PUBLIC_KEY"))
	headers := http.Header{}
	headers.Add("secret_key",os.Getenv("SECRET_KEY"))

	conn, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		fmt.Println("Error connecting to WebSocket:", err)
	}

	message := make(chan string)

	go func() {
		for {

			select {
			case msg := <-message:
				err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					fmt.Println("Error sending message:", err)
				}

			}
		}
	}()

	http.HandleFunc("/send", handler(message, controller))
	http.ListenAndServe(fmt.Sprintf(":%s",os.Getenv("WEBSERVER_PORT")), nil)
}
