package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/maintenance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req MaintenanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := CaptureMaintenance(client, req); err != nil {
			if ie, ok := err.(*InfraiError); ok && ie.Status >= 400 && ie.Status < 500 {
				http.Error(w, ie.Message, ie.Status)
				return
			}
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"group": ErrorGroup("maintenance", req.ID)})
	})
	log.Println("property error service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
