package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	maxLevel = 32 // max level
)

// SkipList .
type SkipList struct {
	head   *node // head node
	back   *node // back node
	level  int
	length int

	rand *rand.Rand
}

// NewSkipList .
func NewSkipList() *SkipList {

	head := newNode(0, math.MinInt32, maxLevel)

	source := rand.NewSource(time.Now().UnixNano())
	return &SkipList{head: head,
		back:   nil,
		level:  1,
		length: 0,
		rand:   rand.New(source)}
}

// Length .
func (sl *SkipList) Length() int {
	return sl.length
}

// 理论来讲，一级索引中元素个数应该占原始数据的 50%，二级索引中元素个数占 25%，三级索引12.5% ，一直到最顶层。
// 因为这里每一层的晋升概率是 50%。对于每一个新插入的节点，都需要调用 randomLevel 生成一个合理的层数。
// 该 randomLevel 方法会随机生成 1~maxLevel 之间的数，且 ：
//        50%的概率返回 1
//        25%的概率返回 2
//      12.5%的概率返回 3 ...
func (sl *SkipList) randomLevel() int {

	level := 1

	for {
		n := sl.rand.Intn(100)
		if n < 50 && level < maxLevel {
			level++
		} else {
			break
		}
	}

	return level
}

// Add .
func (sl *SkipList) Add(score int64, v interface{}) {
	level := sl.randomLevel()
	sl.add(level, score, v)
}

func (sl *SkipList) addFirst(levelForNewNode int, score int64, v interface{}) {

	newNode := newNode(v, score, levelForNewNode) //创建一个新的跳表节点

	for i := 0; i < levelForNewNode; i++ {
		sl.head.forwards[i] = newNode
	}

	sl.back = newNode
	sl.level = levelForNewNode
	sl.length = 1
}

func (sl *SkipList) add(levelForNewNode int, score int64, v interface{}) {

	if sl.length == 0 {
		sl.addFirst(levelForNewNode, score, v)
		return
	}

	var update [maxLevel]*node //记录每层的路径
	if true {
		cur := sl.head //查找插入位置

		for i := maxLevel - 1; i >= 0; i-- {

			for cur.forwards[i] != nil {
				if score < cur.forwards[i].score {
					// score 小于 cur.forwards[i]，  则 新数据在 cur.forwards[i] 的左边
					update[i] = cur
					break
				}

				// score 大于等于 cur.forwards[i]
				cur = cur.forwards[i]
			}

			if cur.forwards[i] == nil {
				update[i] = cur
			}
		}
	}

	if update[0] != sl.head && update[0].score == score {
		// score对应的节点已存在
		update[0].AddValue(v)
		return
	}

	newNode := newNode(v, score, levelForNewNode) //创建一个新的跳表节点

	// 原有节点连接
	for i := 0; i <= levelForNewNode-1; i++ {
		next := update[i].forwards[i]
		update[i].forwards[i] = newNode
		newNode.forwards[i] = next
	}

	if nextNode := newNode.Next(); nextNode != nil {
		newNode.prev = nextNode.prev
		nextNode.prev = newNode
		// sl.back don't change
	} else {
		newNode.prev = sl.back
		sl.back = newNode
	}

	// 如果当前节点的层数大于之前跳表的层数, 更新当前跳表层数
	if levelForNewNode > sl.level {
		sl.level = levelForNewNode
	}

	sl.length++
	return
}

// GetAll .
func (sl *SkipList) GetAll(isReverse bool) []interface{} {
	if sl.length <= 0 {
		return []interface{}{}
	}

	list := make([]interface{}, 0, 10)
	if isReverse {
		for node := sl.back; node != nil; node = node.Prev() {
			len := len(node.vList)
			if len > 0 {
				for i := len - 1; i >= 0; i-- {
					list = append(list, node.vList[i])
				}
			}
		}
	} else {
		for node := sl.head.forwards[0]; node != nil; node = node.Next() {
			list = append(list, node.vList...)
		}
	}

	return list
}

// GetByRange .
func (sl *SkipList) GetByRange(beginScore int64, endScore int64) []interface{} {

	if sl.length == 0 {
		return []interface{}{}
	}

	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forwards[i] != nil {

			if beginScore == cur.forwards[i].score {
				break
			} else if beginScore < cur.forwards[i].score {
				// 要查的数据 在 cur.forwards[i] 的 左边
				break
			} else if beginScore > cur.forwards[i].score {
				// 要查的数据 在 cur.forwards[i] 的 右边
				cur = cur.forwards[i]
			}
		}
	}

	list := make([]interface{}, 0, 10)
	if cur != sl.head && cur.score == beginScore {
		list = append(list, cur.vList...)
	}

	cur = cur.Next()
	for ; cur != nil; cur = cur.Next() {
		if cur.score <= endScore {
			list = append(list, cur.vList...)
		} else {
			break
		}
	}

	return list
}

func (sl *SkipList) getNode(score int64) *node {
	if sl.length == 0 {
		return nil
	}

	cur := sl.head
	for i := maxLevel - 1; i >= 0; i-- {

		for cur.forwards[i] != nil {
			if score == cur.forwards[i].score {
				return cur.forwards[i]
			} else if score < cur.forwards[i].score {
				// score 小于 cur.forwards[i].score， 则 待查数据 处于 cur.forwards[i] 的左边
				break
			} else if score > cur.forwards[i].score {
				// score 大于 cur.forwards[i].score， 则 待查数据 处于 cur.forwards[i] 的右边
				cur = cur.forwards[i]
			}
		}
	}

	return nil
}

// GetByScore .
func (sl *SkipList) GetByScore(score int64) []interface{} {

	existNode := sl.getNode(score)
	if existNode == nil {
		return []interface{}{}
	}

	len := len(existNode.vList)
	list := make([]interface{}, len, len)
	for index, v := range existNode.vList {
		list[index] = v
	}

	return list
}

// RemoveFirst .
func (sl *SkipList) RemoveFirst() {

	if sl.length == 0 {
		return
	} else if sl.length == 1 {
		for i := sl.level - 1; i >= 0; i-- {
			sl.head.forwards[i] = nil
		}

		sl.back.reset() // you must reset，otherwise resource will be leaked
		sl.back = nil
		sl.level = 1
		sl.length = 0
		return
	}

	sl.Remove(sl.head.forwards[0].score)
	return
}

// RemoveBack .
func (sl *SkipList) RemoveBack() {

	if sl.length == 0 {
		return
	} else if sl.length == 1 { // 只有一个节点
		sl.RemoveFirst()
		return
	}

	sl.Remove(sl.back.score)
	return
}

// Remove .
func (sl *SkipList) Remove(score int64) {

	if sl.length == 0 {
		return
	} else if sl.length == 1 { // 只有一个节点
		if sl.head.forwards[0].score == score {
			sl.RemoveFirst()
		}

		return
	}

	var update [maxLevel]*node //记录每层的路径
	if true {
		cur := sl.head //查找插入位置
		for i := sl.level - 1; i >= 0; i-- {
			update[i] = cur
			for cur.forwards[i] != nil {
				if score == cur.forwards[i].score {
					update[i] = cur
					break
				}

				cur = cur.forwards[i]
			}
		}
	}

	var delNode *node
	if update[0].forwards[0] != nil && update[0].forwards[0].score == score {
		delNode = update[0].forwards[0]
	} else {
		return
	}

	for i := delNode.level - 1; i >= 0; i-- {
		if update[i] == sl.head && delNode.forwards[i] == nil {
			sl.level = i // 没有节点比我更高，level降一级
		}

		if update[i].forwards[i] != nil {
			update[i].forwards[i] = update[i].forwards[i].forwards[i]
		}
	}

	// adjust pre node
	if true {
		if nextNode := delNode.Next(); nextNode != nil {
			nextNode.prev = delNode.prev
		}
	}

	// adjust back node
	if delNode == sl.back {
		sl.back = sl.back.prev
	}

	sl.length--
	delNode.reset() // you must reset，otherwise resource will be leaked

	return
}

func (sl *SkipList) String() string {
	return fmt.Sprintf("SkipList level:%+v, length:%+v", sl.level, sl.length)
}
