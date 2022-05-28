package main

import (
	"context"
	"github.com/tardc/prad"
	"github.com/tardc/prad/pkg/interrupt"
	"log"
)

func main() {
	options := prad.ParseOptions()

	client, err := prad.NewClient(options)
	if err != nil {
		log.Fatalf("create client failed: %s\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go interrupt.HandleInterrupt(cancel)

	err = client.Do(ctx, options.Target)
	if err != nil {
		log.Fatalf("run failed: %s\n", err)
	}
}
