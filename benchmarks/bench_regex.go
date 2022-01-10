package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/ayoyu/flashtext/flash"
)

const SIZE = 100000
const CORPUS_SIZE = 5000

func getRandomWord() string {
	rand.Seed(time.Now().UnixNano())
	lenghts := []int{2, 4, 6, 8, 10, 12}
	randLen := lenghts[rand.Intn(len(lenghts))]
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var word []rune = make([]rune, randLen)
	for i := 0; i < randLen; i++ {
		word[i] = letters[rand.Intn(len(letters))]
	}
	return string(word)
}

func main() {
	allWords := make([]string, SIZE)
	for i := 0; i < SIZE; i++ {
		allWords[i] = getRandomWord()
	}
	fmt.Println("keys_size  | FlashText (s) | Regex (s) ")
	for keysSize := 10; keysSize < 20011; keysSize += 1000 {

		tmp := make([]string, CORPUS_SIZE)
		for i := 0; i < CORPUS_SIZE; i++ {
			tmp[i] = allWords[rand.Intn(SIZE)]
		}
		corpus := strings.Join(tmp[:], " ")
		flash := flash.NewFlashKeywords(true)
		keysWords := make([]string, keysSize)
		for i := 0; i < keysSize; i++ {
			tmpKey := allWords[rand.Intn(SIZE)]
			flash.Add(tmpKey)
			keysWords[i] = `\b` + tmpKey + `\b`
		}
		compileString := strings.Join(keysWords[:], "|")
		reCompile := regexp.MustCompile(compileString)

		start1 := time.Now()
		reCompile.FindAllString(corpus, -1)
		reTime := time.Since(start1)

		start2 := time.Now()
		flash.Search(corpus)
		flashTime := time.Since(start2)
		fmt.Println(keysSize, `        |`, flashTime.Seconds(),
			`  |`, reTime.Seconds())
	}

}
