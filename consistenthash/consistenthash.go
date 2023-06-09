package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
// 定义了函数类型 Hash，采取依赖注入的方式，允许用于替换成自定义的 Hash 函数，
// 也方便测试时替换，默认为 crc32.ChecksumIEEE 算法。
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
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

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
