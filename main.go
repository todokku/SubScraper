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

type wordMap map[string]interface{}

func checkLocally(word string) (bool, string) {

	redisServer := os.Getenv("REDIS_SERVER")
	if redisServer != "" {
		status, conn := ConnectToRedis(redisServer)

		if status == true {
			defer conn.Close()
			fmt.Println("Success Redis")
			meaning, status := GetValue(conn, word)
			fmt.Println("Inside checkLocal ", meaning, status)
			if status != true {

				fmt.Println("Meaning not found in redis ")
				return false, ""
			}

			return true, meaning

		}
	}
	return false, ""
}

func getMeaning(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]
	fmt.Println("Fetching meaning")
	status, meaning := checkLocally(word)
	//var meaning string
	//status := false
	if status == false {
		fmt.Println("not found in DB")
		req, err := http.NewRequest("GET", "https://api.dictionaryapi.dev/api/v1/entries/en/"+word, nil)

		if err != nil {
			log.Fatal(err.Error())
			meaning = "MEANINGNOTFOUND"
		} else {
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err.Error())
			}

			if resp.StatusCode != 200 {
				meaning = "MEANINGNOTFOUND"
			} else {
				defer resp.Body.Close()

				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				var dat []wordMap

				_ = json.Unmarshal(bodyBytes, &dat)
				for k, v := range dat[0] {
					if k == "meaning" {
						for _, v1 := range v.(map[string]interface{}) {
							meaning = (v1.([]interface{})[0].(map[string]interface{})["definition"]).(string)
							fmt.Println("WORD: ", word, "MEANING: ", meaning)
							break
						}
					}
				}
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(meaning)
		fmt.Println("Written to response writer")

		if err == nil {
			fmt.Println("writing to DB")
			redisServer := os.Getenv("REDIS_SERVER")
			status, conn := ConnectToRedis(redisServer)

			if status == true {
				defer conn.Close()
				fmt.Println("New word added to redis dictionary")
				SetKey(conn, word, meaning)
			} else {
				fmt.Println("Failed to add the word to dictionary")
			}
		}

	} else {
		fmt.Println("Fetched from redis dictionary  ", word)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(meaning)
		err := json.NewEncoder(w).Encode(string(meaning))
		fmt.Println(err)
	}

}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/meaning/{word}", getMeaning).Methods("GET")

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
