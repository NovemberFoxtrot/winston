package winston

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sir"
	"strings"
	"unicode"
)

type indexData map[string][]*document

type index struct {
	data indexData
}

func (i *index) update(w *document) {
	for _, gram := range w.grams {
		if i.data[gram] == nil {
			i.data[gram] = make([]*document, 0)
		}

		i.data[gram] = append(i.data[gram], w)
	}
}

func Query(query string) []string {
	results := make([]string, 0)

	for key, values := range theindex.data {
		if key == query {
			for _, d := range values {
				results = append(results, d.location)
			}
		}
	}

	return results
}

func IndexDataLen() int {
  return len(theindex.data) 
}

func Add(website string) {
	var d document
	d.location = website
	d.fetchUrl(website)
	d.CalcGrams()
	documents = append(documents, d)

	theindex.update(&d)
}

type document struct {
	location  string
	text      string
	safeText  string
	sentences []int
	grams     []string
	freq      map[string]int
}

func (d1 *document) CommonFreqKeys(d2 *document) []string {
	common := make([]string, 0)

	for key, _ := range d1.freq {
		if d2.freq[key] != 0 {
			common = append(common, key)
		}
	}

	return common
}

func (w *document) FreqSum() (sum int) {
	for _, count := range w.freq {
		sum += count
	}

	return
}

func (w *document) FreqSquare() (sum float64) {
	for _, count := range w.freq {
		sum += math.Pow(float64(count), 2)
	}

	return
}

func (w1 *document) FreqProduct(w2 *document) (sum int) {
	for _, key := range w1.CommonFreqKeys(w2) {
		sum += w1.freq[key] * w2.freq[key]
	}

	return
}

func (w1 *document) Pearson(w2 *document) float64 {
	sum1 := float64(w1.FreqSum())
	sum2 := float64(w2.FreqSum())
	sumsq1 := w1.FreqSquare()
	sumsq2 := w2.FreqSquare()
	sump := float64(w1.FreqProduct(w2))
	n := float64(len(w1.freq))

	num := sump - ((sum1 * sum2) / n)
	den := math.Sqrt((sumsq1 - (math.Pow(sum1, 2))/n) * (sumsq2 - (math.Pow(sum2, 2))/n))

	if den == 0 {
		return 0
	}

	return num / den
}

func (w *document) CleanText() {
	asciiregexp, err := regexp.Compile("[^A-Za-z ]+")
	sir.CheckError(err)

	tagregexp, err := regexp.Compile("<[^>]+>")
	sir.CheckError(err)

	spaceregexp, err := regexp.Compile("[ ]+")
	sir.CheckError(err)

	w.safeText = tagregexp.ReplaceAllString(w.text, " ")
	w.safeText = asciiregexp.ReplaceAllString(w.safeText, " ")
	w.safeText = spaceregexp.ReplaceAllString(w.safeText, " ")
	w.safeText = strings.Trim(w.safeText, "")
	w.safeText = strings.ToLower(w.safeText)
	w.safeText = strings.TrimSpace(w.safeText)
}

func (w *document) MarkSentenceBoundaries() {
	w.sentences = make([]int, 0)

	for index, r := range w.text {
		if !unicode.IsLetter(r) && r == 46 {
			w.sentences = append(w.sentences, index)
		}
	}
}

func (w *document) FetchSentences() {
	for i := 0; i < (len(w.sentences) - 1); i++ {
		fmt.Println(i, w.text[w.sentences[i]:w.sentences[i+1]])
	}
}

func (d *document) CalcGrams() {
	d.CleanText()

	d.MarkSentenceBoundaries()

	d.grams = strings.Split(d.safeText, ` `)
	d.freq = make(map[string]int)

	for _, gram := range d.grams {
		d.freq[gram] += 1
	}
}

func (w *document) fetchUrl(theurl string) {
	var client *http.Client

	if proxy := os.Getenv("http_proxy"); proxy != `` {
		proxyUrl, err := url.Parse(proxy)
		sir.CheckError(err)

		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest(`GET`, theurl, nil)
	sir.CheckError(err)

	resp, err := client.Do(req)
	sir.CheckError(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	sir.CheckError(err)

	w.text = string(body)
}

var documents []document
var theindex index

func init() {
	documents = make([]document, 0)
	theindex.data = make(map[string][]*document)
}
