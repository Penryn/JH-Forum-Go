// 该代码定义了一个推文列表盒子结构，包含推文列表和总数信息。

package cs

// TweetBox 表示推文列表盒子，包含推文列表和总数信息
type TweetBox struct {
	Tweets TweetList // 推文列表
	Total  int64     // 总数
}
