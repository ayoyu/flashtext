package flash

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndSize(t *testing.T) {
	trie := NewFlashKeywords(false)
	keys := []string{
		"key1", "key2", "key3", "key4", "key5",
	}
	count := 0
	for _, key := range keys {
		trie.Add(key)
		count++
		assert.Equal(t, trie.Size(), count)
	}
}

func TestAddKeyWordAndSize(t *testing.T) {
	trie := NewFlashKeywords(false)
	keys2Clean := []struct {
		key       string
		cleanWord string
	}{
		{"Earth", "planet"},
		{"Mars", "planet"},
		{"sun", "star"},
	}
	count := 0
	for _, data := range keys2Clean {
		trie.AddKeyWord(data.key, data.cleanWord)
		count++
		assert.Equal(t, trie.Size(), count)
	}
}

func TestGetAllKeywords(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys2Clean := []struct {
		key       string
		cleanWord string
	}{
		{"Earth", "planet"},
		{"Mars", "planet"},
		{"sun", "star"},
		{"python3.5", "python"},
		{"python3", "python"},
	}
	for _, item := range keys2Clean {
		trie.AddKeyWord(item.key, item.cleanWord)
	}
	assert.Equal(t, trie.Size(), len(keys2Clean))
	allKeysWords := trie.GetAllKeywords()
	for _, item := range keys2Clean {
		cleanWord, ok := allKeysWords[item.key]
		assert.Equal(t, ok, true)
		assert.Equal(t, cleanWord, item.cleanWord)
	}
	t.Logf("allKeysWords: %v", allKeysWords)
}

func TestNumberofNodes(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{
		"Java", "JavaJ2E", "Slice", "SlicePizza", "Chetoos", "ChetoosDoritos",
	}
	for _, key := range keys {
		trie.Add(key)
	}
	assert.Equal(t, trie.nbrNodes, 32)
}

func TestInsertShortThenLongSearch(t *testing.T) {
	trie := NewFlashKeywords(true)
	// small work before the long word. The small word
	// is a prefix of the long word
	keys := []string{"cat", "catch"}
	for _, key := range keys {
		trie.Add(key)
	}
	text := "Try to catch this"
	res := trie.Search(text)
	assert.Equal(t, len(res), 2)
	prefixFlag := true
	for i := 0; i < 2; i++ {
		assert.Equal(t, res[i].Key, keys[i])
		assert.Equal(t, text[res[i].Start:res[i].End+1], keys[i])
		assert.Equal(t, res[i].IsPrefix, prefixFlag)
		t.Logf("Found Key %v, PASS", keys[i])
		prefixFlag = false // next result is not a prefix of something
	}
	t.Logf("res: %v", res)
}

func TestInsertLongThenShortSearch(t *testing.T) {
	trie := NewFlashKeywords(true)
	// long word before the small word while the small word
	// is a prefix of the long word
	keys := []string{"catch", "cat"}
	for _, key := range keys {
		trie.Add(key)
	}
	text := "Try to catch this"
	res := trie.Search(text)
	assert.Equal(t, len(res), 2)
	prefixFlag := true
	for i := 0; i < 2; i++ {
		assert.Equal(t, res[i].Key, keys[(i+1)%len(keys)])
		assert.Equal(t, text[res[i].Start:res[i].End+1], keys[(i+1)%len(keys)])
		assert.Equal(t, res[i].IsPrefix, prefixFlag)
		t.Logf("Found Key %v, PASS", keys[(i+1)%len(keys)])
		prefixFlag = false
	}
	t.Logf("res: %v", res)
}

func TestCaseSensitive(t *testing.T) {
	trie := NewFlashKeywords(false)
	keys := []string{"FOo", "bio1A", "HellO"}
	for _, k := range keys {
		trie.Add(k)
	}
	allKeys := trie.GetAllKeywords()
	for i := 0; i < len(keys); i++ {
		_, ok := allKeys[strings.ToLower(keys[i])]
		assert.Equal(t, ok, true)
	}
	t.Logf("allkeys: %v", allKeys)
}

func TestFalseCaseSensitiveSearch(t *testing.T) {
	trie := NewFlashKeywords(false)
	keys := []string{"FoO", "Banana"}
	for _, k := range keys {
		trie.Add(k)
	}
	text := "foO with baNana"
	res := trie.Search(text)
	assert.Equal(t, len(res), len(keys))
	for i := 0; i < len(keys); i++ {
		assert.Equal(t, res[i].Key, strings.ToLower(keys[i]))
	}
	t.Logf("Search Result: %v", res)
}

func TestTrueCaseSensitiveSearch(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"Foo", "Banana"}
	for _, k := range keys {
		trie.Add(k)
	}
	text := "Foo with banana"
	res := trie.Search(text)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Key, keys[0])
	t.Logf("Search Result: %v", res)
}

func TestNoEnglishSearch(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"北京", "欢迎", "你"}
	for _, k := range keys {
		t.Logf("Insert Key: %v", k)
		trie.Add(k)
	}
	text1 := "北京欢迎你"
	res := trie.Search(text1)
	assert.Equal(t, len(res), len(keys))
	for i := 0; i < len(keys); i++ {
		assert.Equal(t, keys[i], res[i].Key)
	}
	t.Logf("Search Result: %v", res)
	Key := "测试"
	t.Logf("Insert Key: %v", Key)
	trie.Add(Key)
	text2 := "3测试"
	res2 := trie.Search(text2)
	assert.Equal(t, len(res2), 1)
	assert.Equal(t, res2[0].Key, Key)
	t.Logf("res: %v", trie.Search(text2))
}

func TestSimpleRemoveKeys(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"banana", "foo"}
	for _, k := range keys {
		trie.Add(k)
	}
	t.Logf("Size before deleting %v: %v", keys[0], trie.Size())
	trie.RemoveKey(keys[0])
	allKeys := trie.GetAllKeywords()
	assert.Equal(t, trie.Size(), 1)
	assert.Equal(t, len(allKeys), 1)
	assert.Equal(t, trie.nbrNodes, 4)
	t.Logf("All keys: %v", allKeys)
}

func TestOverlapRemoveKeys_1(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"cat", "catch"}
	for _, k := range keys {
		trie.Add(k)
	}
	t.Logf("size before deleting key **%v**: %v %v", keys[0], trie.Size(), trie.GetAllKeywords())
	assert.Equal(t, trie.nbrNodes, len(keys[1])+1)
	// deleting the key `cat` will not drop the nodes because `catch`
	// is still filling the space
	trie.RemoveKey(keys[0])
	allKeys := trie.GetAllKeywords()
	assert.Equal(t, trie.Size(), 1)
	assert.Equal(t, len(allKeys), 1)
	assert.Equal(t, trie.nbrNodes, len(keys[1])+1)
	t.Logf("All keys: %v || nbrNode: %v", allKeys, trie.nbrNodes)
}

func TestOverlapRemoveKeys_2(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"cat", "catch"}
	for _, k := range keys {
		trie.Add(k)
	}
	assert.Equal(t, trie.nbrNodes, len(keys[1])+1)
	t.Logf("size before deleting key **%v**: %v %v nbrNode: %v", keys[1], trie.Size(),
		trie.GetAllKeywords(), trie.nbrNodes)
	// deleting the key `catch` will drop the node(`h`) and node(`c`)
	// but will stop because `cat` is still filling the space
	trie.RemoveKey(keys[1])
	assert.Equal(t, trie.Size(), 1)
	assert.Equal(t, len(trie.GetAllKeywords()), 1)
	assert.Equal(t, trie.nbrNodes, len(keys[0])+1)
	t.Logf("Size: %v, allKeys: %v, nbrNode: %v", trie.Size(), trie.GetAllKeywords(), trie.nbrNodes)

}

func TestOverlapRemoveKeys_3(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"abc", "abd", "af"}
	for _, k := range keys {
		trie.Add(k)
	}
	nodes := trie.nbrNodes
	for len(keys) != 0 {
		k := keys[len(keys)-1]
		t.Logf("%v nbrNode: %v", trie.GetAllKeywords(), trie.nbrNodes)
		trie.RemoveKey(k)
		if len(keys) > 1 {
			nodes--
		} else {
			nodes = nodes - 3
		}
		keys = keys[:len(keys)-1]
		assert.Equal(t, trie.nbrNodes, nodes)
		assert.Equal(t, trie.Size(), len(keys))
		tmp := trie.GetAllKeywords()
		for _, kk := range keys {
			_, ok := tmp[kk]
			assert.Equal(t, ok, true)
		}
	}
	assert.Equal(t, trie.nbrNodes, nodes)
	t.Logf("%v nbrNode: %v %v", trie.GetAllKeywords(), trie.nbrNodes, nodes)

}

func TestGoBackToRootTrick(t *testing.T) {
	keys := []string{"chetoos", "055-5647-3456", "chetoosPiza"}
	// the error is coming from no proper space btw chetoos and 055...
	// knowing that we still (keep=true) going because of the keyword chetoosPiza
	text := "call chetoos055-5647-3456 chetoosPiza"
	trie := NewFlashKeywords(true)
	for _, k := range keys {
		trie.Add(k)
	}
	res := trie.Search(text)
	expected := []struct {
		key        string
		flagPrefix bool
	}{
		{"chetoos", false},
		{"055-5647-3456", false},
		{"chetoos", true},
		{"chetoosPiza", false},
	}
	for i := 0; i < len(res); i++ {
		assert.Equal(t, res[i].Key, expected[i].key)
		assert.Equal(t, res[i].IsPrefix, expected[i].flagPrefix)
		t.Logf("PASS Key %v", res[i].Key)
	}
	t.Logf("res: %v", res)
}

func TestAddFromMap(t *testing.T) {
	trie := NewFlashKeywords(true)
	hMap := map[string][]string{
		"java":               {"java_2e", "java programing"},
		"product management": {"PM", "product manager"},
	}
	testdata := []struct {
		key        string
		originWord string
	}{
		{"java_2e", "java"},
		{"java programing", "java"},
		{"PM", "product management"},
		{"product manager", "product management"},
	}
	trie.AddFromMap(hMap)
	t.Logf("trie Size: %v", trie.Size())
	for _, item := range testdata {
		cleanWord, err := trie.GetKeysWord(item.key)
		t.Logf("key: %v  cleanWord: %v", item.key, cleanWord)
		assert.Nil(t, err)
		assert.Equal(t, item.originWord, cleanWord)
	}
}

func TestAddFromFile(t *testing.T) {
	trie := NewFlashKeywords(true)
	trie.AddFromFile("./../testdata/Keys2Synonyms.txt")
	testdata := []struct {
		key        string
		originWord string
	}{
		{"java_2e", "java"},
		{"java programing", "java"},
		{"python3.3", "python"},
		{"pypy", "python"},
		{"python2", "python"},
		{"SWE", "Software Engineer"},
		{"Developper", "Software Engineer"},
		{"Backend Engineer", "Software Engineer"},
		{"Banana", ""},
		{"Chetoos", ""},
	}
	for _, item := range testdata {
		cleanWord, err := trie.GetKeysWord(item.key)
		t.Logf("key: %v  listWords: %v", item.key, cleanWord)
		assert.Nil(t, err)
		assert.Equal(t, item.originWord, cleanWord)
	}
}

func TestContains(t *testing.T) {
	trie := NewFlashKeywords(false)
	trie.Add("FoO")
	// no case sensitive
	assert.Equal(t, trie.Contains("foo"), true)
	assert.Equal(t, trie.Contains("Anything"), false)
}

func TestAddWordMultipleTimes(t *testing.T) {
	trie := NewFlashKeywords(true)
	key := "Foo"
	trie.Add(key)
	cleanWord, err := trie.GetKeysWord(key)
	assert.Nil(t, err)
	assert.Equal(t, cleanWord, "")
	assert.Equal(t, trie.nbrNodes, len(key)+1)
	assert.Equal(t, trie.Size(), 1)
	t.Logf("curr cleanWord: %v", cleanWord)

	trie.addKeyWord("Foo", "Zoo")
	assert.Equal(t, trie.nbrNodes, len(key)+1)
	assert.Equal(t, trie.Size(), 1)

	cleanWord, err = trie.GetKeysWord(key)
	assert.Nil(t, err)
	assert.Equal(t, cleanWord, "Zoo")
	assert.Equal(t, trie.Size(), 1)
	assert.Equal(t, trie.nbrNodes, len(key)+1)
	t.Logf("curr cleanWord: %v", cleanWord)

	trie.addKeyWord("Foo", "Zoo2")
	cleanWord, err = trie.GetKeysWord(key)
	assert.Nil(t, err)
	assert.Equal(t, cleanWord, "Zoo2")
	assert.Equal(t, trie.Size(), 1)
	assert.Equal(t, trie.nbrNodes, len(key)+1)
	t.Logf("curr cleanWord: %v", cleanWord)
}

func TestReplaceCleanWordLessThenKey(t *testing.T) {
	// len(rune(cleanWord)) < len(rune(key))
	// the key is bigger than the cleanWord that will replace it
	// Chetoos > Cat
	trie := NewFlashKeywords(true)
	trie.AddKeyWord("Chetoos", "Cat")
	text := "With Chetoos in place"
	newText := trie.Replace(text)
	rText := "With Cat in place"
	assert.Equal(t, newText, rText)
	t.Logf("newText: %v", newText)
}

func TestReplaceCleanWordGreaterThenKey(t *testing.T) {
	// len(rune(cleanWord)) > len(rune(key))
	// the key is smaller than the cleanWord that will replace it
	// Cat < Chetoos
	trie := NewFlashKeywords(true)
	trie.AddKeyWord("Cat", "Chetoos")
	text := "With Cat in place"
	newText := trie.Replace(text)
	rText := "With Chetoos in place"
	assert.Equal(t, newText, rText)
	t.Logf("newText: %v", newText)
}

func TestReplaceCleanWordSameLenghtKey(t *testing.T) {
	// len(rune(cleanWord)) == len(rune(key))
	// the key is of the same size as the cleanWord that will replace it
	// Cat == Bee
	trie := NewFlashKeywords(true)
	trie.AddKeyWord("Cat", "Bee")
	text := "With Cat in place"
	newText := trie.Replace(text)
	rText := "With Bee in place"
	assert.Equal(t, newText, rText)
	t.Logf("newText: %v", newText)
}
func TestReplaceKeyAtTheEnd(t *testing.T) {
	trie := NewFlashKeywords(true)
	trie.AddKeyWord("in place", "gone")
	text := "With Cat in place"
	newText := trie.Replace(text)
	rText := "With Cat gone"
	assert.Equal(t, newText, rText)
	t.Logf("newText: %v", newText)

	trie.AddKeyWord("in place", "in Peeeeeeeeeace")
	newText2 := trie.Replace(text)
	rText2 := "With Cat in Peeeeeeeeeace"
	assert.Equal(t, newText2, rText2)
	t.Logf("newText: %v", newText2)
	// random test
	trie2 := NewFlashKeywords(true)
	trie2.addKeyWord("055-5647-3456", "055-XXX")
	text3 := "call chetoos055-5647-3456 chetoosPiza"
	new := trie2.Replace(text3)
	assert.Equal(t, new, "call chetoos055-XXX chetoosPiza")
	t.Logf("newText: %v", new)
}

func TestRepalceWithFalseCaseSensitive(t *testing.T) {
	trie := NewFlashKeywords(false)
	trie.addKeyWord("Foo", "JOJO")
	text := "HU! foo KIWI"
	newText := trie.Replace(text)
	rText := "hu! jojo kiwi"
	assert.Equal(t, newText, rText)
	t.Logf("newText: %v", newText)
}
