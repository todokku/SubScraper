package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

//SetKey setting meaning
func SetKey(conn redis.Conn, key string, value string) error {

	fmt.Print("Setting key value")

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

//GetValue getting meaning
func GetValue(conn redis.Conn, key string) (string, error) {

	fmt.Print("Getting value")

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return value, nil
}

func ConnectToRedis(redisServer string) (bool, redis.Conn) {

	conn, err := redis.Dial("tcp", redisServer)

	if err != nil {

		return false, conn
	}

	return true, conn

}

/*func Connect() {

	redisServer := os.Getenv("REDIS_SERVER")
	if redisServer != "" {
		status, conn := connectToRedis(redisServer)
		defer conn.Close()
		if status == true {
			fmt.Print("Success Redis")

			//err := setKey(conn, "day1", "SUN")
			val, err := getValue(conn, "discourse")
			if err != nil {

				fmt.Print("Failed in set")
			} else {

				//fmt.Print("Key Set")
				fmt.Print("Key Fetched , value : ", val)
			}

		}
	}

}*/
