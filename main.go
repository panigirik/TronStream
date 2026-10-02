package main

import (
	"TronStream/internal/graphs"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", graphs.HealthHandler)

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Error starting HTTP server:", err)
		}
	}()

	select {}

}
