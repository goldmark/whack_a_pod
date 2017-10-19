package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler(),
	}

	log.Printf("starting whack a pod api api")
	srv.ListenAndServe()
}

func handler() http.Handler {

	r := http.NewServeMux()
	r.HandleFunc("/", health)
	r.HandleFunc("/healthz", health)
	r.HandleFunc("/api/healthz", health)
	r.HandleFunc("/api/color", color)
	r.HandleFunc("/api/color-complete", colorComplete)
	r.HandleFunc("/api/color/", color)
	r.HandleFunc("/api/color-complete/", colorComplete)
	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func color(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, hexColorString())
}

type result struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

func colorComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	h, _ := os.Hostname()
	result := result{
		Name:  h,
		Color: hexColorString(),
	}

	b, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("{\"error\":\"could not unmarshap data %v\"}", err)
		sendJSON(w, msg, http.StatusInternalServerError)
	}

	sendJSON(w, string(b), http.StatusOK)
}

func hexColorString() string {
	result := "#"
	for i := 1; i <= 3; i++ {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(256)
		result += strconv.FormatInt(int64(i), 16)
	}
	return result
}

func sendJSON(w http.ResponseWriter, content string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, content)
}
