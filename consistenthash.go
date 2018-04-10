package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash hash方法
type Hash func(data []byte) uint32

// Map map
type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

// New 初始化Map
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty 检查是否为空
func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// Add 添加
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 获取
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	if idx == len(m.keys) {
		idx = 0
	}
	return m.hashMap[m.keys[idx]]
}
