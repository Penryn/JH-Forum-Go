package node

type TrieNodeEx struct {
	Parent       *TrieNodeEx
	Failure      *TrieNodeEx
	Char         int32
	End          bool
	Results      []int
	M_values     map[int32]*TrieNodeEx
	Merge_values map[int32]*TrieNodeEx
	minflag      int32
	maxflag      int32
	Next         int
	Count        int
}

func NewTrieNodeEx() *TrieNodeEx {
	return &TrieNodeEx{
		M_values:     make(map[int32]*TrieNodeEx, 0),
		Merge_values: make(map[int32]*TrieNodeEx, 0),
		Results:      make([]int, 0),
		minflag:      0xffff,
		maxflag:      0,
		Next:         0,
		Count:        0,
	}
}

func (ts *TrieNodeEx) TryGetValue(c int32) (bool, *TrieNodeEx) {
	if ts.minflag <= c && ts.maxflag >= c {
		if val, s := ts.M_values[c]; s {
			return true, val
		}
	}
	return false, nil
}

func (ts *TrieNodeEx) Add(c int32) *TrieNodeEx {
	if val, s := ts.M_values[c]; s {
		return val
	}
	if ts.minflag > c {
		ts.minflag = c
	}
	if ts.maxflag < c {
		ts.maxflag = c
	}
	node := NewTrieNodeEx()
	node.Parent = ts
	node.Char = c
	ts.M_values[c] = node
	ts.Count++
	return node
}

func (ts *TrieNodeEx) SetResults(text int) {
	if !ts.End {
		ts.End = true
	}
	for i := 0; i < len(ts.Results); i++ {
		if ts.Results[i] == text {
			return
		}
	}
	ts.Results = append(ts.Results, text)
}

func (ts *TrieNodeEx) Merge(node *TrieNodeEx) {
	nd := node
	for nd.Char != 0 {
		for key, value := range node.M_values {
			if _, s := ts.M_values[key]; s {
				continue
			}
			if _, s := ts.Merge_values[key]; s {
				continue
			}
			if ts.minflag > key {
				ts.minflag = key
			}
			if ts.maxflag < key {
				ts.maxflag = key
			}
			ts.Merge_values[key] = value
			ts.Count++
		}
		nd = nd.Failure
	}
}

func (ts *TrieNodeEx) Rank(has []*TrieNodeEx) int {
	var seats []bool = make([]bool, len(has))
	var start int = 1
	has[0] = ts
	ts.Rank2(start, seats, has)
	maxCount := len(has) - 1
	for has[maxCount] == nil {
		maxCount--
	}
	return maxCount
}

func (ts *TrieNodeEx) Rank2(start int, seats []bool, has []*TrieNodeEx) int {
	if ts.maxflag == 0 {
		return start
	}
	keys := make([]int32, 0)
	for k := range ts.M_values {
		keys = append(keys, k)
	}
	for k := range ts.Merge_values {
		keys = append(keys, k)
	}

	for has[start] != nil {
		start++
	}
	s := start
	if start < int(ts.minflag) {
		s = int(ts.minflag)
	}

	for i := s; i < len(has); i++ {
		if has[i] == nil {
			next := i - int(ts.minflag)
			if seats[next] {
				continue
			}
			isok := true

			for _, item := range keys {
				if has[next+int(item)] != nil {
					isok = false
					break
				}
			}
			if isok {
				ts.SetSeats(next, seats, has)
				break
			}
		}
	}
	start += len(keys) / 2

	for _, value := range ts.M_values {
		start = value.Rank2(start, seats, has)
	}
	return start
}

func (ts *TrieNodeEx) SetSeats(next int, seats []bool, has []*TrieNodeEx) {
	ts.Next = next
	seats[next] = true
	for key, value := range ts.Merge_values {
		position := next + int(key)
		has[position] = value
	}
	for key, value := range ts.M_values {
		position := next + int(key)
		has[position] = value
	}
}
