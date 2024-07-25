// 该代码定义了多个接口，用于实现推文的检索、管理和辅助服务。

package core

import (
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/core/ms"
)

// TweetService 推文检索服务接口
type TweetService interface {
	GetPostByID(id int64) (*ms.Post, error)                                                                   // 通过ID获取推文
	GetPosts(conditions ms.ConditionsT, offset, limit int) ([]*ms.Post, error)                                // 获取推文列表
	GetPostCount(conditions ms.ConditionsT) (int64, error)                                                    // 获取推文数量
	GetUserPostStar(postID, userID int64) (*ms.PostStar, error)                                               // 获取用户推文点赞信息
	GetUserPostStars(userID int64, limit int, offset int) ([]*ms.PostStar, error)                             // 获取用户推文点赞列表
	GetUserPostStarCount(userID int64) (int64, error)                                                         // 获取用户推文点赞数量
	GetUserPostCollection(postID, userID int64) (*ms.PostCollection, error)                                   // 获取用户推文收藏信息
	GetUserPostCollections(userID int64, offset, limit int) ([]*ms.PostCollection, error)                     // 获取用户推文收藏列表
	GetUserPostCollectionCount(userID int64) (int64, error)                                                   // 获取用户推文收藏数量
	GetPostAttatchmentBill(postID, userID int64) (*ms.PostAttachmentBill, error)                              // 获取推文附件账单
	GetPostContentsByIDs(ids []int64) ([]*ms.PostContent, error)                                              // 通过ID获取推文内容列表
	GetPostContentByID(id int64) (*ms.PostContent, error)                                                     // 通过ID获取推文内容
	ListUserStarTweets(user *cs.VistUser, limit int, offset int) ([]*ms.PostStar, int64, error)               // 列出用户点赞的推文
	ListUserMediaTweets(user *cs.VistUser, limit int, offset int) ([]*ms.Post, int64, error)                  // 列出用户的媒体推文
	ListUserCommentTweets(user *cs.VistUser, limit int, offset int) ([]*ms.Post, int64, error)                // 列出用户评论的推文
	ListUserTweets(userId int64, style uint8, justEssence bool, limit, offset int) ([]*ms.Post, int64, error) // 列出用户推文
	ListFollowingTweets(userId int64, limit, offset int) ([]*ms.Post, int64, error)                           // 列出关注者的推文
	ListIndexNewestTweets(limit, offset int) ([]*ms.Post, int64, error)                                       // 列出最新推文
	ListIndexHotsTweets(limit, offset int) ([]*ms.Post, int64, error)                                         // 列出热门推文
	ListSyncSearchTweets(limit, offset int) ([]*ms.Post, int64, error)                                        // 列出同步搜索推文
}

// TweetManageService 推文管理服务接口，包括创建/删除/更新推文
type TweetManageService interface {
	CreatePost(post *ms.Post) (*ms.Post, error)                            // 创建推文
	DeletePost(post *ms.Post) ([]string, error)                            // 删除推文
	LockPost(post *ms.Post) error                                          // 锁定推文
	StickPost(post *ms.Post) error                                         // 置顶推文
	HighlightPost(userId, postId int64) (int, error)                       // 高亮推文
	VisiblePost(post *ms.Post, visibility cs.TweetVisibleType) error       // 设置推文可见性
	UpdatePost(post *ms.Post) error                                        // 更新推文
	CreatePostStar(postID, userID int64) (*ms.PostStar, error)             // 创建推文点赞
	DeletePostStar(p *ms.PostStar) error                                   // 删除推文点赞
	CreatePostCollection(postID, userID int64) (*ms.PostCollection, error) // 创建推文收藏
	DeletePostCollection(p *ms.PostCollection) error                       // 删除推文收藏
	CreatePostContent(content *ms.PostContent) (*ms.PostContent, error)    // 创建推文内容
	CreateAttachment(obj *ms.Attachment) (int64, error)                    // 创建附件
}

// TweetHelpService 推文辅助服务接口
type TweetHelpService interface {
	RevampPosts(posts []*ms.PostFormated) ([]*ms.PostFormated, error) // 修订推文
	MergePosts(posts []*ms.Post) ([]*ms.PostFormated, error)          // 合并推文
}

// TweetServantA 推文检索服务(版本A)接口
type TweetServantA interface {
	TweetInfoById(id int64) (*cs.TweetInfo, error)                               // 通过ID获取推文信息
	TweetItemById(id int64) (*cs.TweetItem, error)                               // 通过ID获取推文项
	UserTweets(visitorId, userId int64) (cs.TweetList, error)                    // 获取用户推文
	ReactionByTweetId(userId int64, tweetId int64) (*cs.ReactionItem, error)     // 获取推文的反应
	UserReactions(userId int64, limit int, offset int) (cs.ReactionList, error)  // 获取用户的反应列表
	FavoriteByTweetId(userId int64, tweetId int64) (*cs.FavoriteItem, error)     // 获取推文的收藏
	UserFavorites(userId int64, limit int, offset int) (cs.FavoriteList, error)  // 获取用户的收藏列表
	AttachmentByTweetId(userId int64, tweetId int64) (*cs.AttachmentBill, error) // 获取推文的附件账单
}

// TweetManageServantA 推文管理服务(版本A)接口，包括创建/删除/更新推文
type TweetManageServantA interface {
	CreateAttachment(obj *cs.Attachment) (int64, error)                   // 创建附件
	CreateTweet(userId int64, req *cs.NewTweetReq) (*cs.TweetItem, error) // 创建推文
	DeleteTweet(userId int64, tweetId int64) ([]string, error)            // 删除推文
	LockTweet(userId int64, tweetId int64) error                          // 锁定推文
	StickTweet(userId int64, tweetId int64) error                         // 置顶推文
	VisibleTweet(userId int64, visibility cs.TweetVisibleType) error      // 设置推文可见性
	CreateReaction(userId int64, tweetId int64) error                     // 创建反应
	DeleteReaction(userId int64, reactionId int64) error                  // 删除反应
	CreateFavorite(userId int64, tweetId int64) error                     // 创建收藏
	DeleteFavorite(userId int64, favoriteId int64) error                  // 删除收藏
}

// TweetHelpServantA 推文辅助服务(版本A)接口
type TweetHelpServantA interface {
	RevampTweets(tweets cs.TweetList) (cs.TweetList, error) // 修订推文
	MergeTweets(tweets cs.TweetInfo) (cs.TweetList, error)  // 合并推文
}
