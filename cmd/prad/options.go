package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/projectdiscovery/goflags"
	"github.com/tardc/prad/assets"
	"io"
	"log"
	"os"
	"path"
)

type options struct {
	Target     string
	WordFile   string
	Wordlist   goflags.CommaSeparatedStringSlice
	OutputFile string
	Concurrent int
	Proxy      string
	Timeout    int
	NoColor    bool
	QPS        int
}

func parseOptions() *options {
	o := &options{}
	flags := goflags.NewFlagSet()
	flags.SetDescription("web directory and file discovery.")

	flags.SetGroup("input", "input options")
	flags.StringVarP(&o.Target, "url", "u", "", "url to scan").Group("input")
	flags.StringVarP(&o.WordFile, "word-file", "wf", "", "wordlist file").Group("input")
	flags.CommaSeparatedStringSliceVarP(&o.Wordlist, "word-list", "wl", []string{}, "wordlist").Group("input")

	flags.SetGroup("output", "output options")
	flags.BoolVarP(&o.NoColor, "no-color", "nc", false, "disable color in output").Group("output")
	flags.StringVarP(&o.OutputFile, "output-file", "of", "", "output filename").Group("output")

	flags.IntVar(&o.Concurrent, "concurrent", 10, "concurrent goroutines")
	flags.StringVar(&o.Proxy, "proxy", "", "proxy")
	flags.IntVar(&o.Timeout, "timeout", 5, "timeout")
	flags.IntVar(&o.QPS, "qps", 10, "QPS")

	showBanner()
	err := flags.Parse()
	if err != nil {
		log.Fatalf("parse options failed: %s", err)
	}

	if o.Wordlist == nil {
		err = o.ReadWordFile("")
		if err != nil {
			log.Fatalf("read wordlist file failed: %s", err)
		}
	}

	return o
}

func (o *options) ReadConfigFile(filename string) error {
	var (
		fd  *os.File
		err error
	)
	if filename != "" {
		fd, err = os.Open(filename)
	} else {
		fd, err = os.Open("resume.cfg")
	}
	if err != nil {
		return fmt.Errorf("open resume file failed: %s", err)
	}

	err = json.NewDecoder(fd).Decode(o)
	if err != nil {
		return fmt.Errorf("read resume file failed: %s", err)
	}

	return nil
}

func (o *options) ReadWordFile(filename string) error {
	var (
		fr       io.ReadCloser
		err      error
		wordlist []string
	)

	if filename != "" {
		fr, err = os.Open(filename)
	} else if o.WordFile != "" {
		fr, err = os.Open(o.WordFile)
	} else {
		fr, err = assets.Fs.Open(path.Join("wordlist", "common.txt"))
	}
	if err != nil {
		return fmt.Errorf("open wordlist file failed: %s", err)
	}
	fs := bufio.NewScanner(fr)
	fs.Split(bufio.ScanLines)
	for fs.Scan() {
		wordlist = append(wordlist, fs.Text())
	}
	err = fr.Close()
	if err != nil {
		log.Printf("close wordlist file failed: %s\n", err)
	}
	o.Wordlist = wordlist

	return nil
}
