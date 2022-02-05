# FlashText

This package is a golang version of the original python library [flashtext](https://github.com/vi3k6i5/flashtext), based on the [FlashText algorithm](https://arxiv.org/abs/1711.00046) which is a special version of the [Aho-Corasick algorithm](https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm).

The utility of the package is focused on keywords extraction and replacement with fixed strings at **scale**, the time complexity of the algorithm doesn't depend on the number of keys being searched or replaced. For a document of size `N` (characters) and a dictionary of `M` keys to search/replace, the time complexity is `O(N)`.

`Flashtext` doesn't do regular expression and it's not a replacement of `regex`

# Installation

```
$ go get github.com/ayoyu/flashtext
```

# Principale Usage

## Search and extract keywords (caseSensitive=false/true)

```golang
import (
        "fmt"
        "github.com/ayoyu/flashtext/flash"
)
func main(){
        // ************* caseSensitive=false *************************
        flashKeyword1 := flash.NewFlashKeywords(false)
        flashKeyword1.AddKeyWord("Foo", "dummy word") // add the key with it's clean word
        flashKeyword1.Add("Banana") // add the key without a clean word

        fmt.Println("caseSensitive=false: ", flashKeyword1.Search("Got the foo with the Banana"))

        // ************* caseSensitive=true *************************
        flashKeyword2 := flash.NewFlashKeywords(true)
        flashKeyword2.AddKeyWord("Foo", "dummy word")
        flashKeyword2.Add("Banana")

        fmt.Println("caseSensitive=true: ", flashKeyword2.Search("Got the foo with the Banana"))
}

```

```
caseSensitive=false:  [{foo false dummy word 8 10} {banana false  21 26}]

caseSensitive=true:  [{Banana false  21 26}]
```

The format of the resulting output is the following:

```golang
type Result struct {
	Key       string
	IsPrefix  bool
	CleanWord string
	Start     int
	End       int
}
```

```
- `Key`: the string keyword found in the search text

- `IsPrefix` (false/true): indicates if the key A is a prefix of another string(key B)
                           where A and B are both in the dictionary of the flash keywords

- `CleanWord`: the string with which the found key will be replaced in the text.
               We can think of it also like the origin word of the synonym found in the text.

- `Start & End`: span information about the start and end indexes if the key found in the text
```

## Replace keywords (caseSensitive=false/true)

Replace the keys added to the flash keywords with their clean word if it exists. In this example `Foo` and `Zoo` got replaced respectively with their clean words `Dummy word_1` and `Dummy word_2`, but in the case of the `Banana` key, it doesn't get replaced.

```golang
import (
        "fmt"
        "github.com/ayoyu/flashtext/flash"
)
func main(){
        text := "Got the Foo and the Zoo with the Banana"

        // caseSensitive=false
        flashKeyword1 := flash.NewFlashKeywords(true)
        flashKeyword1.AddKeyWord("Foo", "Dummy word_1")
        flashKeyword1.AddKeyWord("Zoo", "Dummy word_2")
        flashKeyword1.Add("Banana")
        fmt.Println("New text(caseSensitive=true): ", flashKeyword1.Replace(text))

        // caseSensitive=false
        flashKeyword2 := flash.NewFlashKeywords(false)
        flashKeyword2.AddKeyWord("Foo", "Dummy word_1")
        flashKeyword2.AddKeyWord("Zoo", "Dummy word_2")
        flashKeyword2.Add("Banana")
        fmt.Println("New text(caseSensitive=false): ", flashKeyword2.Replace(text))

}
```

```
New text(caseSensitive=true):  Got the Dummy word_1 and the Dummy word_2 with the Banana

New text(caseSensitive=false):  got the dummy word_1 and the dummy word_2 with the banana
```

## Other docs

```golang
Add(word string)
AddFromFile(filePath string) error
AddFromMap(keys2synonyms map[string][]string)
AddKeyWord(word string, cleanWord string)
Contains(word string) bool
GetAllKeywords() map[string]string
GetKeysWord(word string) (string, error)
RemoveKey(word string) bool
Replace(text string) string
Search(text string) []Result
Size() int
```

To check the documentation of all the methods and the functions in your browser, type in your terminal:

```
$ godoc -http=:8080
```

Then browse to: `http://localhost:8080/pkg/github.com/ayoyu/flash`

# Benchmark and performance

```
$ go run benchmarks/bench_regex.go
```

```
keys_size | FlashText (s) |  Regex (s)
---------------------------------------
   10     | 0.00053121    | 0.007381449
  1010    | 0.000902698   | 1.021105121
  2010    | 0.001164453   | 2.155324188
  3010    | 0.001272009   | 3.189556999
  4010    | 0.001415052   | 4.489287341
  5010    | 0.00151844    | 5.662644436
  6010    | 0.001601235   | 6.820220812
  7010    | 0.001711219   | 7.845579981
  8010    | 0.001785076   | 9.740038207
...
```

# Differences & Improvements with the original python library

...

# Citation

The original paper: https://arxiv.org/abs/1711.00046

```
@ARTICLE{2017arXiv171100046S,
   author = {{Singh}, V.},
    title = "{Replace or Retrieve Keywords In Documents at Scale}",
  journal = {ArXiv e-prints},
archivePrefix = "arXiv",
   eprint = {1711.00046},
 primaryClass = "cs.DS",
 keywords = {Computer Science - Data Structures and Algorithms},
     year = 2017,
    month = oct,
   adsurl = {http://adsabs.harvard.edu/abs/2017arXiv171100046S},
  adsnote = {Provided by the SAO/NASA Astrophysics Data System}
}
```
