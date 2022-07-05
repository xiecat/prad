package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/xiecat/prad"
	"github.com/xiecat/prad/internal/output"
	"github.com/xiecat/prad/pkg/interrupt"
)

func main() {

	options := parseOptions()

	client, err := newClient(options)
	if err != nil {
		gologger.Fatal().Msgf("create client failed: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go interrupt.HandleInterrupt(cancel)

	resultChan, err := client.Do(ctx, options.Target)
	if err != nil {
		gologger.Fatal().Msgf("run failed: %s", err)
	}

	var w output.Writer
	if options.OutputFile != "" {
		w = output.NewMultiOut(options.NoColor, options.OutputFile)
	} else {
		w = output.NewStdout(options.NoColor)
	}
	defer w.Close()

	for r := range resultChan {
		if r.Error != nil {
			gologger.Debug().Msgf("check failed: %s", r.Error)
			continue
		}

		if options.FilterStatusCode != nil {
			for _, statusCode := range options.FilterStatusCode {
				if statusCode == strconv.Itoa(r.Result.Code) {
					w.Write(r.Result)
					break
				}
			}
		} else if options.ExcludeStatusCode != nil {
			var shouldOutput = true
			for _, statusCode := range options.ExcludeStatusCode {
				if statusCode == strconv.Itoa(r.Result.Code) {
					shouldOutput = false
					break
				}
			}
			if shouldOutput {
				w.Write(r.Result)
			}
		} else {
			w.Write(r.Result)
		}

		options.ProcessedNum++
	}

	if options.ProcessedNum != len(options.Wordlist) {
		if options.ResumeFile == "" {
			rand.Seed(time.Now().Unix())
			options.ResumeFile = fmt.Sprintf("resume-%d.cfg", rand.Int())
		}

		err = options.WriteConfigFile(options.ResumeFile)
		if err != nil {
			gologger.Fatal().Msgf("read wordlist file failed: %s", err)
		}
	} else {
		if options.ResumeFile != "" {
			err = os.Remove(options.ResumeFile)
			if err != nil {
				gologger.Fatal().Msgf("remove resume file failed: %s", err)
			}
		}
	}
}

func newClient(o *options) (*prad.Client, error) {
	client, err := prad.NewClient(o.Wordlist[o.ProcessedNum:])
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
