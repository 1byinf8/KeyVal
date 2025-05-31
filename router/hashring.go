package router

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type HashRing struct {
	replicas   int
	hashCircle map[uint32]string
	sortedKeys []uint32
	nodes      map[string]bool
	mu         sync.RWMutex
}

func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		replicas:   replicas,
		hashCircle: make(map[uint32]string),
		nodes:      make(map[string]bool),
	}
}

func HashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (h *HashRing) AddNode(nodeID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.nodes[nodeID] {
		return
	}

	h.nodes[nodeID] = true
	for i := 0; i < h.replicas; i++ {
		virtualNode := nodeID + "#" + strconv.Itoa(i)
		hash := HashKey(virtualNode)
		h.hashCircle[hash] = virtualNode
		h.sortedKeys = append(h.sortedKeys, hash)
	}
	sort.Slice(h.sortedKeys, func(i, j int) bool {
		return h.sortedKeys[i] < h.sortedKeys[j]
	})
}

func (h *HashRing) RemoveNode(nodeID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.nodes[nodeID] {
		return
	}

	delete(h.nodes, nodeID)
	newSortedKeys := []uint32{}

	for _, hash := range h.sortedKeys {
		if h.hashCircle[hash] != nodeID {
			newSortedKeys = append(newSortedKeys, hash)
		} else {
			delete(h.hashCircle, hash)
		}
	}
	h.sortedKeys = newSortedKeys
}

func (h *HashRing) GetNodes(key string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.sortedKeys) == 0 {
		return nil
	}

	hash := HashKey(key)
	result := make([]string, 0, 3)
	seen := make(map[string]bool)

	startIdx := sort.Search(len(h.sortedKeys), func(i int) bool {
		return h.sortedKeys[i] >= hash
	})

	if startIdx == len(h.sortedKeys) {
		startIdx = 0
	}

	for i := 0; i < len(h.sortedKeys) && len(result) < 3; i++ {
		idx := (startIdx + i) % len(h.sortedKeys)
		node := h.hashCircle[h.sortedKeys[idx]]

		if !seen[node] {
			result = append(result, node)
			seen[node] = true
		}
	}

	return result
}

func (h *HashRing) GetAllNodes() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	nodes := make([]string, 0, len(h.nodes))
	for node := range h.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (h *HashRing) IsEmpty() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.nodes) == 0
}
