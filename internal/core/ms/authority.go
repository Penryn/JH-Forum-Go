// 该代码定义了一些常量、类型和方法，用于描述用户行为和好友过滤。

package ms

import (
	"JH-Forum/pkg/types"
)

// act 表示用户动作的枚举类型
const (
	ActRegisterUser act = iota
	ActCreatePublicTweet
	ActCreatePublicAttachment
	ActCreatePublicPicture
	ActCreatePublicVideo
	ActCreatePrivateTweet
	ActCreatePrivateAttachment
	ActCreatePrivatePicture
	ActCreatePrivateVideo
	ActCreateFriendTweet
	ActCreateFriendAttachment
	ActCreateFriendPicture
	ActCreateFriendVideo
	ActCreatePublicComment
	ActCreatePublicPicureComment
	ActCreateFriendComment
	ActCreateFriendPicureComment
	ActCreatePrivateComment
	ActCreatePrivatePicureComment
	ActStickTweet
	ActTopTweet
	ActLockTweet
	ActVisibleTweet
	ActDeleteTweet
	ActCreateActivationCode
)

type (
	act uint8

	// FriendFilter 好友过滤器，用于描述用户的好友关系
	FriendFilter map[int64]types.Empty
	// FriendSet 好友集合，用于描述用户的好友关系
	FriendSet map[string]types.Empty

	// Action 表示用户执行的动作
	Action struct {
		Act    act   // 用户动作类型
		UserId int64 // 用户ID
	}
)

// IsFriend 判断用户是否是好友
func (f FriendFilter) IsFriend(userId int64) bool {
	_, ok := f[userId]
	return ok
}

// IsAllow 判断用户是否允许执行某个动作
func (a act) IsAllow(user *User, userId int64, isFriend bool, isActivation bool) bool {
	if user.IsAdmin {
		return true
	}

	if user.ID == userId && isActivation {
		switch a {
		case ActCreatePublicTweet, ActCreatePublicAttachment, ActCreatePublicPicture, ActCreatePublicVideo,
			ActCreatePrivateTweet, ActCreatePrivateAttachment, ActCreatePrivatePicture, ActCreatePrivateVideo,
			ActCreateFriendTweet, ActCreateFriendAttachment, ActCreateFriendPicture, ActCreateFriendVideo,
			ActCreatePrivateComment, ActCreatePrivatePicureComment, ActStickTweet, ActLockTweet, ActVisibleTweet,
			ActDeleteTweet:
			return true
		}
	}

	if user.ID == userId && !isActivation {
		switch a {
		case ActCreatePrivateTweet, ActCreatePrivateComment, ActStickTweet, ActLockTweet, ActDeleteTweet:
			return true
		}
	}

	if isFriend && isActivation {
		switch a {
		case ActCreatePublicComment, ActCreatePublicPicureComment, ActCreateFriendComment, ActCreateFriendPicureComment:
			return true
		}
	}

	if !isFriend && isActivation {
		switch a {
		case ActCreatePublicComment, ActCreatePublicPicureComment:
			return true
		}
	}

	return false
}
