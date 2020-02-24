package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func GetMeaning(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]
	req, err := http.NewRequest("GET", "https://owlbot.info/api/v4/dictionary/" + word, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Set("Authorization", "Token 23db56a1aace2a60d425f5456265e8d988819728")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bodyString)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/meaning/{word}", GetMeaning).Methods("GET")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)) //Launch the app, visit localhost:8000/api

	if err != nil {
		fmt.Print(err)
	}
}