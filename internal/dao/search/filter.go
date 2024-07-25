// 包 search 实现了与 Tweet 相关的搜索过滤功能。

package search

import (
	"JH-Forum/internal/core"
	"JH-Forum/internal/core/ms"
	"JH-Forum/pkg/types"
)

// tweetSearchFilter 定义了 Tweet 搜索结果的过滤器。
type tweetSearchFilter struct {
	ams core.AuthorizationManageService // 授权管理服务
}

// filterResp 根据用户权限过滤搜索结果。
func (s *tweetSearchFilter) filterResp(user *ms.User, resp *core.QueryResp) {
	// 管理员不过滤
	if user != nil && user.IsAdmin {
		return
	}

	var item *ms.PostFormated
	items := resp.Items
	latestIndex := len(items) - 1
	if user == nil {
		// 未认证用户过滤非公开内容
		for i := 0; i <= latestIndex; i++ {
			item = items[i]
			if item.Visibility != core.PostVisitPublic {
				items[i] = items[latestIndex]
				items = items[:latestIndex]
				resp.Total--
				latestIndex--
				i--
			}
		}
	} else {
		// 认证用户过滤好友和私密内容
		var cutFriend, cutPrivate bool
		friendFilter := s.ams.BeFriendFilter(user.ID)
		friendFilter[user.ID] = types.Empty{}
		for i := 0; i <= latestIndex; i++ {
			item = items[i]
			cutFriend = (item.Visibility == core.PostVisitFriend && !friendFilter.IsFriend(item.UserID))
			cutPrivate = (item.Visibility == core.PostVisitPrivate && user.ID != item.UserID)
			if cutFriend || cutPrivate {
				items[i] = items[latestIndex]
				items = items[:latestIndex]
				resp.Total--
				latestIndex--
				i--
			}
		}
	}

	resp.Items = items
}
