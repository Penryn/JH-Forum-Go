package node

type TrieNode struct {
	Index    int
	Layer    int
	End      bool
	Char     int
	Results  []int
	M_values map[int]*TrieNode
	Failure  *TrieNode
	Parent   *TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		End:      false,
		M_values: make(map[int]*TrieNode),
		Results:  make([]int, 0),
	}
}

func (s *TrieNode) Add(c int) *TrieNode {
	if val, s := s.M_values[c]; s {
		return val
	}
	node := NewTrieNode()
	node.Parent = s
	node.Char = c
	s.M_values[c] = node
	return node
}

func (s *TrieNode) SetResults(text int) {
	if !s.End {
		s.End = true
	}
	for i := 0; i < len(s.Results); i++ {
		if s.Results[i] == text {
			return
		}
	}
	s.Results = append(s.Results, text)
}
