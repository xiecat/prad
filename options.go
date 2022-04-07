package prad

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/projectdiscovery/goflags"
)

type Options struct {
	Target      string
	WordFile    string
	Wordlist    goflags.CommaSeparatedStringSlice
	Extension   string
	Concurrent  int
	Prefix      string
	Suffix      string
	Proxy       string
	FilterCode  int
	ExcludeCode int
	Timeout     int
	NoColor     bool
	QPS         int
	BasicAuth   string
	UserAgent   string
	Headers     goflags.CommaSeparatedStringSlice
	Resume      bool
	Offset      int
}

func ParseOptions() *Options {
	o := &Options{}
	flags := goflags.NewFlagSet()
	flags.SetDescription("web directory and file discovery.")

	flags.SetGroup("input", "input options")
	flags.StringVarP(&o.Target, "url", "u", "", "url to scan").Group("input")
	flags.StringVarP(&o.WordFile, "word-file", "wf", "", "wordlist file").Group("input")
	flags.CommaSeparatedStringSliceVarP(&o.Wordlist, "word-list", "wl", []string{}, "wordlist").Group("input")
	flags.BoolVar(&o.Resume, "resume", false, "resume task from resume.cfg").Group("input")

	flags.SetGroup("word", "word options")
	flags.StringVarP(&o.Extension, "word-ext", "we", "", "word extension").Group("word")
	flags.StringVarP(&o.Prefix, "word-prefix", "wp", "", "word prefix").Group("word")
	flags.StringVarP(&o.Suffix, "word-suffix", "ws", "", "word suffix").Group("word")

	flags.SetGroup("output", "output options")
	flags.IntVarP(&o.FilterCode, "filter-code", "fc", 0, "filter by status code").Group("output")
	flags.IntVarP(&o.ExcludeCode, "exclude-code", "ec", 0, "exclude by status code").Group("output")
	flags.BoolVarP(&o.NoColor, "no-color", "nc", false, "disable color in output")

	flags.IntVar(&o.Concurrent, "concurrent", 10, "concurrent goroutines")
	flags.StringVar(&o.Proxy, "proxy", "", "proxy")
	flags.IntVar(&o.Timeout, "timeout", 5, "timeout")
	flags.IntVar(&o.QPS, "qps", 10, "QPS")
	flags.StringVar(&o.BasicAuth, "basic-auth", "", "basic auth user:pass")
	flags.StringVar(&o.UserAgent, "user-agent", "", "user agent")
	flags.CommaSeparatedStringSliceVar(&o.Headers, "headers", []string{}, "custom headers")

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

	if o.Resume {
		err = o.ReadConfigFile("")
		if err != nil {
			log.Fatalf("read config file failed: %s", err)
		}
	}

	err = o.CheckOptions()
	if err != nil {
		log.Fatalf("check options failed: %s", err)
	}

	return o
}

func (o *Options) CheckOptions() error {
	if o.BasicAuth != "" && !strings.Contains(o.BasicAuth, ":") {
		return fmt.Errorf("incorrect basic auth format: %s", o.BasicAuth)
	}

	for _, header := range o.Headers {
		if !strings.Contains(header, ":") {
			return fmt.Errorf("incorrect custom header: %s", header)
		}
	}

	return nil
}

func (o *Options) ReadConfigFile(filename string) error {
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

func (o *Options) ReadWordFile(filename string) error {
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
		fr, err = Fs.Open(path.Join("wordlist", "common.txt"))
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
