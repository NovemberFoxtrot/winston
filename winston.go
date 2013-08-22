package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func checkerror(err error) {
	if err != nil {
		log.Fatal("[ERROR]", err)
	}
}

type Winston struct {
	Text  string
	Grams []string
	Freq  map[string]int
}

func (w *Winston) CleanText() {
	asciiregexp, err := regexp.Compile("[^A-Za-z ]+")
	checkerror(err)

	tagregexp, err := regexp.Compile("<[^>]+>")
	checkerror(err)

	spaceregexp, err := regexp.Compile("[ ]+")
	checkerror(err)

	w.Text = tagregexp.ReplaceAllString(w.Text, " ")
	w.Text = asciiregexp.ReplaceAllString(w.Text, " ")
	w.Text = spaceregexp.ReplaceAllString(w.Text, " ")
	w.Text = strings.Trim(w.Text, "")
	w.Text = strings.ToLower(w.Text)
	w.Text = strings.TrimSpace(w.Text)
}

func (w *Winston) CalcGrams() {
	w.CleanText()

	w.Grams = strings.Split(w.Text, ` `)
	w.Freq = make(map[string]int)

	for _, gram := range w.Grams {
		w.Freq[gram] += 1
	}
}

func (w *Winston) FetchUrl(theurl string) {
	var client *http.Client

	if proxy := os.Getenv("http_proxy"); proxy != `` {
		proxyUrl, err := url.Parse(proxy)
		checkerror(err)

		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest(`GET`, theurl, nil)
	checkerror(err)

	resp, err := client.Do(req)
	checkerror(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkerror(err)

	w.Text = string(body)
}

func main() {
	winstons := make([]Winston, 0)

	for _, url := range os.Args[1:] {
		var w Winston
		w.FetchUrl(url)
		w.CalcGrams()
		fmt.Println(len(w.Text), len(w.Grams), len(w.Freq))
		winstons = append(winstons, w)
	}
}
