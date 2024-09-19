package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhilashdk2016/order-poller/storage"
	"github.com/joho/godotenv"
)

func pollAndUpdateEventsInOutbox(r *storage.Repository) {
	err := r.GetOutbox()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Could not connect to database")
	}
	fmt.Println("Connected to transactional_pattern_demo")
	r := storage.Repository{
		DB: db,
	}

	// utils.NewTimedExecutor(5*time.Second, 10*time.Second).Start(func() {
	// 	pollAndUpdateEventsInOutbox(&r)
	// }, true)

	done := make(chan int)
	ticker := time.NewTicker(50 * time.Second)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-ticker.C:
				pollAndUpdateEventsInOutbox(&r)
			case <-sigs:
				ticker.Stop()
				done <- 1
				return
			}
		}
	}()
	<-done
}
