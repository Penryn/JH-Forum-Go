package words

import (
	"bytes"
	"encoding/binary"
	"errors"
	"JH-Forum/pkg/sensitive/node"
	"os"
)

type IllegalWordsSearch struct {
	node.BaseSearchEx
	UseSkipWordFilter             bool
	_skipBitArray                 []bool
	UseDuplicateWordFilter        bool
	UseBlacklistFilter            bool
	_blacklist                    []int
	UseDBCcaseConverter           bool
	UseSimplifiedChineseConverter bool
	UseIgnoreCase                 bool
}

func NewIllegalWordsSearch() *IllegalWordsSearch {
	_skipList := " \t\r\n~!@#$%^&*()_+-=【】、[]{}|;':\"，。、《》？αβγδεζηθικλμνξοπρστυφχψωΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ。，、；：？！…—·ˉ¨‘’“”々～‖∶＂＇｀｜〃〔〕〈〉《》「」『』．〖〗【】（）［］｛｝ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫ⒈⒉⒊⒋⒌⒍⒎⒏⒐⒑⒒⒓⒔⒕⒖⒗⒘⒙⒚⒛㈠㈡㈢㈣㈤㈥㈦㈧㈨㈩①②③④⑤⑥⑦⑧⑨⑩⑴⑵⑶⑷⑸⑹⑺⑻⑼⑽⑾⑿⒀⒁⒂⒃⒄⒅⒆⒇≈≡≠＝≤≥＜＞≮≯∷±＋－×÷／∫∮∝∞∧∨∑∏∪∩∈∵∴⊥∥∠⌒⊙≌∽√§№☆★○●◎◇◆□℃‰€■△▲※→←↑↓〓¤°＃＆＠＼︿＿￣―♂♀┌┍┎┐┑┒┓─┄┈├┝┞┟┠┡┢┣│┆┊┬┭┮┯┰┱┲┳┼┽┾┿╀╁╂╃└┕┖┗┘┙┚┛━┅┉┤┥┦┧┨┩┪┫┃┇┋┴┵┶┷┸┹┺┻╋╊╉╈╇╆╅╄"
	_skipBitArray := make([]bool, 0x10000)
	for _, val := range _skipList {
		_skipBitArray[int(val)] = true
	}

	return &IllegalWordsSearch{
		UseSkipWordFilter:             true,
		_skipBitArray:                 _skipBitArray,
		UseDuplicateWordFilter:        false,
		UseBlacklistFilter:            false,
		_blacklist:                    make([]int, 0),
		UseDBCcaseConverter:           true,
		UseSimplifiedChineseConverter: true,
		UseIgnoreCase:                 true,
	}
}

// 在文本中查找所有的关键字
func (s *IllegalWordsSearch) FindAll(text string) []*IllegalWordsSearchResult {
	return s.FindAll2(text, 0xffffffff)
}

func (s *IllegalWordsSearch) FindAll2(text string, flag int) []*IllegalWordsSearchResult {
	results := make([]*IllegalWordsSearchResult, 0)
	pIndex := make([]int, len(text))
	p := 0
	findIndex := 0
	var pChar int32 = 0

	var i int
	for _, c := range text {
		if p != 0 {
			pIndex[i] = p
			if findIndex != 0 {
				for _, item := range s.I_guides[findIndex] {
					r := s.getIllegalResult(item, i-1, text, p, pIndex, flag)
					if r != nil {
						results = append(results, r)
					}
				}
			}
		}
		if s.UseSkipWordFilter && s._skipBitArray[c] { //使用跳词
			findIndex = 0
			i++
			continue
		}
		t := s.I_dict[c]
		if t == 0 { //不在字表中，跳过
			p = 0
			pChar = c
			i++
			continue
		}
		next := s.I_next[p] + t
		find := s.I_key[next] == t
		if !find {
			if s.UseDuplicateWordFilter && pChar == c {
				i++
				continue
			}
			if p != 0 {
				p = 0
				next = s.I_next[0] + t
				find = s.I_key[next] == t
			}
		}
		if find {
			findIndex = s.I_check[next]
			p = next
		}
		pChar = c
		i++
	}
	if findIndex != 0 {
		for _, item := range s.I_guides[findIndex] {
			r := s.getIllegalResult(item, i-1, text, p, pIndex, flag)
			if r != nil {
				results = append(results, r)
			}
		}
	}
	return results
}

func (s *IllegalWordsSearch) FindFirst(text string) *IllegalWordsSearchResult {
	return s.FindFirst2(text, 0xffffffff)
}

func (s *IllegalWordsSearch) FindFirst2(text string, flag int) *IllegalWordsSearchResult {
	pIndex := make([]int, len(text))
	p := 0
	findIndex := 0
	var pChar int32 = 0

	var i int
	for _, c := range text {
		if p != 0 {
			pIndex[i] = p
			if findIndex != 0 {
				for _, item := range s.I_guides[findIndex] {
					r := s.getIllegalResult(item, i-1, text, p, pIndex, flag)
					if r != nil {
						return r
					}
				}
			}
		}
		if s.UseSkipWordFilter && s._skipBitArray[c] { //使用跳词
			findIndex = 0
			i++
			continue
		}
		t := s.I_dict[c]
		if t == 0 { //不在字表中，跳过
			p = 0
			pChar = c
			i++
			continue
		}
		next := s.I_next[p] + t
		find := s.I_key[next] == t
		if !find {
			if s.UseDuplicateWordFilter && pChar == c {
				i++
				continue
			}
			if p != 0 {
				p = 0
				next = s.I_next[0] + t
				find = s.I_key[next] == t
			}
		}
		if find {
			findIndex = s.I_check[next]
			p = next
		}
		pChar = c
		i++
	}
	if findIndex != 0 {
		for _, item := range s.I_guides[findIndex] {
			r := s.getIllegalResult(item, i-1, text, p, pIndex, flag)
			if r != nil {
				return r
			}
		}
	}
	return nil
}

func (s *IllegalWordsSearch) ContainsAny(text string) bool {
	return s.ContainsAny2(text, 0xffffffff)
}

func (s *IllegalWordsSearch) ContainsAny2(text string, flag int) bool {
	pIndex := make([]int, len(text))
	p := 0
	findIndex := 0
	var pChar int32 = 0

	var i int = 0
	for _, c := range text {
		if p != 0 {
			pIndex[i] = p
			if findIndex != 0 {
				for _, item := range s.I_guides[findIndex] {
					r := s.getIllegalResult(item, i-1, text, p, pIndex, flag)
					if r != nil {
						return true
					}
				}
			}
		}
		if s.UseSkipWordFilter && s._skipBitArray[c] { //使用跳词
			findIndex = 0
			i++
			continue
		}
		t := s.I_dict[c]
		if t == 0 { //不在字表中，跳过
			p = 0
			pChar = c
			i++
			continue
		}
		next := s.I_next[p] + t
		find := s.I_key[next] == t
		if !find {
			if s.UseDuplicateWordFilter && pChar == c {
				i++
				continue
			}
			if p != 0 {
				p = 0
				next = s.I_next[0] + t
				find = s.I_key[next] == t
			}
		}
		if find {
			findIndex = s.I_check[next]
			p = next
		}
		pChar = c
		i++
	}
	if findIndex != 0 {
		for _, item := range s.I_guides[findIndex] {
			r := s.getIllegalResult(item, i-1, text, p, pIndex, flag)
			if r != nil {
				return true
			}
		}
	}
	return false
}

func (s *IllegalWordsSearch) Replace(text string, replaceChar rune) string {
	return s.Replace2(text, replaceChar, 0xffffffff)
}

func (s *IllegalWordsSearch) Replace2(text string, replaceChar rune, flag int) string {
	result := []rune(text)

	pIndex := make([]int, len(text))
	p := 0
	findIndex := 0
	var pChar int32 = 0

	var i int
	for _, c := range text {
		if p != 0 {
			pIndex[i] = p
			if findIndex != 0 {
				for _, item := range s.I_guides[findIndex] {
					r := s.getIllegalResult(item, i-1, text, p, pIndex, flag)
					if r != nil {
						for j := r.Start; j < i; j++ {
							result[j] = replaceChar
						}
						break
					}
				}
			}
		}
		if s.UseSkipWordFilter && s._skipBitArray[c] { //使用跳词
			findIndex = 0
			i++
			continue
		}
		t := s.I_dict[c]
		if t == 0 { //不在字表中，跳过
			p = 0
			pChar = c
			i++
			continue
		}
		next := s.I_next[p] + t
		find := s.I_key[next] == t
		if !find {
			if s.UseDuplicateWordFilter && pChar == c {
				i++
				continue
			}
			if p != 0 {
				p = 0
				next = s.I_next[0] + t
				find = s.I_key[next] == t
			}
		}
		if find {
			findIndex = s.I_check[next]
			p = next
		}
		pChar = c
		i++
	}
	if findIndex != 0 {
		for _, item := range s.I_guides[findIndex] {
			r := s.getIllegalResult(item, i-1, text, p, pIndex, flag)
			if r != nil {
				for j := r.Start; j < i; j++ {
					result[j] = replaceChar
				}
				break
			}
		}
	}
	return string(result)
}

func (s *IllegalWordsSearch) findStart(keyword string, end int, srcText string, p int, pIndex []int) int {
	_srcText := []rune(srcText)
	if end+1 < len(_srcText) {
		en1 := isEnglishOrNumber(_srcText[end+1])
		en2 := isEnglishOrNumber(_srcText[end])
		if en1 && en2 {
			return -1
		}
	}
	n := len([]rune(keyword))
	start := end
	pp := p
	for n > 0 {
		pi := pIndex[start]
		start--
		if pi != pp {
			n--
			pp = pi
		}
		if start == -1 {
			return 0
		}
	}

	sn1 := isEnglishOrNumber(_srcText[start])
	start++
	sn2 := isEnglishOrNumber(_srcText[start])
	if sn1 && sn2 {
		return -1
	}
	return start
}

func isEnglishOrNumber(c int32) bool {
	if c < 128 {
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			return true
		}
	}
	return false
}

func (s *IllegalWordsSearch) getIllegalResult(index int, end int, srcText string, p int, pIndex []int, flag int) *IllegalWordsSearchResult {
	if s.UseBlacklistFilter {
		b := s._blacklist[index]
		if (b | flag) != b {
			return nil
		}
	}
	keyword := s.I_keywords[index]
	_keyword := []rune(keyword)
	_srcText := []rune(srcText)
	if len(_keyword) == 1 {
		if !(s.toSenseWord(_srcText[end]) == s.toSenseWord(_keyword[0])) {
			return nil
		}
		return NewIllegalWordsSearchResult(keyword, end, end, srcText, 0xffffffff)
	}
	start := s.findStart(keyword, end, srcText, p, pIndex)
	if start == -1 {
		return nil
	}
	if !(s.toSenseWord(_srcText[start]) == s.toSenseWord(_keyword[0])) {
		return nil
	}
	if s.UseBlacklistFilter {
		return NewIllegalWordsSearchResult(keyword, start, end, srcText, s._blacklist[index])
	}
	return NewIllegalWordsSearchResult(keyword, start, end, srcText, 0xffffffff)
}
func (s *IllegalWordsSearch) toSenseWord(c int32) int32 {
	if s.UseIgnoreCase {
		if c >= 'A' && c <= 'Z' {
			return int32(c | 0x20)
		}
	}
	if s.UseDBCcaseConverter {
		if c == 12288 {
			return int32(' ')
		}
		if c >= 65280 && c < 65375 {
			k := int32(c - 65248)
			if s.UseIgnoreCase {
				if 'A' <= k && k <= 'Z' {
					k = int32(k | 0x20)
				}
			}
			return k
		}
	}
	if s.UseSimplifiedChineseConverter {
		if c >= 0x4e00 && c <= 0x9fa5 {
			return node.GetSimplified(int(c - 0x4e00))
		}
	}
	return c
}
func (s *IllegalWordsSearch) toSenseWord2(text string) string {
	_text := []rune(text)
	stringBuilder := make([]int32, len(_text))
	for i, c := range _text {
		stringBuilder[i] = s.toSenseWord(c)
	}
	return string(stringBuilder)
}

func (s *IllegalWordsSearch) SetSkipWords(skipList string) {
	s._skipBitArray = make([]bool, 0x10000)
	for _, c := range skipList {
		s._skipBitArray[c] = true
	}
}

func (s *IllegalWordsSearch) SetBlacklist(blacklist []int) error {
	if len(s.I_keywords) != len(blacklist) {
		return errors.New("请关键字与黑名单列表的长度要一样长！")
	}
	s._blacklist = blacklist
	return nil
}

func (s *IllegalWordsSearch) SetKeywords(keywords []string) {
	list := make([]string, 0)
	for _, val := range keywords {
		list = append(list, s.toSenseWord2(val))
	}
	s.BaseSearchEx.SetKeywords(list)
}

func (s *IllegalWordsSearch) Save2(f *os.File) {
	s.BaseSearchEx.Save2(f)

	f.Write(s.intToBytes(s.boolToInt(s.UseSkipWordFilter)))
	f.Write(s.intToBytes(len(s._skipBitArray)))
	for _, key := range s._skipBitArray {
		f.Write(s.intToBytes(s.boolToInt(key)))
	}
	f.Write(s.intToBytes(s.boolToInt(s.UseDuplicateWordFilter)))
	f.Write(s.intToBytes(s.boolToInt(s.UseBlacklistFilter)))

	f.Write(s.intToBytes(len(s._blacklist)))
	for _, key := range s._blacklist {
		f.Write(s.intToBytes(key))
	}
	f.Write(s.intToBytes(s.boolToInt(s.UseDBCcaseConverter)))
	f.Write(s.intToBytes(s.boolToInt(s.UseSimplifiedChineseConverter)))
	f.Write(s.intToBytes(s.boolToInt(s.UseIgnoreCase)))
}

func (s *IllegalWordsSearch) Load2(f *os.File) {
	s.BaseSearchEx.Load2(f)
	intBs := make([]byte, 4)

	f.Read(intBs)
	s.UseSkipWordFilter = s.intToBool(s.bytesToInt(intBs))

	f.Read(intBs)
	length := s.bytesToInt(intBs)
	s._skipBitArray = make([]bool, length)
	for i := 0; i < length; i++ {
		f.Read(intBs)
		s._skipBitArray[i] = s.intToBool(s.bytesToInt(intBs))
	}

	f.Read(intBs)
	s.UseDuplicateWordFilter = s.intToBool(s.bytesToInt(intBs))
	f.Read(intBs)
	s.UseBlacklistFilter = s.intToBool(s.bytesToInt(intBs))

	f.Read(intBs)
	length = s.bytesToInt(intBs)
	s._blacklist = make([]int, length)
	for i := 0; i < length; i++ {
		f.Read(intBs)
		s._blacklist[i] = s.bytesToInt(intBs)
	}
	f.Read(intBs)
	s.UseDBCcaseConverter = s.intToBool(s.bytesToInt(intBs))

	f.Read(intBs)
	s.UseSimplifiedChineseConverter = s.intToBool(s.bytesToInt(intBs))

	f.Read(intBs)
	s.UseIgnoreCase = s.intToBool(s.bytesToInt(intBs))
}

func (s *IllegalWordsSearch) intToBytes(i int) []byte {
	x := int32(i)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
func (s *IllegalWordsSearch) bytesToInt(bs []byte) int {
	bytesBuffer := bytes.NewBuffer(bs)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
func (s *IllegalWordsSearch) boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
func (s *IllegalWordsSearch) intToBool(i int) bool {
	return i != 0
}
