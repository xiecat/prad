package main

import (
	"context"
	"log"
	"strconv"

	"github.com/tardc/prad"
	"github.com/tardc/prad/internal/output"
	"github.com/tardc/prad/pkg/interrupt"
)

func main() {

	options := parseOptions()

	client, err := newClient(options)
	if err != nil {
		log.Fatalf("create client failed: %s\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go interrupt.HandleInterrupt(cancel)

	resultChan, err := client.Do(ctx, options.Target)
	if err != nil {
		log.Fatalf("run failed: %s\n", err)
	}

	var w output.Writer
	if options.OutputFile != "" {
		w = output.NewMultiOut(options.NoColor, options.OutputFile)
	} else {
		w = output.NewStdout(options.NoColor)
	}
	defer w.Close()
	for r := range resultChan {
		if options.filterStatusCode != nil {
			for _, statusCode := range options.filterStatusCode {
				if statusCode == strconv.Itoa(r.Code) {
					w.Write(r)
					break
				}
			}
		} else if options.excludeStatusCode != nil {
			var shouldOutput = true
			for _, statusCode := range options.excludeStatusCode {
				if statusCode == strconv.Itoa(r.Code) {
					shouldOutput = false
					break
				}
			}
			if shouldOutput {
				w.Write(r)
			}
		} else {
			w.Write(r)
		}
	}
}

func newClient(o *options) (*prad.Client, error) {
	client, err := prad.NewClient(o.Wordlist)
	if err != nil {
		return nil, err
	}

	if o.Proxy != "" {
		err := client.SetProxy(o.Proxy)
		if err != nil {
			return nil, err
		}
	}
	err = client.SetTimeout(o.Timeout)
	if err != nil {
		return nil, err
	}
	err = client.SetQPS(o.QPS)
	if err != nil {
		return nil, err
	}
	err = client.SetConcurrent(o.Concurrent)
	if err != nil {
		return nil, err
	}

	return client, err
}
