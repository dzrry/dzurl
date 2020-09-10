package main

import (
	"fmt"
	"github.com/dzrry/dzurl/domain"
	"github.com/dzrry/dzurl/repo/redis"
	"log"
)

func main() {
	r, err := redis.NewRepo("localhost", "6379", "")
	if err != nil {
		log.Fatal("13" + err.Error())
	}

	rdrct := &domain.Redirect{
		Key:       "12345",
		URL:       "123456789",
		CreatedAt: 1000000,
	}
	err = r.Store(rdrct)
	if err != nil {
		log.Fatal("23" + err.Error())
	}

	rdrct, err = r.Load("12345")
	if err != nil {
		log.Fatal("28" + err.Error())
	}
	fmt.Println("30" + rdrct.URL)

	rdrct.URL = "0123456789"
	err = r.Store(rdrct)
	if err != nil {
		log.Fatal("35" + err.Error())
	}
	fmt.Println(rdrct.URL)
	newRdrct, err := r.Load("12345")
	if err != nil {
		log.Fatal("40" + err.Error())
	}
	fmt.Println("42" + newRdrct.URL)
}