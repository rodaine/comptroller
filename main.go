package main

import (
	"log"

	"sync"

	"github.com/rodaine/comptroller/config"
	"github.com/rodaine/comptroller/gh"
)

func main() {
	ctx, err := config.Init()
	checkErr(err)

	_, err = gh.Init(ctx)
	checkErr(err)

	evts, errs := gh.Ingest(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for evt := range evts {
			log.Print(evt)
		}
		wg.Done()
	}()

	go func() {
		for err := range errs {
			log.Print(err)
		}
		wg.Done()
	}()

	wg.Wait()
	log.Println("All good.")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
