package prad

import (
	"log"

	"github.com/projectdiscovery/goflags"
)

type Options struct {
	Target      string
	WordFile    string
	Extension   string
	Concurrent  int
	Prefix      string
	Suffix      string
	Proxy       string
	FilterCode  int
	ExcludeCode int
	Timeout     int
}

func ParseOptions() *Options {
	o := &Options{}
	flags := goflags.NewFlagSet()
	flags.SetDescription("web directory and file discovery.")

	flags.SetGroup("input", "input options")
	flags.StringVarP(&o.Target, "url", "u", "", "url to scan").Group("input")
	flags.StringVarP(&o.WordFile, "word-file", "wf", "", "wordlist file").Group("input")

	flags.SetGroup("word", "word options")
	flags.StringVarP(&o.Extension, "word-ext", "we", "", "word extension").Group("word")
	flags.StringVarP(&o.Prefix, "word-prefix", "wp", "", "word prefix").Group("word")
	flags.StringVarP(&o.Suffix, "word-suffix", "ws", "", "word suffix").Group("word")

	flags.SetGroup("output", "output options")
	flags.IntVarP(&o.FilterCode, "filter-code", "fc", 0, "filter by status code").Group("output")
	flags.IntVarP(&o.ExcludeCode, "exclude-code", "ec", 0, "exclude by status code").Group("output")

	flags.IntVar(&o.Concurrent, "concurrent", 10, "concurrent goroutines")
	flags.StringVar(&o.Proxy, "proxy", "", "proxy")
	flags.IntVar(&o.Timeout, "timeout", 5, "timeout")

	showBanner()
	err := flags.Parse()
	if err != nil {
		log.Fatalf("parse options failed: %s", err)
	}

	return o
}