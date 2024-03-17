package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jatinkumar-me/go-lru/lru"
)

type CacheEntry struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

func routeHandler() *http.ServeMux {
	router := http.NewServeMux()

	// Initialize the LRU cache
	cache := lru.NewLRU[int, string](1<<10, time.Second*5)

	router.HandleFunc("/cache/set", func(w http.ResponseWriter, r *http.Request) {
		var data CacheEntry
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		cache.Put(data.Key, data.Value)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Value set in cache for key %d\n", data.Key)
	})

	router.HandleFunc("/cache/get", func(w http.ResponseWriter, r *http.Request) {
		keyStr := r.URL.Query().Get("key")
		key, err := strconv.Atoi(keyStr)
		if err != nil {
			http.Error(w, "Invalid key", http.StatusBadRequest)
			return
		}

		value, ok := cache.Get(key)
		if !ok {
			http.Error(w, "Key not found in cache", http.StatusNotFound)
			return
		}

		// Respond with value
		response := struct {
			Key   int    `json:"key"`
			Value string `json:"value"`
		}{Key: key, Value: value}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	router.HandleFunc("/log", logHandler)
	return router
}
