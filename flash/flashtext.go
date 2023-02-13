package flash

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const separator string = "=>"

type TrieNode struct {
	selfRune  rune
	children  map[rune]*TrieNode
	isWord    bool
	cleanWord string
	keep      bool
	key       string
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

type FlashKeywords struct {
	root          *TrieNode
	size          int // nbr of keys
	nbrNodes      int
	caseSensitive bool
}

// Instantiate a new Instance of the `FlashKeywords` with
// a case sensitive true or false
func NewFlashKeywords(caseSensitive bool) *FlashKeywords {
	return &FlashKeywords{
		root:          newTrieNode(),
		nbrNodes:      1,
		caseSensitive: caseSensitive,
	}
}

// Returns the number of the keys inside the keys dictionary
func (tree *FlashKeywords) Size() int {
	return tree.size
}

// Returns a map of all the keys in the trie with their `cleanWord`
func (tree *FlashKeywords) GetAllKeywords() map[string]string {
	key2Clean := make(map[string]string, tree.size)
	stack := make([]*TrieNode, 0, tree.nbrNodes)
	stack = append(stack, tree.root)
	_size := 1
	for _size != 0 {
		node := stack[_size-1]
		stack = stack[:_size-1]
		_size--
		if node.isWord {
			key2Clean[node.key] = node.cleanWord
		}
		for _, child := range node.children {
			stack = append(stack, child)
			_size++
		}
	}

	return key2Clean
}

func (tree *FlashKeywords) addKeyWord(word string, cleanWord string) {
	if !tree.caseSensitive {
		word = strings.ToLower(word)
		cleanWord = strings.ToLower(cleanWord)
	}
	currentNode := tree.root
	for _, char := range word {
		if currentNode.isWord {
			currentNode.keep = true
		}
		if _, ok := currentNode.children[char]; !ok {
			currentNode.children[char] = newTrieNode()
			tree.nbrNodes++
		}
		currentNode = currentNode.children[char]
		currentNode.selfRune = char
	}
	if !currentNode.isWord {
		tree.size++
		currentNode.isWord = true
		if len(currentNode.children) != 0 {
			currentNode.keep = true
		}
		currentNode.key = word
		currentNode.cleanWord = cleanWord
	} else if cleanWord != "" {
		if currentNode.cleanWord != "" {
			log.Printf("Warning: overwrite the clean word of %s from %s to %s",
				word, currentNode.cleanWord, cleanWord)
		}
		currentNode.cleanWord = cleanWord
	}
}

// Add the key `word` into the trie
func (tree *FlashKeywords) Add(word string) {
	tree.addKeyWord(word, "")
}

// Add the key `word` into the trie with the corresponding `cleanWord`
func (tree *FlashKeywords) AddKeyWord(word string, cleanWord string) {
	tree.addKeyWord(word, cleanWord)
}

// Add Multiple Keywords simultaneously from a map example:
//	keyword_dict = {
//  	"java": ["java_2e", "java programing"],
//  	"product management": ["PM", "product manager"]
// 	}
// 	trie.AddFromMap(keyword_dict)
func (tree *FlashKeywords) AddFromMap(keys2synonyms map[string][]string) {

	for key, listSynonyms := range keys2synonyms {
		for _, synonym := range listSynonyms {
			tree.addKeyWord(synonym, key)
		}
	}
}

// Add Multiple Keywords simultaneously from a file by providing the `filePath`
func (tree *FlashKeywords) AddFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		synonym2key := strings.Split(line, separator)
		if len(synonym2key) == 2 {
			tree.addKeyWord(synonym2key[0], synonym2key[1])
		} else if len(synonym2key) == 1 {
			tree.addKeyWord(synonym2key[0], "")
		} else {
			// skip line
			log.Printf("Warning: Skipped malformed line %s correct format: key1=>key2", line)
		}
	}
	return nil
}

// Returns the corresponding `cleanWord` for the key `word` from the trie
func (tree *FlashKeywords) GetKeysWord(word string) (string, error) {
	currentNode := tree.root
	for _, char := range word {
		currentNode = currentNode.children[char]
		if currentNode == nil {
			return "", fmt.Errorf("The word %s doesn't exists in the keywords dictionnary", word)
		}
	}
	if !currentNode.isWord {
		return "", fmt.Errorf("The word %s doesn't exists in the keywords dictionnary", word)
	}

	return currentNode.cleanWord, nil
}

// Check if the key `word` exists in the trie dictionary
func (tree *FlashKeywords) Contains(word string) bool {
	currentNode := tree.root
	for _, char := range word {
		currentNode = currentNode.children[char]
		if currentNode == nil {
			return false
		}
	}
	return currentNode.isWord
}

// Remove the key `word` from the trie dictionary
func (tree *FlashKeywords) RemoveKey(word string) bool {
	var nextNode *TrieNode
	parent := make(map[*TrieNode]*TrieNode)
	currentNode := tree.root
	for _, currChar := range word {
		if _, ok := currentNode.children[currChar]; !ok {
			return false
		}
		nextNode = currentNode.children[currChar]
		parent[nextNode] = currentNode
		currentNode = nextNode
	}
	if !currentNode.isWord {
		return false
	}
	currentNode.isWord = false
	tree.size--
	var parentNode *TrieNode
	var childRune rune
	for currentNode != tree.root && len(currentNode.children) == 0 && !currentNode.isWord {
		parentNode = parent[currentNode]
		childRune = currentNode.selfRune
		currentNode = nil
		tree.nbrNodes--
		delete(parentNode.children, childRune)
		currentNode = parentNode
	}
	return true
}

// the resulting output struct:
//	- `Key`: the string keyword found in the search text
//	- `IsPrefix` (false/true): indicates if the key A is a prefix of another string(key B)
//		where A and B are both in the dictionary of the flash keywords
//	- `CleanWord`: the string with which the found key will be replaced in the text.
//               We can think of it also like the origin word of the synonym found in the text.
//	- `Start & End`: span information about the start and end indexes if the key found in the text
type Result struct {
	Key       string
	IsPrefix  bool // support for key the smallest(the prefix) and the longest match
	CleanWord string
	Start     int
	End       int
}

// Search in the text for the stored keys in the trie and
// returns a slice of `Result`
func (tree *FlashKeywords) Search(text string) []Result {
	if !tree.caseSensitive {
		text = strings.ToLower(text)
	}
	n := len(text)
	var res []Result
	currentNode := tree.root
	start := 0
	for idx, char := range text {
		currentNode = currentNode.children[char]
		if currentNode == nil {
			currentNode = tree.root
			start = idx + 1
		} else {
			if currentNode.isWord {
				isPrefix := false
				if currentNode.keep {
					// possibility to be a prefix of another continous word
					if idx+1 < n {
						if _, ok := currentNode.children[rune(text[idx+1])]; ok {
							isPrefix = true
						}
					}
				}
				res = append(res, Result{
					Key:       currentNode.key,
					IsPrefix:  isPrefix,
					CleanWord: currentNode.cleanWord,
					Start:     start,
					End:       idx,
				})
				if !isPrefix {
					// go back to root with 2 conditions (see TestGoBackToRootTrick):
					// 	- simple one if keep=false (isPrefix=false by default)
					// 	- keep can be true but when we look one step ahead
					// 	  no node is founded => Go back to root
					currentNode = tree.root
					start = idx + 1
				}
			}
		}

	}
	return res
}

// Replace the keys found in the text with their `cleanWord` if it exists
// and returns a new string with the replaced keys
func (tree *FlashKeywords) Replace(text string) string {
	if !tree.caseSensitive {
		text = strings.ToLower(text)
	}
	n := len(text)
	var buf []rune = make([]rune, 0, n)
	bufSize := 0
	currentNode := tree.root
	// track the tail of the buf to know if append or set with a new rune buf[lastChange]
	// in the case lenght of key is different from lenght of cleanWord
	lastChange := 0
	for idx, char := range text {
		if lastChange < bufSize {
			buf[lastChange] = char
			lastChange++
		} else {
			buf = append(buf, char)
			bufSize++
			lastChange = bufSize
		}
		currentNode = currentNode.children[char]
		if currentNode == nil {
			currentNode = tree.root
		} else if currentNode.isWord {
			if currentNode.cleanWord != "" {
				// repalce opp `leftmost match first`(replace key with the cleanWord)
				runeKeySize := len([]rune(currentNode.key))
				start := bufSize - runeKeySize
				lastChange = start
				for _, cChar := range currentNode.cleanWord {
					if start < bufSize {
						buf[start] = cChar
						start++
						lastChange++
					} else {
						buf = append(buf, cChar)
						bufSize++
						lastChange = bufSize
						start = bufSize
					}

				}
				// done with replacement Go back to root
				currentNode = tree.root
			} else if currentNode.keep {
				// in case the currentNode(isWord) doesn't have a cleanWord
				// worth to check if it is a prefix of another `big word`
				// to make the replacement opp on the `big word`
				if _, ok := currentNode.children[rune(text[idx+1])]; !ok {
					// nothing found go back to root
					currentNode = tree.root
				}
			} else {
				currentNode = tree.root
			}

		}

	}

	return string(buf[:lastChange])
}
