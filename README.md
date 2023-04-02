# FlashText

This package is a golang version of the original python library [flashtext](https://github.com/vi3k6i5/flashtext), based on the [FlashText algorithm](https://arxiv.org/abs/1711.00046) which is a special version of the [Aho-Corasick algorithm](https://en.wikipedia.org/wiki/Aho%E2%80%93Corasick_algorithm).

The utility of the package is focused on keywords **extraction and replacement** with fixed strings at **scale**, the time complexity of the algorithm doesn't depend on the number of keys being searched or replaced. For a document of size `N` (characters) and a dictionary of `M` keys to **search/replace**, the time complexity is `O(N)`.

`Flashtext` doesn't do regular expression and it's not a replacement of `regex`

# Installation

```
$ go get github.com/ayoyu/flashtext
```

# Usage Overview

## Search and extract keywords

#### caseSensitive=false:

```golang
package main

import (
	"fmt"

	"github.com/ayoyu/flashtext"
)

func main(){
	flashKeys := flashtext.NewFlashKeywords(false)
	// add to dictionnary the key `Apple` with its cleanWord `Fruit`.
	flashKeys.AddKeyWord("Apple", "Fruit")
	// add to dictionnary the key `Apple` with its cleanWord `Company`.
	flashKeys.AddKeyWord("FootBall", "Sport")
	// add to dictionnary the key without a cleanWord.
	flashKeys.Add("Banana")

	// It depends on the context of the use case, the cleanWord can be
    // seen as a synonym to its key or a label/entity to describe its key,...etc.
	// (Similar to the Elasticsearch `Synonym token filter` functionality)

	// The search will extract the keys that exists in the flash keyword dictionnary
	// with their `cleanWords`
	res := flashKeys.Search("I played football, while eating my apple")
	for _, item := range res {
		fmt.Println(item)
	}
}
```

- Output

```bash
{football false sport 9 16}
{apple false fruit 35 39}
```

#### caseSensitive=true:

From the previous example with the same added keys:

```golang
flashKeys := flashtext.NewFlashKeywords(true)
// ...
res := flashKeys.Search("I played football, while eating my Apple")
for _, item := range res {
	fmt.Println(item)
}
```

- Ouput

```
{Apple false Fruit 35 39}
```

The structure of the resulting output is the following:

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

- `IsPrefix`: Indicates if the key A is a prefix of another key B
			  where A and B are both in the dictionary of the flash keywords

- `CleanWord`: The string with which the found key will be replaced in the text.
               It depends on the context of the use case, the cleanWord can be
               seen as a synonym to its key or a label/entity to describe its key,...etc (Similar to the Elasticsearch `Synonym token filter` functionality)

- `Start & End`: span information about the start and end indexes if the key found in the text
```

## Replace keywords

Replace the keys added to the flash keywords with their `clean words` if they exist in the document.

In this example `FootBall` and `Apple` will be replaced respectively with their clean words `Sport` and `Fruit`, but in the case of the key `Banana`, it doesn't get replaced because the key has no `cleanWord`.

#### caseSensitive=false:

```golang
import (
	"fmt"

	"github.com/ayoyu/flashtext"
)

func main() {
	flashKeys := flashtext.NewFlashKeywords(false)
	flashKeys.AddKeyWord("Apple", "Fruit")
	flashKeys.AddKeyWord("FootBall", "Sport")
	flashKeys.AddKeyWord("🔥", "💪")
	flashKeys.Add("Banana")

	text := "I played football, while eating my Apple 🔥"
	fmt.Println("New text: ", flashKeys.Replace(text))
}
```

- Output

```
New text:  i played sport, while eating my fruit 💪
```

#### caseSensitive=true:

```golang
import (
	"fmt"

	"github.com/ayoyu/flashtext"
)

func main() {
	flashKeys := flashtext.NewFlashKeywords(true)
	flashKeys.AddKeyWord("Apple", "Fruit")
	flashKeys.AddKeyWord("FootBall", "Sport")
	flashKeys.AddKeyWord("🔥", "💪")
	flashKeys.Add("Banana")

	text := "I played football, while eating my Apple 🔥"
	fmt.Println("New text: ", flashKeys.Replace(text))
}
```

- Output:

```
New text:  I played football, while eating my Fruit 💪
```

To check the documentation of all the methods and the functions in your browser, type in your terminal:

```
$ godoc -http=:8080
```

Then browse to: `http://localhost:8080/pkg/github.com/ayoyu/flash`

# Benchmarks

First generate some random data or use the ones that are located in the `testdata` folder

```shell
$ go run testdata/gen_testdata.go
```

`cd` to the `benchmarks` folder.

- Benchmark FlashText Search

```shell
$ go test -benchmem -run=^$ -bench ^BenchmarkFlashTextSearch$
Search on a corpus text of size:  669317
goos: linux
goarch: amd64
pkg: github.com/ayoyu/flashtext/benchmarks
cpu: Intel(R) Core(TM) i5-9300H CPU @ 2.40GHz
BenchmarkFlashTextSearch/key_size=10-8         	     162	   7449500 ns/op	     848 B/op	       4 allocs/op
BenchmarkFlashTextSearch/key_size=1010-8       	      85	  22213111 ns/op	 2653670 B/op	      19 allocs/op
BenchmarkFlashTextSearch/key_size=2010-8       	      50	  31000992 ns/op	 5930451 B/op	      22 allocs/op
BenchmarkFlashTextSearch/key_size=3010-8       	      32	  38498148 ns/op	 7642576 B/op	      23 allocs/op
BenchmarkFlashTextSearch/key_size=4010-8       	      39	  33514732 ns/op	 9797074 B/op	      24 allocs/op
BenchmarkFlashTextSearch/key_size=5010-8       	      36	  34218947 ns/op	12508624 B/op	      25 allocs/op
BenchmarkFlashTextSearch/key_size=6010-8       	      33	  42046412 ns/op	15916498 B/op	      26 allocs/op
BenchmarkFlashTextSearch/key_size=7010-8       	      28	  55664084 ns/op	15916496 B/op	      26 allocs/op
BenchmarkFlashTextSearch/key_size=8010-8       	      30	  36441695 ns/op	20192720 B/op	      27 allocs/op
BenchmarkFlashTextSearch/key_size=9010-8       	      33	  57783636 ns/op	20192720 B/op	      27 allocs/op
BenchmarkFlashTextSearch/key_size=10010-8      	      37	  45879421 ns/op	20192720 B/op	      27 allocs/op
BenchmarkFlashTextSearch/key_size=11010-8      	      33	  52046617 ns/op	25550288 B/op	      28 allocs/op
BenchmarkFlashTextSearch/key_size=12010-8      	      27	  38160474 ns/op	25550291 B/op	      28 allocs/op
BenchmarkFlashTextSearch/key_size=13010-8      	      33	  51766075 ns/op	25550290 B/op	      28 allocs/op
BenchmarkFlashTextSearch/key_size=14010-8      	      30	  52940062 ns/op	25550288 B/op	      28 allocs/op
BenchmarkFlashTextSearch/key_size=15010-8      	      20	  59203761 ns/op	32259536 B/op	      29 allocs/op
BenchmarkFlashTextSearch/key_size=16010-8      	      31	  45628554 ns/op	32259536 B/op	      29 allocs/op
BenchmarkFlashTextSearch/key_size=17010-8      	      26	  45658666 ns/op	32259536 B/op	      29 allocs/op
BenchmarkFlashTextSearch/key_size=18010-8      	      27	  39414959 ns/op	32259536 B/op	      29 allocs/op
BenchmarkFlashTextSearch/key_size=19010-8      	      25	  49583917 ns/op	32259536 B/op	      29 allocs/op
BenchmarkFlashTextSearch/key_size=20010-8      	      27	  47829831 ns/op	40664531 B/op	      30 allocs/op
PASS
ok  	github.com/ayoyu/flashtext/benchmarks	42.381s

```

- Benchmark Regex Search

```shell
$ go test -benchmem -run=^$ -bench ^BenchmarkRegexSearch$
Search on a corpus text of size:  669317
goos: linux
goarch: amd64
pkg: github.com/ayoyu/flashtext/benchmarks
cpu: Intel(R) Core(TM) i5-9300H CPU @ 2.40GHz
BenchmarkRegexSearch/key_size=10-8         	       9	 116758673 ns/op	    2760 B/op	      10 allocs/op
BenchmarkRegexSearch/key_size=1010-8       	       1	16455596425 ns/op	  779272 B/op	    2785 allocs/op
BenchmarkRegexSearch/key_size=2010-8       	       1	34206170218 ns/op	 1505056 B/op	    5548 allocs/op
BenchmarkRegexSearch/key_size=3010-8       	       1	52015943197 ns/op	 2204912 B/op	    8185 allocs/op
BenchmarkRegexSearch/key_size=4010-8       	       1	72208914019 ns/op	 2968072 B/op	   10852 allocs/op
BenchmarkRegexSearch/key_size=5010-8       	       1	87424594909 ns/op	 3648112 B/op	   13442 allocs/op
BenchmarkRegexSearch/key_size=6010-8       	       1	111404534897 ns/op	 4469488 B/op	   16163 allocs/op
BenchmarkRegexSearch/key_size=7010-8       	       1	129792983299 ns/op	 5183560 B/op	   18770 allocs/op
BenchmarkRegexSearch/key_size=8010-8       	       1	153291348479 ns/op	 6027576 B/op	   21370 allocs/op
SIGQUIT: quit
PC=0x468301 m=0 sigcode=0

goroutine 0 [idle]:
runtime.futex()
	/usr/local/go/src/runtime/sys_linux_amd64.s:559 +0x21 fp=0x7ffda567d168 sp=0x7ffda567d160 pc=0x468301
runtime.futexsleep(0x467ef3?, 0xa567d1e8?, 0x444eeb?)
	/usr/local/go/src/runtime/os_linux.go:69 +0x36 fp=0x7ffda567d1b8 sp=0x7ffda567d168 pc=0x432256
runtime.notesleep(0x60e5c8)
	/usr/local/go/src/runtime/lock_futex.go:160 +0x87 fp=0x7ffda567d1f0 sp=0x7ffda567d1b8 pc=0x40b787
runtime.mPark(...)
	/usr/local/go/src/runtime/proc.go:1457
runtime.stopm()
	/usr/local/go/src/runtime/proc.go:2239 +0x8c fp=0x7ffda567d220 sp=0x7ffda567d1f0 pc=0x43cb2c
runtime.findRunnable()
	/usr/local/go/src/runtime/proc.go:2866 +0x9e8 fp=0x7ffda567d310 sp=0x7ffda567d220 pc=0x43e1c8
runtime.schedule()
	/usr/local/go/src/runtime/proc.go:3206 +0xbe fp=0x7ffda567d348 sp=0x7ffda567d310 pc=0x43effe
runtime.goschedImpl(0xc000202820)
	/usr/local/go/src/runtime/proc.go:3370 +0xc5 fp=0x7ffda567d380 sp=0x7ffda567d348 pc=0x43f625
runtime.gosched_m(0xc000202820?)
	/usr/local/go/src/runtime/proc.go:3378 +0x31 fp=0x7ffda567d3a0 sp=0x7ffda567d380 pc=0x43f7b1
runtime.mcall()
	/usr/local/go/src/runtime/asm_amd64.s:448 +0x43 fp=0x7ffda567d3b0 sp=0x7ffda567d3a0 pc=0x464303
...

*** Test killed with quit: ran too long (11m0s).
exit status 2
FAIL	github.com/ayoyu/flashtext/benchmarks	660.014s

```

```
$ go run simple_flash_regex_bench.go
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

# Citation

- The python library: https://github.com/vi3k6i5/flashtext
- The original paper: https://arxiv.org/abs/1711.00046

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

## TODO

- [ ] Make the data structure Thread-safe to be able to use it in a concurrent environment.
