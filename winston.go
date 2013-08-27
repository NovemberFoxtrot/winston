package winston

import (
	"fmt"
	"math"
	"regexp"
	"sir"
	"strings"
	"unicode"
)

type Documents []*Document

type Document struct {
	Location  string
	Text      string
	SafeText  string
	Sentences []int
	Grams     []string
	Freq      map[string]int
	BiFreq    map[string]int
}

func (d1 *Document) CommonFreqKeys(d2 *Document) []string {
	common := make([]string, 0)

	for key, _ := range d1.Freq {
		if d2.Freq[key] != 0 {
			common = append(common, key)
		}
	}

	return common
}

func (w *Document) FreqSum() (sum int) {
	for _, count := range w.Freq {
		sum += count
	}

	return
}

func (w *Document) FreqSquare() (sum float64) {
	for _, count := range w.Freq {
		sum += math.Pow(float64(count), 2)
	}

	return
}

func (w1 *Document) FreqProduct(w2 *Document) (sum int) {
	for _, key := range w1.CommonFreqKeys(w2) {
		sum += w1.Freq[key] * w2.Freq[key]
	}

	return
}

func (w1 *Document) Pearson(w2 *Document) float64 {
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

func (w *Document) CleanText() {
	asciiregexp, err := regexp.Compile("[^A-Za-z ]+")
	sir.CheckError(err)

	tagregexp, err := regexp.Compile("<[^>]+>")
	sir.CheckError(err)

	spaceregexp, err := regexp.Compile("[ ]+")
	sir.CheckError(err)

	w.SafeText = tagregexp.ReplaceAllString(w.Text, " ")
	w.SafeText = asciiregexp.ReplaceAllString(w.SafeText, " ")
	w.SafeText = spaceregexp.ReplaceAllString(w.SafeText, " ")
	w.SafeText = strings.Trim(w.SafeText, "")
	w.SafeText = strings.ToLower(w.SafeText)
	w.SafeText = strings.TrimSpace(w.SafeText)
}

func (w *Document) MarkSentenceBoundaries() {
	w.Sentences = make([]int, 0)

	for index, r := range w.Text {
		if !unicode.IsLetter(r) && r == 46 {
			w.Sentences = append(w.Sentences, index)
		}
	}
}

func (w *Document) FetchSentences() {
	for i := 0; i < (len(w.Sentences) - 1); i++ {
		fmt.Println(i, w.Text[w.Sentences[i]:w.Sentences[i+1]])
	}
}

func (d *Document) CalcGrams() {
	d.CleanText()

	d.MarkSentenceBoundaries()

	d.Grams = strings.Split(d.SafeText, ` `)
	d.Freq = make(map[string]int)

	for _, gram := range d.Grams {
		d.Freq[gram] += 1
	}
}

var TheDocuments []Document

func init() {
	TheDocuments = make([]Document, 0)
}
