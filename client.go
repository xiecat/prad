package prad

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/logrusorgru/aurora"
	"golang.org/x/time/rate"
)

type Client struct {
	Wordlist      []string
	Client        *http.Client
	Options       *Options
	ResultHandler func(*Result)
	RateLimiter   *rate.Limiter
}

func NewClient(options *Options) (*Client, error) {
	ht := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	if options.Proxy != "" {
		u, err := url.Parse(options.Proxy)
		if err != nil {
			return nil, err
		}
		ht.Proxy = http.ProxyURL(u)
	}
	hc := &http.Client{
		Timeout:   time.Second * time.Duration(options.Timeout),
		Transport: ht,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var (
		fr       io.ReadCloser
		err      error
		wordlist []string
	)
	if options.WordFile != "" {
		fr, err = os.Open(options.WordFile)
	} else {
		fr, err = Fs.Open(path.Join("wordlist", "common.txt"))
	}
	if err != nil {
		return nil, fmt.Errorf("open wordlist file failed: %s", err)
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

	c := &Client{
		Wordlist:    wordlist,
		Client:      hc,
		Options:     options,
		RateLimiter: rate.NewLimiter(rate.Limit(options.QPS), 1),
	}

	c.ResultHandler = func(r *Result) {
		var (
			output string
			code   string
		)
		if r.Redirect != "" {
			output = fmt.Sprintf("%s -> %s", r.URL, r.Redirect)
		} else {
			output = r.URL
		}

		if options.NoColor {
			code = strconv.Itoa(r.Code)
		} else {
			switch r.Code {
			case http.StatusNotFound:
				code = aurora.BrightRed(r.Code).String()
			case http.StatusOK:
				code = aurora.BrightGreen(r.Code).String()
			default:
				code = aurora.BrightYellow(r.Code).String()
			}
		}

		fmt.Printf("%s - %s\n", code, output)
	}

	return c, nil
}

func (c *Client) Do(ctx context.Context, target string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	wordChan := make(chan string, c.Options.Concurrent)
	go func() {
		defer close(wordChan)

		for _, word := range c.Wordlist[c.Options.Offset:] {
			select {
			case <-ctx.Done():
				return
			default:
			}

			wordChan <- word
		}
	}()

	var counter uint32
	wg := &sync.WaitGroup{}
	for i := 0; i < c.Options.Concurrent && i < len(c.Wordlist); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				err := c.RateLimiter.Wait(ctx)
				if err != nil {
					log.Printf("rate limiter failed when wait: %s\n", err)
				}

				word, ok := <-wordChan
				if !ok {
					break
				}
				resp, err := c.Check(ctx, target, word)
				atomic.AddUint32(&counter, 1)
				if err != nil {
					log.Printf("check %s %s failed: %s\n", target, word, err)
					continue
				}

				if c.Options.FilterCode != 0 && resp.Code != c.Options.FilterCode {
					continue
				}
				if c.Options.ExcludeCode != 0 && resp.Code == c.Options.ExcludeCode {
					continue
				}

				if resp != nil && c.ResultHandler != nil {
					c.ResultHandler(resp)
				}
			}
		}()
	}

	wg.Wait()

	if int(counter) != len(c.Wordlist) {
		c.Options.Offset = int(counter)

		fd, err := os.Create("resume.cfg")
		if err != nil {
			return fmt.Errorf("create resume file failed: %s", err)
		}

		err = json.NewEncoder(fd).Encode(c.Options)
		if err != nil {
			return fmt.Errorf("save resume file failed: %s", err)
		}
	}

	return nil
}

func (c *Client) Check(ctx context.Context, target, word string) (*Result, error) {
	var newWord string
	if c.Options.Extension != "" {
		newWord = fmt.Sprintf("%s%s%s.%s", c.Options.Prefix, word, c.Options.Suffix, c.Options.Extension)
	} else {
		newWord = fmt.Sprintf("%s%s%s", c.Options.Prefix, word, c.Options.Suffix)
	}

	var u string
	if strings.Contains(target, "{{") {
		reg := regexp.MustCompile(`{{.*?}}`)
		u = reg.ReplaceAllString(target, newWord)
	} else {
		u = fmt.Sprintf("%s/%s",
			strings.TrimSuffix(target, "/"),
			strings.TrimPrefix(newWord, "/"),
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	if c.Options.BasicAuth != "" {
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.Options.BasicAuth)))
	}

	if c.Options.UserAgent != "" {
		req.Header.Set("User-Agent", c.Options.UserAgent)
	}

	if len(c.Options.Headers) > 0 {
		for _, header := range c.Options.Headers {
			kvs := strings.SplitN(header, ":", 2)
			req.Header.Set(kvs[0], kvs[1])
		}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	result := &Result{
		URL:      u,
		Code:     resp.StatusCode,
		Redirect: resp.Header.Get("Location"),
	}

	return result, nil
}
