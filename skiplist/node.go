package skiplist

import "fmt"

func newNode(v interface{}, score int64, level int) *node {
	return &node{
		vList:    []interface{}{v},
		score:    score,
		level:    level,
		forwards: make([]*node, level, level),
		prev:     nil,
	}
}

// 跳表节点
type node struct {
	vList    []interface{} //跳表保存的值
	score    int64         //用于排序的分值
	level    int           //层高
	forwards []*node       //每层前进指针

	prev *node // 上一个节点
}

func (node *node) AddValue(v interface{}) {
	node.vList = append(node.vList, v)
}

func (node *node) Prev() *node {
	return node.prev
}

// Next
func (node *node) Next() *node {
	if len(node.forwards) == 0 {
		return nil
	}

	return node.forwards[0]
}

func (node *node) reset() {
	node.vList = nil
	node.forwards = nil
	node.prev = nil
}

func (node *node) String() string {
	return fmt.Sprintf("skipListNode v:%+v, score:%+v, level:%+v, forwards length:%+v", node.vList, node.score, node.level, len(node.forwards))
}
