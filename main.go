package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	connStr := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/addpayment", add).Methods("GET")
	r.HandleFunc("/removepayment", remove).Methods("DELETE")
	r.HandleFunc("/getpayments", get).Methods("GET")

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

func add(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println(err)
		http.Error(w, "Can't read body.", http.StatusBadRequest)
		return
	}

	fmt.Println(body)

	if _, err := url.ParseRequestURI(body["url"]); err != nil {
		log.Println(err)
		http.Error(w, "Invalid URL.", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO paymethods(url) VALUES($1);"
	if _, err := db.Query(query, body["url"]); err != nil {
		log.Println(err)
		http.Error(w, "We coundn't update the database at the moment.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}

func remove(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	if _, err := url.ParseRequestURI(vars.Get("url")); err != nil {
		log.Println(err)
		http.Error(w, "Invalid URL.", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM paymethods WHERE url = $1"
	if _, err := db.Query(query, vars.Get("url")); err != nil {
		log.Println(err)
		http.Error(w, "We weren't able to remove the given URL.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}

func get(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM paymethods;"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		http.Error(w, "We coudn't retrieve the URLs.", http.StatusInternalServerError)
		return
	}

	urls, err := readRows(rows)
	if err != nil {
		log.Println(err)
		http.Error(w, "We coudn't parse the URLs.", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(urls); err != nil {
		log.Println(err)
		http.Error(w, "We weren't able to encode the URLs.", http.StatusInternalServerError)
		return
	}
}

func readRows(rows *sql.Rows) ([]map[string]string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	rawResult := make([][]byte, len(cols))
	var ret []map[string]string
	dest := make([]interface{}, len(cols))

	// load the interface with pointers to get the data
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		// read data into dest that hold the pointers
		err := rows.Scan(dest...)
		if err != nil {
			return nil, err
		}

		result := make(map[string]string)
		// read the pointers and add them to the result
		for i, raw := range rawResult {
			result[cols[i]] = string(raw)
		}

		ret = append(ret, result)
	}
	return ret, nil
}
