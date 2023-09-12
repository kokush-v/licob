package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cam-per/licob"
	"github.com/cam-per/licob/config"
	"github.com/cam-per/licob/utils"
	"github.com/cam-per/licob/web"
)

func main() {
	config.Lookup()

	s, err := licob.Auth()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Authorized as", s.State.User.String())

	fmt.Println("\nWeb server addr:")
	for _, addr := range utils.LookupIP() {
		fmt.Printf("http://%s:%s/channels/389413733329272837\n", addr, config.Port)
	}

	web.Init(s)
	web.RunAsync()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
