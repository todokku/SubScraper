package main

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
)

func setKey(conn redis.Conn, key string, value string) error {

	fmt.Print("Setting key value")

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

func getValue(conn redis.Conn, key string) (string, error) {

	fmt.Print("Getting value")

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return value, nil
}

func connectToRedis(redisServer string) (bool, redis.Conn) {

	conn, err := redis.Dial("tcp", redisServer)

	if err != nil {

		return false, conn
	}

	return true, conn

}

func main() {

	redisServer := os.Getenv("REDIS_SERVER")
	if redisServer != "" {
		status, conn := connectToRedis(redisServer)
		defer conn.Close()
		if status == true {
			fmt.Print("Success Redis")

			//err := setKey(conn, "day1", "SUN")
			val, err := getValue(conn, "zealous")
			if err != nil {

				fmt.Print("Failed in set")
			} else {

				//fmt.Print("Key Set")
				fmt.Print("Key Fetched , value : ", val)
			}

		}
	}

}
