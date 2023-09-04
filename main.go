package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
	"strings"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id := uuid.New().String()
			text, _ := io.ReadAll((r.Body))
			err := rdb.Set(ctx, id, text, 0).Err()

			if err != nil {
				fmt.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to write in db"))
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				data := make(map[string]string)
				data["id"] = id
				jsonResp, _ := json.Marshal(data)
				w.Write(jsonResp)
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Not implemented"))
		}
	})
	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/view/")
		if id != "" {
			val, err := rdb.Get(ctx, id).Result()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Failed to read db"))
			} else {
				fmt.Fprintf(w, val)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No id provided"))

		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
