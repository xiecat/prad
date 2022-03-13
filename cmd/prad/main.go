package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/tardc/prad"
)

func main() {
	options := prad.ParseOptions()

	client, err := prad.NewClient(options)
	if err != nil {
		log.Fatalf("create client failed: %s\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer close(signalChan)
	go func() {
		count := 0
		for {
			_, ok := <-signalChan
			if ok {
				if count == 0 {
					count += 1
					cancel()
				} else {
					os.Exit(1)
				}
			} else {
				return
			}
		}
	}()

	err = client.Do(ctx, options.Target)
	if err != nil {
		log.Fatalf("run failed: %s\n", err)
	}
}
