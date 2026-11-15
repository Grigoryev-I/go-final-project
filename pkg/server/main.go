package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"planner/pkg/api"
	"strconv"
)

const defaultPort = "7540"
const webDir = "./web"

func Start() error {
	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":" + getPort()
	fmt.Printf("Сервер запущен: http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}

func getPort() string {
	port := defaultPort

	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if _, err := strconv.Atoi(envPort); err == nil {
			port = envPort
			log.Printf("Использую переменную окружения TODO_PORT=%s", port)
		} else {
			log.Printf("Некорректное значение TODO_PORT = %q, использую %s", envPort, port)
		}
	}

	return port
}
