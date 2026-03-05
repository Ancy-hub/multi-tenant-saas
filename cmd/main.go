package main

import (
	"log"
	"net/http"
	"github.com/ancy-shibu/multi-tenant-saas/internal/db"
	"github.com/go-chi/chi/v5"
)

func main() {
	database,err:=db.NewDB();
	if err!=nil{
		log.Fatal("Failed to connect to db:",err)
	}
	defer database.Close()
	log.Println("Database connected")

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("Server running on port 8081")

	err= http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatal(err)
	}
}