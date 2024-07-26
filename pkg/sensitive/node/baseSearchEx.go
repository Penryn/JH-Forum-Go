package node

import "sort"
import (
	"bytes"
	"encoding/binary"
	"os"
)

type BaseSearchEx struct {
	I_keywords []string
	I_guides   [][]int
	I_key      []int
	I_next     []int
	I_check    []int
	I_dict     []int
}

func (s *BaseSearchEx) checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func (s *BaseSearchEx) Save(filename string) {
	f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	s.Save2(f)
	f.Close()
}

func (s *BaseSearchEx) Save2(f *os.File) {
	f.Write(s.intToBytes(len(s.I_keywords)))
	for _, key := range s.I_keywords {
		f.Write(s.intToBytes(len(key)))
		f.Write([]byte(key))
	}

	f.Write(s.intToBytes(len(s.I_guides)))
	for _, key := range s.I_guides {
		f.Write(s.intToBytes(len(key)))
		for _, item := range key {
			f.Write(s.intToBytes(item))
		}
	}

	f.Write(s.intToBytes(len(s.I_key)))
	for _, item := range s.I_key {
		f.Write(s.intToBytes(item))
	}

	f.Write(s.intToBytes(len(s.I_next)))
	for _, item := range s.I_next {
		f.Write(s.intToBytes(item))
	}

	f.Write(s.intToBytes(len(s.I_check)))
	for _, item := range s.I_check {
		f.Write(s.intToBytes(item))
	}

	f.Write(s.intToBytes(len(s.I_dict)))
	for _, item := range s.I_dict {
		f.Write(s.intToBytes(item))
	}
}

func (s *BaseSearchEx) Load(filename string) {
	f, _ := os.OpenFile(filename, os.O_RDONLY, 0666)
	s.Load2(f)
	f.Close()
}
func (s *BaseSearchEx) Load2(f *os.File) {
	intBs := make([]byte, 4)

	f.Read(intBs)
	length := s.bytesToInt(intBs)

	s.I_keywords = make([]string, length)
	for i := 0; i < length; i++ {
		f.Read(intBs)
		l := s.bytesToInt(intBs)
		temp := make([]byte, l)
		f.Read(temp)
		s.I_keywords[i] = string(temp)
	}

	f.Read(intBs)
	length = s.bytesToInt(intBs)
	s.I_guides = make([][]int, length)
	for i := 0; i < length; i++ {
		f.Read(intBs)
		l := s.bytesToInt(intBs)
		ls := make([]int, l)
		for j := 0; j < l; j++ {
			f.Read(intBs)
			ls[j] = s.bytesToInt(intBs)
		}
		s.I_guides[i] = ls
	}

	f.Read(intBs)
	length = s.bytesToInt(intBs)
	s.I_key = make([]int, length)
	for i := 0; i < length; i++ {
		f.Read(intBs)
		s.I_key[i] = s.bytesToInt(intBs)
	}

	f.Read(intBs)
	length = s.bytesToInt(intBs)
	s.I_next = make([]int, length)
	for i := 0; i < length; i++ {
		f.Read(intBs)
		s.I_next[i] = s.bytesToInt(intBs)
	}

	f.Read(intBs)
	length = s.bytesToInt(intBs)
	s.I_check = make([]int, length)
	for i := 0; i < length; i++ {
		f.Read(intBs)
		s.I_check[i] = s.bytesToInt(intBs)
	}

	f.Read(intBs)
	length = s.bytesToInt(intBs)
	s.I_dict = make([]int, length)
	for i := 0; i < length; i++ {
		f.Read(intBs)
		s.I_dict[i] = s.bytesToInt(intBs)
	}
}

func (s *BaseSearchEx) intToBytes(i int) []byte {
	x := int32(i)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
func (s *BaseSearchEx) bytesToInt(bs []byte) int {
	bytesBuffer := bytes.NewBuffer(bs)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func (s *BaseSearchEx) SetKeywords(keywords []string) {
	s.I_keywords = keywords
	length := s.CreateDict(keywords)
	root := NewTrieNodeEx()

	for i, keyword := range keywords {
		nd := root
		p := []rune(keyword)
		for _, c := range string(p) {
			nd = nd.Add(int32(s.I_dict[c]))
		}
		nd.SetResults(i)
	}
	nodes := make([]*TrieNodeEx, 0)
	for _, value := range root.M_values {
		value.Failure = root
		for _, trans := range value.M_values {
			nodes = append(nodes, trans)
		}
	}
	for len(nodes) > 0 {
		newNodes := make([]*TrieNodeEx, 0)
		for _, nd := range nodes {
			r := nd.Parent.Failure
			c := nd.Char
			for {
				if r == nil {
					break
				}
				if _, s := r.M_values[c]; s {
					break
				}
				r = r.Failure
			}
			if r == nil {
				nd.Failure = root
			} else {
				nd.Failure = r.M_values[c]
				for _, result := range nd.Failure.Results {
					nd.SetResults(result)
				}
			}

			for _, child := range nd.M_values {
				newNodes = append(newNodes, child)
			}
		}
		nodes = newNodes
	}
	root.Failure = root
	for _, item := range root.M_values {
		s.tryLinks(item)
	}
	s.build(root, length)
}

func (s *BaseSearchEx) tryLinks(node *TrieNodeEx) {
	node.Merge(node.Failure)
	for _, item := range node.M_values {
		s.tryLinks(item)
	}
}

func (s *BaseSearchEx) build(root *TrieNodeEx, length int) {
	has := make([]*TrieNodeEx, 0x00FFFFFF)

	length = root.Rank(has) + length + 1
	s.I_key = make([]int, length)
	s.I_next = make([]int, length)
	s.I_check = make([]int, length)
	var guides [][]int
	first := make([]int, 1)
	first[0] = 0
	guides = append(guides, first)

	for i := 0; i < length; i++ {
		item := has[i]
		if item == nil {
			continue
		}
		s.I_key[i] = int(item.Char)
		s.I_next[i] = item.Next
		if item.End {
			s.I_check[i] = len(guides)
			guides = append(guides, item.Results)
		}
	}
	s.I_guides = guides
}

func (s *BaseSearchEx) CreateDict(keywords []string) int {
	dictionary := make(map[int32]int, 0)

	for _, keyword := range keywords {
		for _, item := range keyword {
			if v, s := dictionary[item]; s {
				if v > 0 {
					dictionary[item] = dictionary[item] + 2
				}
			} else {
				if v > 0 {
					dictionary[item] = 2
				} else {
					dictionary[item] = 1
				}
			}
		}
	}
	list := s.sortMap(dictionary)

	index1 := make([]int, 0)
	for i := 0; i < len(list); i = i + 2 {
		index1 = append(index1, i)
	}
	length := len(index1)
	for i := 0; i < length/2; i++ {
		index1[i], index1[length-i-1] = index1[length-i-1], index1[i]
	}
	for i := 1; i < len(list); i = i + 2 {
		index1 = append(index1, i)
	}

	list2 := make([]int32, 0)
	for i := 0; i < len(list); i++ {
		list2 = append(list2, list[index1[i]])
	}

	s.I_dict = make([]int, 0x10000)
	for i, v := range list2 {
		s.I_dict[v] = i + 1
	}
	return len(dictionary)
}
func (s *BaseSearchEx) sortMap(mp map[int32]int) []int32 {
	var newMp = make([]int, 0)
	var newMpKey = make([]int32, 0)
	for oldk, v := range mp {
		newMp = append(newMp, v)
		newMpKey = append(newMpKey, oldk)
	}
	sort.Ints(newMp)

	list := make([]int32, 0)
	for k := range newMp {
		list = append(list, newMpKey[k])
	}
	return list
}
