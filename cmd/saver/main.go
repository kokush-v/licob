package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cam-per/licob"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ds, err := licob.Auth()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Authorized as", ds.State.User.String())

	saver := NewSaver(ds)
	saver.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
