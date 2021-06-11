package main

import (
	"fmt"
	"log"
	"math"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/hasher"
)

func main() {
	fmt.Println("hello from custom bitly. Clear this file when complete.")
	fmt.Println(uint64(math.MaxUint64))

	log.Println(string(hasher.GenHash(uint64(math.MaxUint64))))
	log.Println(hasher.GenClear("pIrkgbKrQ8v"))

}
