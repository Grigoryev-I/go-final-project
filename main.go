package main

import (
	"log"
	"planner/pkg/db"
	"planner/pkg/server"
)

func main() {
	if err := db.Init(db.GetDBFile()); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
