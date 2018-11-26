package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"twitter/framework"
)

func main(){
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func(){
		sig := <- sigs
		log.Printf("Receive program end signal-------------: ", sig)
		done <- true
	}()


	for i := 0; i < 1; i++{
		go framework.Crawl()
		go framework.Parse()
		go framework.Store()
	}

	<- done
	log.Printf("Program exit-------------: ")


}
