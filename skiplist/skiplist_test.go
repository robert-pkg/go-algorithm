package skiplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	assert := assert.New(t)

	sl := NewSkipList()

	sl.add(10, 10, "tom")
	sl.add(2, 20, "charles")
	sl.add(5, 15, "robert")
	sl.add(10, 10, "candy")

	if true {
		tests := []struct {
			score int64
			value string
		}{
			{10, "tom"},
			{10, "candy"},
			{15, "robert"},
			{20, "charles"},
		}

		list := sl.GetAll(false)

		assert.Equal(len(tests), len(list), "length should be equal")
		for i := 0; i < len(list); i++ {
			assert.Equal(tests[i].value, list[i], "value should be equal")
		}
	}

	if true {
		tests := []struct {
			score int64
			value string
		}{
			{20, "charles"},
			{15, "robert"},
			{10, "candy"},
			{10, "tom"},
		}

		list := sl.GetAll(true)

		assert.Equal(len(tests), len(list), "length should be equal")
		for i := 0; i < len(list); i++ {
			assert.Equal(tests[i].value, list[i], "value should be equal")
		}
	}

}

func Test_GetByScore(t *testing.T) {
	assert := assert.New(t)

	sl := NewSkipList()

	sl.add(10, 10, "tom")
	sl.add(2, 20, "charles")
	sl.add(5, 15, "robert")
	sl.add(10, 10, "candy")

	tests := []struct {
		score int64
		value string
	}{
		{10, "tom"},
		{10, "candy"},
	}

	list := sl.GetByScore(10)

	assert.Equal(len(tests), len(list), "length should be equal")
	for i := 0; i < len(list); i++ {
		assert.Equal(tests[i].value, list[i], "value should be equal")
	}
}

func Test_RemoveFirst(t *testing.T) {

	assert := assert.New(t)

	sl := NewSkipList()

	sl.add(10, 10, "tom")
	sl.add(5, 15, "robert")
	sl.add(2, 20, "charles")

	sl.RemoveFirst()

	assert.Equal(5, sl.level, "skiplist level should be equal")
	assert.Equal(2, sl.length, "skiplist length should be equal")

	t.Logf("skip list: %v", sl)

	t.Log("遍历:")
	node := sl.head.forwards[0]
	for node != nil {
		t.Logf("node: %v", node)
		node = node.Next()
	}
}

func Test_RemoveBack(t *testing.T) {

	assert := assert.New(t)

	tests := []struct {
		sl     func() *SkipList
		level  int
		length int
	}{
		{sl: func() *SkipList {
			list := NewSkipList()
			list.add(10, 10, "tom")
			list.add(5, 15, "robert")
			list.add(2, 20, "charles") // 该节点将被删除
			return list
		},
			level:  10,
			length: 2,
		},
	}

	for _, r := range tests {
		sl := r.sl()
		sl.RemoveBack()

		assert.Equal(r.level, sl.level, "skiplist level should be equal")
		assert.Equal(r.length, sl.length, "skiplist length should be equal")

		t.Logf("skip list: %v", sl)

		t.Log("遍历:")
		node := sl.head.forwards[0]
		for node != nil {
			t.Logf("node: %v", node)
			node = node.Next()
		}
	}
}

func Test_Remove(t *testing.T) {
}
