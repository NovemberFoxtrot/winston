package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func checkerror(err error) {
	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(1)
	}
}

type Winston struct {
	Text  string
	words []string
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
	var w Winston

	w.FetchUrl(os.Args[1])

	fmt.Println(w.Text)
}
