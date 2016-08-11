package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

/*
// TODO
func deployment(descriptor []byte, ws *websocket.Conn) {
	for i := 0; i < len(descriptor); i++ {
		time.Sleep(100 * time.Millisecond)
	}
}
*/

func deploy(w http.ResponseWriter, r *http.Request) {
	log.Println("Entered deploy")

	// key websocket call, not entirely sure what this does yet.
	// It might be setting up the reader and writer for the websocket.
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	// defer means do this AFTER this deploy method returns
	defer ws.Close()

	// run forever until we hit a break
	for {
		log.Println("top of for loop")

		// Read from the websocket
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// We could call a goroutine here to launch a "deployment"
		// We could have a channel feeding us information from the deployment
		// goroutine to write status out to the websocket

		// write back to the websocket
		err = ws.WriteMessage(websocket.TextMessage, []byte("We got it"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println("Entered serveHome")

	// verify we didn't get here because of a different url
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}

	// We only support GET for now, but could support other
	// methods in the future
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
	}

	// Tell the caller we're returning html
	w.Header().Set("Content-Type", "text/html")

	// load the file and write it
	template.Must(template.ParseFiles("deployment.html")).Execute(w, r.Host)
}

func main() {
	// get the address from the command line
	var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
	flag.Parse()

	// map urls to methods
	http.HandleFunc("/ws", deploy)
	http.HandleFunc("/", serveHome)

	// listen and wait on the given address
	log.Fatal(http.ListenAndServe(*addr, nil))
}
