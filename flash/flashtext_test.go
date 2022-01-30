package flash

import (
	"strings"
	"testing"
)

func TestAddAndSize(t *testing.T) {
	trie := NewFlashKeywords(false)
	keys := []string{
		"key1", "key2", "key3", "key4", "key5",
	}
	for i, key := range keys {
		trie.Add(key)
		if trie.Size() != i+1 {
			t.Errorf("trie.Add(%v) FAILD, EXPECTED %d got %d", key, i+1, trie.Size())
		} else {
			t.Logf("trie.Add(%v) | Size: %d PASS", key, trie.Size())
		}
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
	for i, data := range keys2Clean {
		trie.AddKeyWord(data.key, data.cleanWord)
		if trie.Size() != i+1 {
			t.Errorf("trie.AddKeyWord(%v, %v) FAILD, EXPECTED %d got %d", data.key, data.cleanWord, i+1, trie.Size())
		} else {
			t.Logf("trie.AddKeyWord(%v, %v) | Size: %d PASS", data.key, data.cleanWord, trie.Size())
		}
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
	if trie.Size() != len(keys2Clean) {
		t.Errorf("trie.AddKeyWord FAILED, Size is expected to be equal: %d", len(keys2Clean))
	}

	allKeysWords := trie.GetAllKeywords()

	for _, item := range keys2Clean {
		if cleanWord, ok := allKeysWords[item.key]; !ok || cleanWord != item.cleanWord {
			t.Errorf("FAILED for key: %s", item.key)
		} else {
			t.Logf("Key %v with CleanWord %v PASS", item.key, item.cleanWord)
		}
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
	if trie.nbrNodes != 32 {
		t.Errorf("trie.nbrNodes FAILED, EXPECTED %d, got %d", 32, trie.nbrNodes)
	} else {
		t.Logf("trie.nbrNodes = %d, PASS", trie.nbrNodes)
	}
}

func TestInsertShortThenLongSearch(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"cat", "catch"}
	for _, key := range keys {
		trie.Add(key)
	}
	text := "Try to catch this"
	res := trie.Search(text)
	if len(res) != 2 {
		t.Errorf("Res should have 2 results cat and catch")
	} else {
		prefixFlag := true
		for i := 0; i < 2; i++ {
			if res[i].Key == keys[i] && text[res[i].Start:res[i].End+1] == keys[i] && res[i].IsPrefix == prefixFlag {
				t.Logf("Found Key %v, PASS", keys[i])
				prefixFlag = false
			} else {
				t.Errorf("Error Key %v FAILED %v", keys[i], res[i])
			}
		}
	}
	t.Logf("res: %v", res)
}

func TestInsertLongThenShortSearch(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"catch", "cat"}
	for _, key := range keys {
		trie.Add(key)
	}
	text := "Try to catch this"
	res := trie.Search(text)
	if len(res) != 2 {
		t.Errorf("Res should have 2 results cat and catch")
	} else {
		prefixFlag := true
		for i := 0; i < 2; i++ {
			if res[i].Key == keys[(i+1)%len(keys)] && text[res[i].Start:res[i].End+1] == keys[(i+1)%len(keys)] && res[i].IsPrefix == prefixFlag {
				t.Logf("Found Key %v, PASS", keys[(i+1)%len(keys)])
				prefixFlag = false
			} else {
				t.Errorf("Error Key %v FAILED %v", keys[(i+1)%len(keys)], res[i])
			}
		}
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
		if _, ok := allKeys[strings.ToLower(keys[i])]; !ok {
			t.Errorf("FAILED key: %v", keys[i])
		}
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
	if len(res) != len(keys) {
		t.Errorf("FAILED len(res) != len(keys)")
	}
	for i := 0; i < len(keys); i++ {
		if res[i].Key != strings.ToLower(keys[i]) {
			t.Errorf("FAILED res[i].key != strings.ToLower(keys[i]) ")
		}
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
	if len(res) != 1 {
		t.Errorf("FAILED len(res) != 1")
	}
	if res[0].Key != keys[0] {
		t.Errorf("FAILED key: %v", keys[0])
	}
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
	if len(keys) != len(res) {
		t.Errorf("FAILED result")
	}
	for i := 0; i < len(keys); i++ {
		if keys[i] != res[i].Key {
			t.Errorf("FAILED key: %v", keys[i])
		}
	}
	t.Logf("Search Result: %v", res)
	tmpKey := "测试"
	t.Logf("Insert Key: %v", tmpKey)
	trie.Add(tmpKey)
	text2 := "3测试"
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
	if trie.Size() != 1 && len(allKeys) != 1 && trie.nbrNodes != 4 {
		t.Errorf("Size should be equal 1, FAILED deleting key %v", keys[0])
	} else {
		t.Logf("All keys: %v", allKeys)
	}
}

func TestOverlapRemoveKeys_1(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"cat", "catch"}
	for _, k := range keys {
		trie.Add(k)
	}
	t.Logf("size before deleting key **%v**: %v %v", keys[0], trie.Size(), trie.GetAllKeywords())
	trie.RemoveKey(keys[0])
	allKeys := trie.GetAllKeywords()
	if trie.Size() != 1 && len(allKeys) != 1 && trie.nbrNodes != len(keys[1])+1 {
		t.Errorf("Size should be equal 1, FAILED deleting key %v", keys[0])
	} else {
		t.Logf("All keys: %v || nbrNode: %v", allKeys, trie.nbrNodes)
	}
}

func TestOverlapRemoveKeys_2(t *testing.T) {
	trie := NewFlashKeywords(true)
	keys := []string{"cat", "catch"}
	for _, k := range keys {
		trie.Add(k)
	}
	t.Logf("size before deleting key **%v**: %v %v nbrNode: %v",
		keys[1], trie.Size(), trie.GetAllKeywords(), trie.nbrNodes)
	trie.RemoveKey(keys[1])
	if trie.Size() != 1 && len(trie.GetAllKeywords()) != 1 && trie.nbrNodes != len(keys[0])+1 {
		t.Errorf("FAILED deleting key: %v", keys[1])
	} else {
		t.Logf("Size: %v, allKeys: %v, nbrNode: %v",
			trie.Size(), trie.GetAllKeywords(), trie.nbrNodes)
	}

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
		if trie.nbrNodes != nodes {
			t.Errorf("FAILED delete: nbrNodes")
		}
		if trie.Size() != len(keys) {
			t.Errorf("FAILED delete: Size")
		}
		tmp := trie.GetAllKeywords()
		for _, kk := range keys {
			if _, ok := tmp[kk]; !ok {
				t.Errorf("FAILED delete, expected key %v", kk)
			}
		}
	}
	if trie.nbrNodes != nodes {
		t.Errorf("FAILED delete: nbrNodes")
	}
	t.Logf("%v nbrNode: %v %v", trie.GetAllKeywords(), trie.nbrNodes, nodes)

}

func TestGoBackToRootTrick(t *testing.T) {
	keys := []string{"chetoos", "055-5647-3456", "chetoosPiza"}
	// error no proper space ' ' btw chetoos and 055...
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
		if res[i].Key == expected[i].key && res[i].IsPrefix == expected[i].flagPrefix {
			t.Logf("PASS Key %v", res[i].Key)
		} else {
			t.Errorf("FAILED Key %v", res[i].Key)
		}
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
		if err != nil {
			t.Error(err)
		}
		if item.originWord != cleanWord {
			t.Errorf("FAILED Key: %v", item.key)
		}
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
		if err != nil {
			t.Error(err)
		}
		if item.originWord != cleanWord {
			t.Errorf("FAILED Key: %v", item.key)
		}

	}
}

func TestContains(t *testing.T) {
	trie := NewFlashKeywords(false)
	trie.Add("FoO")
	// no case sensitive
	if !trie.Contains("foo") {
		t.Errorf("FAILED Key Foo")
	}
	if trie.Contains("Anything") {
		t.Errorf("FAILED")
	}
}

func TestAddWordMultipleTimes(t *testing.T) {
	trie := NewFlashKeywords(true)
	key := "Foo"
	trie.Add(key)
	cleanWord, err := trie.GetKeysWord(key)
	if (cleanWord != "" || err != nil) && trie.nbrNodes != len(key)+1 && trie.Size() != 1 {
		t.Errorf("FAIL nbrNodes must equal %v & size=1", len(key)+1)
	} else {
		t.Logf("curr cleanWord: %v", cleanWord)
	}
	trie.addKeyWord("Foo", "Zoo")
	if trie.nbrNodes != len(key)+1 && trie.Size() != 1 {
		t.Errorf("FAIL nbrNodes must equal %v & size=1", len(key)+1)
	}
	cleanWord, err = trie.GetKeysWord(key)
	if (cleanWord != "Zoo" || err != nil) && trie.Size() != 1 && trie.nbrNodes != len(key)+1 {
		t.Errorf("FAIL")
	} else {
		t.Logf("curr cleanWord: %v", cleanWord)
	}
	trie.addKeyWord("Foo", "Zoo2")
	cleanWord, err = trie.GetKeysWord(key)
	if (cleanWord != "Zoo2" || err != nil) && trie.Size() != 1 && trie.nbrNodes != len(key)+1 {
		t.Errorf("FAIL")
	} else {
		t.Logf("curr cleanWord: %v", cleanWord)
	}
}

func TestReplaceCleanWordLessThenKey(t *testing.T) {
	// len(rune(cleanWord)) < len(rune(key))
	trie := NewFlashKeywords(true)
	trie.AddKeyWord("Chetoos", "Cat")
	text := "With Chetoos in place"
	newText := trie.Replace(text)
	rText := "With Cat in place"
	if newText != rText {
		t.Errorf("FAIL %v != %v", newText, rText)
	} else {
		t.Logf("newText: %v", newText)
	}
}

func TestReplaceCleanWordGreaterThenKey(t *testing.T) {
	// len(rune(cleanWord)) > len(rune(key))
	trie := NewFlashKeywords(true)
	trie.AddKeyWord("Cat", "Chetoos")
	text := "With Cat in place"
	newText := trie.Replace(text)
	rText := "With Chetoos in place"
	if newText != rText {
		t.Errorf("FAIL %v != %v", newText, rText)
	} else {
		t.Logf("newText: %v", newText)
	}
}

func TestReplaceCleanWordSameLenghtKey(t *testing.T) {
	// len(rune(cleanWord)) == len(rune(key))
	trie := NewFlashKeywords(true)
	trie.AddKeyWord("Cat", "Bee")
	text := "With Cat in place"
	newText := trie.Replace(text)
	rText := "With Bee in place"
	if newText != rText {
		t.Errorf("FAIL %v != %v", newText, rText)
	} else {
		t.Logf("newText: %v", newText)
	}
}
func TestReplaceKeyAtTheEnd(t *testing.T) {
	trie := NewFlashKeywords(true)
	trie.AddKeyWord("in place", "gone")
	text := "With Cat in place"
	newText := trie.Replace(text)
	rText := "With Cat gone"
	if newText != rText {
		t.Errorf("FAIL %v != %v", newText, rText)
	} else {
		t.Logf("newText: %v", newText)
	}
	trie.AddKeyWord("in place", "in Peeeeeeeeeace")
	newText2 := trie.Replace(text)
	rText2 := "With Cat in Peeeeeeeeeace"
	if newText2 != rText2 {
		t.Errorf("FAIL %v != %v", newText2, rText2)
	} else {
		t.Logf("newText: %v", newText2)
	}
	// random test
	trie2 := NewFlashKeywords(true)
	trie2.addKeyWord("055-5647-3456", "055-XXX")
	text3 := "call chetoos055-5647-3456 chetoosPiza"
	new := trie2.Replace(text3)
	if new != "call chetoos055-XXX chetoosPiza" {
		t.Errorf("FAIL")
	}
	t.Logf("%v", new)

}

func TestRepalceWithFalseCaseSensitive(t *testing.T) {
	trie := NewFlashKeywords(false)
	trie.addKeyWord("Foo", "JOJO")
	text := "HU! foo KIWI"
	newText := trie.Replace(text)
	rText := "hu! jojo kiwi"
	if newText != rText {
		t.Errorf("FAIL %v != %v", newText, rText)
	} else {
		t.Logf("%v", newText)
	}

}
