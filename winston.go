package winston

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal("[ERROR]", err)
	}
}

type Winston struct {
	Location string
	Text     string
	SafeText string
	Grams    []string
	Freq     map[string]int
}

func (w1 *Winston) CommonFreqKeys(w2 *Winston) []string {
	common := make([]string, 0)

	for key, _ := range w1.Freq {
		if w2.Freq[key] != 0 {
			common = append(common, key)
		}
	}

	return common
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
	for _, key := range w1.CommonFreqKeys(w2) {
		sum += w1.Freq[key] * w2.Freq[key]
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
	CheckError(err)

	tagregexp, err := regexp.Compile("<[^>]+>")
	CheckError(err)

	spaceregexp, err := regexp.Compile("[ ]+")
	CheckError(err)

	w.SafeText = tagregexp.ReplaceAllString(w.Text, " ")
	w.SafeText = asciiregexp.ReplaceAllString(w.SafeText, " ")
	w.SafeText = spaceregexp.ReplaceAllString(w.SafeText, " ")
	w.SafeText = strings.Trim(w.SafeText, "")
	w.SafeText = strings.ToLower(w.SafeText)
	w.SafeText = strings.TrimSpace(w.SafeText)
}

func (w *Winston) CalcGrams() {
	w.CleanText()

	w.Grams = strings.Split(w.SafeText, ` `)
	w.Freq = make(map[string]int)

	for _, gram := range w.Grams {
		w.Freq[gram] += 1
	}
}

func (w *Winston) FetchUrl(theurl string) {
	var client *http.Client

	if proxy := os.Getenv("http_proxy"); proxy != `` {
		proxyUrl, err := url.Parse(proxy)
		CheckError(err)

		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest(`GET`, theurl, nil)
	CheckError(err)

	resp, err := client.Do(req)
	CheckError(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)

	w.Text = string(body)
}
