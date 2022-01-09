package flash

import (
	"strings"
)

type TreeNode struct {
	selfRune   rune
	children   map[rune]*TreeNode
	isWord     bool
	cleanWords map[string]struct{}
	sizeCleans int
	keep       bool
	key        string
}

func newTreeNode() *TreeNode {
	return &TreeNode{
		children:   make(map[rune]*TreeNode),
		cleanWords: make(map[string]struct{}),
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

func (tree *FlashKeywords) GetAllKeywords() map[string]map[string]struct{} {
	key2Clean := make(map[string]map[string]struct{}, tree.size)
	stack := make([]*TreeNode, 0, tree.nbrNodes)
	stack = append(stack, tree.root)
	_size := 1
	for _size != 0 {
		node := stack[_size-1]
		stack = stack[:_size-1]
		_size--
		if node.isWord {
			key2Clean[node.key] = node.cleanWords
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
	if cleanWord != "" {
		if _, ok := currentNode.cleanWords[cleanWord]; !ok {
			currentNode.sizeCleans++
		}
		currentNode.cleanWords[cleanWord] = struct{}{}
	}
	if !currentNode.isWord {
		tree.size++
		currentNode.isWord = true
		if len(currentNode.children) != 0 {
			currentNode.keep = true
		}
		currentNode.key = word
	}
}

func (tree *FlashKeywords) Add(word string) {
	tree.addKeyWord(word, "")
}

func (tree *FlashKeywords) AddKeyWord(word string, cleanWord string) {
	tree.addKeyWord(word, cleanWord)
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
	cleanWord string
	start     int
	end       int
}

func (tree *FlashKeywords) Search(text string) []Result {
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
				if currentNode.sizeCleans > 0 {
					for clean := range currentNode.cleanWords {
						res = append(res, Result{
							key:       currentNode.key,
							cleanWord: clean,
							start:     start,
							end:       idx,
						})
					}
				} else {
					res = append(res, Result{
						key:   currentNode.key,
						start: start,
						end:   idx,
					})
				}
				if !currentNode.keep {
					currentNode = tree.root
					start = idx + 1
				}
			}
		}

	}
	return res
}
