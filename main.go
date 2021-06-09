package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go_backend_shorturl/configuration"
	"go_backend_shorturl/hasher"
	"log"
	"math"
)

var FSO configuration.OsFileSystem

func main() {
	fmt.Println("hello from custom bitly. Clear this file when complete.")
	fmt.Println(uint64(math.MaxUint64))

	log.Println(string(hasher.GenHash(uint64(math.MaxUint64))))
	log.Println(hasher.GenClear("pIrkgbKrQ8v"))

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

	c, err := configuration.ReadConfig("config.json", FSO)
	log.Println(c)

}
