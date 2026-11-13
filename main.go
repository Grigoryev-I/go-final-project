package main

import (
	"fmt"
	"log"
	"planner/pkg/db"
	"planner/pkg/server"
)

func main() {
	if err := db.Init(db.GetDBFile()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Запускаем сервер")
	server.Start()
}
