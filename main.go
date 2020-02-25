package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func checkLocally(word string) (bool, string) {

	redisServer := os.Getenv("REDIS_SERVER")
	if redisServer != "" {
		status, conn := ConnectToRedis(redisServer)

		if status == true {
			defer conn.Close()
			fmt.Println("Success Redis")
			meaning, err := GetValue(conn, word)
			if err != nil {

				fmt.Print("Failed in set")
				return false, ""
			} else {

				return true, meaning

			}

		}
	}
	return false, ""
}

func GetMeaning(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	word := vars["word"]
	fmt.Println("Fetching meaning")
	status, meaning := checkLocally(word)
	//var meaning string
	//status := false
	if status == false {
		fmt.Println("not found in DB")
		// req, err := http.NewRequest("GET", "https://owlbot.info/api/v4/dictionary/"+word, nil)
		// if err != nil {
		// log.Fatal(err.Error())
		// }
		// req.Header.Set("Authorization", "Token 23db56a1aace2a60d425f5456265e8d988819728")

		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// 	log.Fatal(err.Error())
		// }
		// defer resp.Body.Close()

		// bodyBytes, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// bodyString := string(bodyBytes)
		cmd := exec.Command("./main.py", word)

		meaning, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Output: ", string(meaning))
		bodyString := string(meaning)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(bodyString)
		fmt.Println("Written to response writer")
		if err == nil {
			fmt.Println("writing to DB")
			redisServer := os.Getenv("REDIS_SERVER")
			status, conn := ConnectToRedis(redisServer)

			if status == true {
				defer conn.Close()
				fmt.Println("New word added to redis dictionary")
				SetKey(conn, word, bodyString)
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
