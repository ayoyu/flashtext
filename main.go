package main

import (
	"fmt"

	"github.com/ayoyu/flashtext/flash"
)

func main() {
	trie := flash.NewFlashKeywords(false)
	fmt.Println(trie.Size())
}
