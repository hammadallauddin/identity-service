package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"

// 	"honnef.co/go/tools/config"
// )

// func main() {
// 	cfg := config.Load()

// 	// Define route with JSON response
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"service": cfg.ServiceName,
// 			"status":  "running",
// 			"message": "DDD server with minimal setup",
// 		})
// 	})

// 	// Start server
// 	log.Printf("🚀 Server %s listening on :%s", cfg.ServiceName, cfg.Port)
// 	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
// }
