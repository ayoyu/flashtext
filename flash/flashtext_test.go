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
	keys2Cleans := map[string][]string{
		"Earth":     {"planet"},
		"Mars":      {"planet"},
		"sun":       {"star", "Sun"},
		"python3.5": {"python"},
		"python3":   {"python"},
	}
	for key, data := range keys2Cleans {
		for _, clean := range data {
			trie.AddKeyWord(key, clean)
		}
	}
	if trie.Size() != len(keys2Cleans) {
		t.Errorf("trie.AddKeyWord FAILED, Size is expected to be equal: %d", len(keys2Cleans))
	}

	allKeysWords := trie.GetAllKeywords()

	for key, data := range keys2Cleans {
		setCleans, ok := allKeysWords[key]
		if !ok {
			t.Errorf("the Key %v does not exist in the tree: ", key)
		}
		for _, item := range data {
			if _, ok := setCleans[item]; !ok {
				t.Errorf("Clean Name %v does not exist for the key %v", item, key)
			}
		}
		t.Logf("Key %v with setCleanNames %v PASS", key, setCleans)
	}

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
		for i := 0; i < 2; i++ {
			if res[i].key == keys[i] && text[res[i].start:res[i].end+1] == keys[i] {
				t.Logf("Found Key %v, PASS", keys[i])
			} else {
				t.Errorf("Error Key %v FAILED %v", keys[i], res[i])
			}
		}
	}

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
		t.Logf("Modulo %v", 2%len(keys))
		for i := 0; i < 2; i++ {
			if res[i].key == keys[(i+1)%len(keys)] && text[res[i].start:res[i].end+1] == keys[(i+1)%len(keys)] {
				t.Logf("Found Key %v, PASS", keys[(i+1)%len(keys)])
			} else {
				t.Errorf("Error Key %v FAILED %v", keys[(i+1)%len(keys)], res[i])
			}
		}
	}
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
		if res[i].key != strings.ToLower(keys[i]) {
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
	if res[0].key != keys[0] {
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
		if keys[i] != res[i].key {
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
