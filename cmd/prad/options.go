package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/xiecat/prad/assets"
)

type options struct {
	Target            string
	Wordlist          goflags.CommaSeparatedStringSlice
	OutputFile        string
	Concurrent        int
	Proxy             string
	Timeout           int
	NoColor           bool
	QPS               int
	FilterStatusCode  goflags.CommaSeparatedStringSlice
	ExcludeStatusCode goflags.CommaSeparatedStringSlice

	ResumeFile   string
	ProcessedNum int
}

func parseOptions() *options {
	var (
		wordFile string
		verbose  bool
	)
	o := &options{}
	flags := goflags.NewFlagSet()
	flags.SetDescription("web directory and file discovery.")

	flags.SetGroup("input", "input options")
	flags.StringVarP(&o.Target, "url", "u", "", "url to scan").Group("input")
	flags.StringVarP(&wordFile, "word-file", "wf", "", "wordlist file").Group("input")
	flags.CommaSeparatedStringSliceVarP(&o.Wordlist, "word-list", "wl", []string{}, "wordlist").Group("input")

	flags.SetGroup("output", "output options")
	flags.BoolVarP(&o.NoColor, "no-color", "nc", false, "disable color in output").Group("output")
	flags.StringVarP(&o.OutputFile, "output-file", "of", "", "output filename").Group("output")
	flags.CommaSeparatedStringSliceVarP(&o.FilterStatusCode, "filter-status", "fs", []string{}, "filtering using status codes").Group("output")
	flags.CommaSeparatedStringSliceVarP(&o.ExcludeStatusCode, "exclude-status", "es", []string{}, "excluding using status codes").Group("output")

	flags.IntVar(&o.Concurrent, "concurrent", 10, "concurrent goroutines")
	flags.StringVar(&o.Proxy, "proxy", "", "proxy")
	flags.IntVar(&o.Timeout, "timeout", 5, "timeout")
	flags.IntVar(&o.QPS, "qps", 10, "QPS")
	flags.StringVar(&o.ResumeFile, "resume", "", "resume from config file")
	flags.BoolVarP(&verbose, "verbose", "V", false, "verbose")

	showBanner()
	err := flags.Parse()
	if err != nil {
		gologger.Fatal().Msgf("parse options failed: %s", err)
	}

	if flags.CommandLine.NFlag() < 1 {
		flags.CommandLine.Usage()
		os.Exit(1)
	}

	gologger.DefaultLogger.SetFormatter(formatter.NewCLI(o.NoColor))
	if verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	}

	if o.ResumeFile != "" {
		err = o.ReadConfigFile(o.ResumeFile)
		if err != nil {
			gologger.Fatal().Msgf("resume failed from %s: %s", o.ResumeFile, err)
		}
	}

	if o.Target == "" {
		gologger.Fatal().Msg("target must be set")
	} else {
		if matched, err := regexp.MatchString(`(?i)^https?:\/\/`, o.Target); !matched || err != nil {
			gologger.Fatal().Msg("unsupported protocol scheme")
		}
	}

	if o.Wordlist == nil {
		err = o.ReadWordFile(wordFile)
		if err != nil {
			gologger.Fatal().Msgf("read wordlist file failed: %s", err)
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
	defer fd.Close()

	err = json.NewDecoder(fd).Decode(o)
	if err != nil {
		return fmt.Errorf("read resume file failed: %s", err)
	}

	return nil
}

func (o *options) WriteConfigFile(filename string) error {
	var (
		fd  *os.File
		err error
	)
	if filename != "" {
		fd, err = os.Create(filename)
	} else {
		fd, err = os.Open("resume.cfg")
	}
	if err != nil {
		return fmt.Errorf("open resume file failed: %s", err)
	}
	defer fd.Close()

	err = json.NewEncoder(fd).Encode(o)
	if err != nil {
		return fmt.Errorf("write resume file failed: %s", err)
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
		gologger.Warning().Msgf("close wordlist file failed: %s", err)
	}
	o.Wordlist = wordlist

	return nil
}
