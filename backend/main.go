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

	http.HandleFunc("/publiccache/", func(w http.ResponseWriter, r *http.Request) {
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

	// Start server
	fmt.Println("Server is listening on port 8081...")
	http.ListenAndServe(":8081", nil)
}

func PublicCacheHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/publiccache/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	payload, ok := payloads[id]
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=180")
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
				payload[j] = 'x'
			}
			payloads[i] = payload
		}
	})
}
