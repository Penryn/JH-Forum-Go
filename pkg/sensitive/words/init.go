package words

import (
	"fmt"
	"os"
	"strings"
)

var Search *IllegalWordsSearch

func init() {
	fmt.Println("init")
	bs, err := os.ReadFile("./pkg/sensitive/data/badWord.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	s := string(bs)
	s = strings.Replace(s, "\r\n", "\n", -1)
	s = strings.Replace(s, "\r", "\n", -1)
	sp := strings.Split(s, "\n") // 如果文件中使用的是 '\n' 分隔符
	list := make([]string, 0)
	for _, item := range sp {
		// 去掉每行最后一个逗号
		if strings.HasSuffix(item, ",") {
			item = item[:len(item)-1]
		item = strings.TrimSuffix(item, "，")
			item = item[:len(item)-3]
		}
		list = append(list, item)
	}
	Search = NewIllegalWordsSearch()
	Search.SetKeywords(list)
}
