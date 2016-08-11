package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func deploy(w http.ResponseWriter, r *http.Request) {
	log.Println("Entered deploy")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	for {
		log.Println("top of for loop")
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = ws.WriteMessage(websocket.TextMessage, []byte("We got it"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println("Entered serveHome")
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
	}
	w.Header().Set("Content-Type", "text/html")
	template.Must(template.ParseFiles("deployment.html")).Execute(w, r.Host)
}

func main() {
	var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
	flag.Parse()
	http.HandleFunc("/ws", deploy)
	http.HandleFunc("/", serveHome)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
