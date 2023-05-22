package main

import (
	"Assignment/pkg/store/postgres"
	"fmt"
)

func main() {
	db, err := postgres.ConnectDB("localhost", "5432", "postgres", "postgres", "MIcroservices")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database")
	}
	defer db.Close()
}
