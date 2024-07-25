// 该代码定义了一个结构体类型 IndexTweetList。

package ms

type IndexTweetList struct {
	Tweets []*PostFormated // 推文列表
	Total  int64           // 总数
}
