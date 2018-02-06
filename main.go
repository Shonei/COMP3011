package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	temp, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// Routes in place for testing purposes
	r.Handle("/ping", test())

	// listen on the router
	http.Handle("/", r)

	log.Println(http.ListenAndServe(":"+port, nil))
}

func test() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
}
