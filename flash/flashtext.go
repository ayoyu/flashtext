package flash

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const separator string = "=>"

type TreeNode struct {
	selfRune  rune
	children  map[rune]*TreeNode
	isWord    bool
	cleanWord string
	keep      bool
	key       string
}

func newTreeNode() *TreeNode {
	return &TreeNode{
		children: make(map[rune]*TreeNode),
	}
}

type FlashKeywords struct {
	root          *TreeNode
	size          int // nbr of keys
	nbrNodes      int
	caseSensitive bool
}

func NewFlashKeywords(caseSensitive bool) *FlashKeywords {
	return &FlashKeywords{
		root:          newTreeNode(),
		nbrNodes:      1,
		caseSensitive: caseSensitive,
	}
}

func (tree *FlashKeywords) Size() int {
	return tree.size
}

func (tree *FlashKeywords) GetAllKeywords() map[string]string {
	key2Clean := make(map[string]string, tree.size)
	stack := make([]*TreeNode, 0, tree.nbrNodes)
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
			currentNode.children[char] = newTreeNode()
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

func (tree *FlashKeywords) Add(word string) {
	tree.addKeyWord(word, "")
}

func (tree *FlashKeywords) AddKeyWord(word string, cleanWord string) {
	tree.addKeyWord(word, cleanWord)
}

func (tree *FlashKeywords) AddFromMap(keys2synonyms map[string][]string) {
	// keyword_dict = {
	//  "java": ["java_2e", "java programing"],
	//  "product management": ["PM", "product manager"]
	// }
	for key, listSynonyms := range keys2synonyms {
		for _, synonym := range listSynonyms {
			tree.addKeyWord(synonym, key)
		}
	}
}

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

func (tree *FlashKeywords) RemoveKey(word string) bool {
	var nextNode *TreeNode
	parent := make(map[*TreeNode]*TreeNode)
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
	for currentNode != tree.root && len(currentNode.children) == 0 && !currentNode.isWord {
		parentNode := parent[currentNode]
		childRune := currentNode.selfRune
		currentNode = nil
		tree.nbrNodes--
		delete(parentNode.children, childRune)
		currentNode = parentNode
	}
	return true
}

type Result struct {
	key       string
	isPrefix  bool // support for key the smallest(the prefix) and the longest match
	cleanWord string
	start     int
	end       int
}

func (tree *FlashKeywords) Search(text string) []Result {
	n := len(text)
	if !tree.caseSensitive {
		text = strings.ToLower(text)
	}
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
					key:       currentNode.key,
					isPrefix:  isPrefix,
					cleanWord: currentNode.cleanWord,
					start:     start,
					end:       idx,
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
				// fmt.Println("runSize: ", runeKeySize)
				start := bufSize - runeKeySize
				//fmt.Println("start: ", start, buf)
				lastChange = start
				for _, cChar := range currentNode.cleanWord {
					//fmt.Println("cChar: ", cChar, string(cChar))
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
					//fmt.Println("start: ", start, buf)

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
