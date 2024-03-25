package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

var payloads map[int][]byte
var once sync.Once

func main() {
	generatePayloads()

	http.HandleFunc("/secondarycache/", func(w http.ResponseWriter, r *http.Request) {
		PublicCacheHandler(w, r)
	})

	http.HandleFunc("/nocache", func(w http.ResponseWriter, r *http.Request) {
		NoCacheHandler(w, r)
	})

	http.HandleFunc("/privatecache", func(w http.ResponseWriter, r *http.Request) {
		PrivateCacheHandler(w, r)
	})

	http.HandleFunc("/getresponse", func(w http.ResponseWriter, r *http.Request) {
		getresponseWithoutHeaders(w, r)
	})

	http.HandleFunc("/queryresource", func(w http.ResponseWriter, r *http.Request) {
		QueryResourceHandler(w, r)
	})

	// Start server
	fmt.Println("Server is listening on port 8082...")
	http.ListenAndServe(":8082", nil)
}

func QueryResourceHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	param1 := queryParams.Get("param1")
	param2 := queryParams.Get("param2")
	param3 := queryParams.Get("param3")

	if param1 == "" || param2 == "" || param3 == "" {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("param1: %s, param2: %s, param3: %s", param1, param2, param3)

	// Write response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}


func PublicCacheHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/secondarycache/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	payload, ok := payloads[id]
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=60")
	w.Write(payload)
	additionalContent := []byte(" cached as public for ID " + strconv.Itoa(id))
	w.Write(additionalContent)
}

func PrivateCacheHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "private, max-age=3600")
   // w.Write(payload)
    // sleepBeforeRespond()
}

func getresponseWithoutHeaders(w http.ResponseWriter, r *http.Request) {
   // w.Write(payload)
    // sleepBeforeRespond()
}

func NoCacheHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "no-store")
   // w.Write(payload)
    // sleepBeforeRespond()
}

func generatePayloads() {
	once.Do(func() {
		payloads = make(map[int][]byte)
		for i := 1; i <= 50; i++ {
			payload := make([]byte, 102400)
			for j := range payload {
				payload[j] = 'y'
			}
			payloads[i] = payload
		}
	})
}
