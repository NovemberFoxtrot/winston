package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
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

func (w *Winston) FreqSum() (sum int) {
	for _, count := range w.Freq {
		sum += count
	}

	return
}

func (w *Winston) FreqSquare() (sum float64) {
	for _, count := range w.Freq {
		sum += math.Pow(float64(count), 2)
	}

	return
}

func (w1 *Winston) FreqProduct(w2 *Winston) (sum int) {
	for key, count := range w1.Freq {
		sum += count * w2.Freq[key]
	}

	return
}

func (w1 *Winston) Pearson(w2 *Winston) float64 {
	sum1 := float64(w1.FreqSum())
	sum2 := float64(w2.FreqSum())
	sumsq1 := w1.FreqSquare()
	sumsq2 := w2.FreqSquare()
	sump := float64(w1.FreqProduct(w2))

	n := float64(len(w1.Freq))

	num := sump - ((sum1 * sum2) / n)
	den := math.Sqrt((sumsq1 - (math.Pow(sum1, 2))/n) * (sumsq2 - (math.Pow(sum2, 2))/n))

	if den == 0 {
		return 0
	}

	return num / den
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

	if len(winstons) < 1 {
		os.Exit(1)
	}

	n := 0

	for n < (len(winstons) - 1) {
		sim := winstons[n].Pearson(&winstons[n+1])
		fmt.Println(sim)
		n += 1
	}
}
