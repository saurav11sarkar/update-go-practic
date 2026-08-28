package main

import (
	"fmt"
	"go-prictic/internal/config"
	"go-prictic/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := config.ConnectionDB(cfg)
	if err != nil {
		panic(err)
	}

	if err := server.Start(cfg, db); err != nil {
		fmt.Println(err)
	}
}
