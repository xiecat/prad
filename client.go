package prad

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Client struct {
	wordlist    []string
	concurrent  int
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

func NewClient(wordlist []string) (*Client, error) {
	ht := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	hc := &http.Client{
		Transport: ht,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	c := &Client{
		wordlist:   wordlist,
		concurrent: 10,
		httpClient: hc,
	}

	return c, nil
}

func (c *Client) Do(ctx context.Context, target string) (<-chan *Result, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	wordChan := make(chan string, c.concurrent)
	go func() {
		defer close(wordChan)

		for _, word := range c.wordlist {
			select {
			case <-ctx.Done():
				return
			default:
			}

			wordChan <- word
		}
	}()

	resultChan := make(chan *Result, c.concurrent)
	wg := &sync.WaitGroup{}
	for i := 0; i < c.concurrent && i < len(c.wordlist); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if c.rateLimiter != nil {
					err := c.rateLimiter.Wait(ctx)
					if err != nil {
						log.Debugf("rate limiter failed when wait: %s\n", err)
					}
				}

				word, ok := <-wordChan
				if !ok {
					break
				}
				resp, err := c.Check(ctx, target, word)
				if err != nil {
					log.Debugf("check %s %s failed: %s\n", target, word, err)
					continue
				}

				resultChan <- resp
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan, nil
}

func (c *Client) Check(ctx context.Context, target, word string) (*Result, error) {
	var u string
	if strings.Contains(target, "{{") {
		reg := regexp.MustCompile(`{{.*?}}`)
		u = reg.ReplaceAllString(target, word)
	} else {
		u = fmt.Sprintf("%s/%s",
			strings.TrimSuffix(target, "/"),
			strings.TrimPrefix(word, "/"),
		)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &Result{
		URL:      u,
		Code:     resp.StatusCode,
		Redirect: resp.Header.Get("Location"),
	}

	return result, nil
}

func (c *Client) SetProxy(proxy string) error {
	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			return err
		}
		c.httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(u)
	} else {
		return errors.New("empty string for proxy")
	}

	return nil
}

func (c *Client) SetTimeout(timeout int) error {
	if timeout < 0 {
		return errors.New("invalid timeout")
	}
	c.httpClient.Timeout = time.Second * time.Duration(timeout)

	return nil
}

func (c *Client) SetQPS(qps int) error {
	if qps <= 0 {
		return errors.New("invalid qps")
	}
	c.rateLimiter = rate.NewLimiter(rate.Limit(qps), 1)

	return nil
}

func (c *Client) SetConcurrent(concurrent int) error {
	if concurrent <= 0 {
		return errors.New("invalid concurrent")
	}
	c.concurrent = concurrent

	return nil
}

type Result struct {
	URL      string
	Code     int
	Redirect string
}

func (r *Result) String() string {
	var output string
	if r.Redirect != "" {
		output = fmt.Sprintf("%s -> %s", r.URL, r.Redirect)
	} else {
		output = r.URL
	}

	return fmt.Sprintf("%d - %s", r.Code, output)
}
