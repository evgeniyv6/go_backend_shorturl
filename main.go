package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go_backend_shorturl/configuration"
	"log"
	"math"
)

var FS configuration.OsFileSystem

func main() {
	fmt.Println("hello from custom bitly")
	fmt.Println(uint64(math.MaxUint64))

	// try redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Printf("Connect redis err: %s", err)
	}
	log.Println("Success redis connect!")
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Close redis err: %s\n", err)
		}
	}()

	rep, err := conn.Do("ping")
	if err != nil {
		log.Printf("Cmd DO err: %s", err)
	}
	log.Println(rep)

	log.Fatal(configuration.ReadConfig("config.json", FS))
}
