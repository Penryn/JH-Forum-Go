// 该代码定义了两个接口，用于实现话题服务及其版本A。

package core

import (
	"JH-Forum/internal/core/cs"
)

// TopicService 话题服务接口
type TopicService interface {
	UpsertTags(userId int64, tags []string) (cs.TagInfoList, error)        // 插入或更新标签
	DecrTagsById(ids []int64) error                                        // 根据ID减少标签
	ListTags(typ cs.TagType, limit int, offset int) (cs.TagList, error)    // 列出标签
	TagsByKeyword(keyword string) (cs.TagInfoList, error)                  // 根据关键词搜索标签
	GetHotTags(userId int64, limit int, offset int) (cs.TagList, error)    // 获取热门标签
	GetNewestTags(userId int64, limit int, offset int) (cs.TagList, error) // 获取最新标签
	GetFollowTags(userId int64, limit int, offset int) (cs.TagList, error) // 获取关注的标签
	FollowTopic(userId int64, topicId int64) error                         // 关注话题
	UnfollowTopic(userId int64, topicId int64) error                       // 取消关注话题
	StickTopic(userId int64, topicId int64) (int8, error)                  // 置顶话题
}

// TopicServantA 话题服务(版本A)接口
type TopicServantA interface {
	UpsertTags(userId int64, tags []string) (cs.TagInfoList, error)     // 插入或更新标签
	DecrTagsById(ids []int64) error                                     // 根据ID减少标签
	ListTags(typ cs.TagType, limit int, offset int) (cs.TagList, error) // 列出标签
	TagsByKeyword(keyword string) (cs.TagInfoList, error)               // 根据关键词搜索标签
}
