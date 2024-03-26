package loadbalance

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// ConsistentHashBalance 一致性hash算法实现负载均衡。
// 用于解决简单哈希算法在增删节点时，重新映射带来的效率低下问题。
// 	一致性/单调性：以uint32范围作为哈希表，新增或者删减节点时，
//	  【不影响】系统正常运行，解决哈希表的动态伸缩问题
// 	分散性：数据应该分散地存放在（分布式集群中的）各个节点，不必每个节点都存储所有的数据
// 	平衡性：采用虚拟节点解决hash环偏斜问题。
//	  hash的结果应该平均分配到各个节点，从算法层面解决负载均衡问题
// 实现步骤：
// 	1.计算存储节点（服务器）哈希值，将其存储空间抽象成一个环（0 - 2^32 -1）
// 	2.对数据（URL、IP）进行哈希计算，按顺时针方向将其映射到距离最近的节点上

type ConsistentHashBalance struct {
	// hash 函数,支持自定义,默认使用crc32.ChecksumIEEE
	hash Hash

	// 服务器节点hash列表,按照从小到大排序
	HashKeys UInt32Slice

	// 服务器节点 host值与服务器真实地址映射表
	hashMap map[uint32]string

	// 虚拟节点倍数
	// 解决平衡数问题
	replicas int
	// 读写锁
	mux sync.RWMutex
}

// Hash 默认使用 crc32.CheckIEEE
type Hash func(data []byte) uint32

type UInt32Slice []uint32

// NewConsistentHashBalance 初始化一致性Hash结构体
func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	consistentHashBalance := &ConsistentHashBalance{
		hash:     fn,
		hashMap:  make(map[uint32]string),
		replicas: replicas,
	}
	if consistentHashBalance.hash == nil {
		// 返回一个32位无符号整数
		consistentHashBalance.hash = crc32.ChecksumIEEE
	}
	return consistentHashBalance
}

// Add 添加服务器节点,参数为服务器地址
// eg. "http://ip:port/
// 对每一个真实的addr,对应创建 replcas 虚拟节点
// 最后,对 hashkey进行排序
func (c *ConsistentHashBalance) Add(servers ...string) error {
	if len(servers) == 0 {
		return errors.New("servers length at least 1")
	}

	c.mux.Lock()
	defer c.mux.Unlock()
	for _, addr := range servers {
		// 对每个addr创建replicas虚拟节点
		for i := 0; i < c.replicas; i++ {
			// 算出的hash都是一样的,没有意义,所以要加i 以便于进行区分
			//hash := c.hash([]byte(addr))
			hash := c.hash([]byte(strconv.Itoa(i) + addr))
			c.HashKeys = append(c.HashKeys, hash) // 将hash值追加到hash列表中
			// 实现一个 addr 对应多个 hash值 , key: hash值, value: addr
			c.hashMap[hash] = addr
		}
	}
	sort.Sort(c.HashKeys)
	//sort.Slice(c.hashKeys, func(i, j int) bool {
	//	return c.hashKeys[i] < c.hashKeys[j]
	//})
	return nil
}

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Get 获取指定的key最靠近它的那个服务器节点
// 返回服务器节点的hash值 >= key的hash值 (也可能穿过起止环)

// 实现步骤:
// 1.计算key的hash值
// 2.通过二分查找最有服务器节点
// 3.取出服务器
func (c *ConsistentHashBalance) Get(key string) (string, error) {
	length := len(c.HashKeys)
	if length == 0 {
		return "", errors.New("node list is Empty")
	}
	// 1.计算key的hash值
	hash := c.hash([]byte(key))
	fmt.Println(hash)

	// 2.通过二分查找最有服务器节点
	index := sort.Search(length, func(i int) bool {
		return c.HashKeys[i] >= hash
	})
	// 如果返回的index==lenght,就意味着没有查到,所以直接返回第一个节点
	if index == length {
		index = 0
	}

	// 3.取出服务地址,并返回
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.hashMap[c.HashKeys[index]], nil
}
