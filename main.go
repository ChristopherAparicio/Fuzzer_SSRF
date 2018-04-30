package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var uuidStorage map[string]bool
var mutex *sync.Mutex

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// http:host/uuid
	uuid := r.URL.Path[1:]
	if uuid != "" {
		mutex.Lock()
		uuidStorage[uuid] = true
		mutex.Unlock()
	}
	fmt.Fprintf(w, "Hello World !")
	return
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid != "" {
		mutex.Lock()
		uuidIsPresent := uuidStorage[uuid]
		mutex.Unlock()
		if uuidIsPresent {
			data := map[string]interface{}{
				"uuid":    uuid,
				"present": true,
			}
			jsonResponse(w, 200, data)
			return
		}
		data := map[string]interface{}{
			"uuid":    uuid,
			"present": false,
		}
		jsonResponse(w, 200, data)
		return
	}
	jsonResponse(w, 400, nil)
	return
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	// Init UUID Storage and Mutex
	uuidStorage = make(map[string]bool)
	mutex = &sync.Mutex{}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/read", readHandler)
	log.Printf("Server up and running on port 9100")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func jsonResponse(w http.ResponseWriter, code int, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
