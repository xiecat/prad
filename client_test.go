package prad

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

var wordlist = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

func TestNewClient(t *testing.T) {
	_, err := NewClient(wordlist)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}
}

func TestClient_Check(t *testing.T) {
	c, err := NewClient(wordlist)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}

	target := "https://github.com"
	word := "1"
	r, err := c.Check(context.Background(), target, word)
	if err != nil {
		t.Fatal("Check failed:", err)
	}
	if r.URL != target+"/"+word {
		t.Fatal("Generate url wrong:", r.URL)
	}

	target = "https://github.com/{{}}/2"
	word = "1"
	r, err = c.Check(context.Background(), target, word)
	if err != nil {
		t.Fatal("Check failed:", err)
	}
	if r.URL != strings.ReplaceAll(target, "{{}}", word) {
		t.Fatal("Generate url wrong:", r.URL)
	}
}

func TestClient_Do(t *testing.T) {
	c, err := NewClient(wordlist)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}

	target := "https://github.com"
	resultChan, err := c.Do(context.Background(), target)
	if err != nil {
		t.Fatal("Do failed:", err)
	}

	for r := range resultChan {
		fmt.Println(r)
	}
}

func TestClient_SetProxy(t *testing.T) {
	c, err := NewClient(wordlist)
	if err != nil {
		t.Fatal("NewClient failed:", err)
	}
	err = c.SetProxy("http://127.0.0.1:8080")
	if err != nil {
		t.Fatal("Set proxy failed:", err)
	}
	result, err := c.Do(context.Background(), "https://github.com")
	if err != nil {
		t.Fatal("Do failed:", err)
	}
	for range result {
	}
}
