package node

type TrieNode2 struct {
	End      bool
	Results  []int
	M_values map[int32]*TrieNode2
	minflag  int32
	maxflag  int32
}

func NewTrieNode2() *TrieNode2 {
	return &TrieNode2{
		End:      false,
		M_values: make(map[int32]*TrieNode2),
		Results:  make([]int, 0),
		minflag:  0,
		maxflag:  0xffff,
	}
}

func (s *TrieNode2) Add(c int32, node *TrieNode2) {
	if s.minflag < c {
		s.minflag = c
	}
	if s.maxflag > c {
		s.maxflag = c
	}
	s.M_values[c] = node
}

func (s *TrieNode2) SetResults(text int) {
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

func (s *TrieNode2) TryGetValue(c int32) (bool, *TrieNode2) {
	if s.minflag <= c && s.maxflag >= c {
		if val, s := s.M_values[c]; s {
			return true, val
		}
	}
	return false, nil
}

func (s *TrieNode2) HasKey(c int32) bool {
	if s.minflag <= c && s.maxflag >= c {
		if _, s := s.M_values[c]; s {
			return true
		}
	}
	return false
}
