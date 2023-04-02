package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ayoyu/flashtext"
)

const WORDS_FILE_PATH = "./../testdata/words_benchmark_test.txt"
const CORPUS_FILE_PATH = "./../testdata/corpus_benchmark_test.txt"

// test -benchmem -run=^$ -bench ^BenchmarkRegextSearch$
func readGenWordsTestData(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var words []string
	var scanner *bufio.Scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words, nil
}

func readGenCorpusTestData(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", nil
	}
	return string(content), nil
}

func BenchmarkFlashTextSearch(b *testing.B) {
	words, _ := readGenWordsTestData(WORDS_FILE_PATH)
	corpus, _ := readGenCorpusTestData(CORPUS_FILE_PATH)
	fmt.Println("Search on a corpus text of size: ", len(corpus))
	for keysSize := 10; keysSize < 20011; keysSize += 1000 {
		var flash *flashtext.FlashKeywords = flashtext.NewFlashKeywords(true)
		for i := 0; i < keysSize; i++ {
			flash.Add(words[rand.Intn(len(words))])
		}
		b.ResetTimer()
		b.Run(
			fmt.Sprintf("key_size=%d", keysSize), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					flash.Search(corpus)
				}
			},
		)
	}
}

func BenchmarkRegexSearch(b *testing.B) {
	words, _ := readGenWordsTestData(WORDS_FILE_PATH)
	corpus, _ := readGenCorpusTestData(CORPUS_FILE_PATH)
	fmt.Println("Search on a corpus text of size: ", len(corpus))
	for keysSize := 10; keysSize < 20011; keysSize += 1000 {
		var keysWords []string = make([]string, keysSize)
		for i := 0; i < keysSize; i++ {
			keysWords[i] = `\b` + words[rand.Intn(len(words))] + `\b`
		}
		compileString := strings.Join(keysWords[:], "|")
		reCompile := regexp.MustCompile(compileString)
		b.ResetTimer()
		b.Run(
			fmt.Sprintf("key_size=%d", keysSize), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					reCompile.FindAllString(corpus, -1)
				}
			},
		)
	}
}
