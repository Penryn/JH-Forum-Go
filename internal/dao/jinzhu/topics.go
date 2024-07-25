// package jinzhu 实现了话题服务，包括标签管理、话题关注、置顶等功能。

package jinzhu

import (
	"errors"
	"strings"

	"JH-Forum/internal/core"           // 引入核心服务包
	"JH-Forum/internal/core/cs"        // 引入核心服务-通用包
	"JH-Forum/internal/core/ms"        // 引入核心服务-消息服务包
	"JH-Forum/internal/dao/jinzhu/dbr" // 引入Jinzhu数据库访问包
	"gorm.io/gorm"                     // 引入gorm包
)

var (
	_ core.TopicService  = (*topicSrv)(nil)
	_ core.TopicServantA = (*topicSrvA)(nil)
)

type topicSrv struct {
	db             *gorm.DB
	tnTopicUser    string // 话题用户表名称
	tnDotTopicUser string // 话题用户表名称前缀
}

type topicSrvA struct {
	db *gorm.DB
}

type topicInfo struct {
	TopicId int64
	IsTop   int8
}

// newTopicService 创建新的话题服务实例
func newTopicService(db *gorm.DB) core.TopicService {
	return &topicSrv{
		db:             db,
		tnTopicUser:    db.NamingStrategy.TableName("TopicUser"),
		tnDotTopicUser: db.NamingStrategy.TableName("TopicUser") + ".",
	}
}

// newTopicServantA 创建新的话题服务实例（管理员）
func newTopicServantA(db *gorm.DB) core.TopicServantA {
	return &topicSrvA{
		db: db,
	}
}

// UpsertTags 新增或更新用户的标签
func (s *topicSrv) UpsertTags(userId int64, tags []string) (_ cs.TagInfoList, err error) {
	db := s.db.Begin()
	defer func() {
		if err == nil {
			db.Commit()
		} else {
			db.Rollback()
		}
	}()
	return createTags(db, userId, tags)
}

// DecrTagsById 根据标签ID减少标签数量
func (s *topicSrv) DecrTagsById(ids []int64) (err error) {
	db := s.db.Begin()
	defer func() {
		if err == nil {
			db.Commit()
		} else {
			db.Rollback()
		}
	}()
	return decrTagsByIds(db, ids)
}

// ListTags 按类型列出标签列表
func (s *topicSrv) ListTags(typ cs.TagType, offset, limit int) (res cs.TagList, err error) {
	conditions := &ms.ConditionsT{}
	switch typ {
	case cs.TagTypeHot:
		// 热门标签
		conditions = &ms.ConditionsT{
			"ORDER": "quote_num DESC",
		}
	case cs.TagTypeNew:
		// 最新标签
		conditions = &ms.ConditionsT{
			"ORDER": "id DESC",
		}
	}
	return s.listTags(conditions, limit, offset)
}

// GetHotTags 获取热门标签列表
func (s *topicSrv) GetHotTags(userId int64, limit int, offset int) (cs.TagList, error) {
	tags, err := s.listTags(&ms.ConditionsT{
		"ORDER": "quote_num DESC",
	}, limit, offset)
	if err != nil {
		return nil, err
	}
	return s.tagsFormatA(userId, tags)
}

// GetNewestTags 获取最新标签列表
func (s *topicSrv) GetNewestTags(userId int64, limit int, offset int) (cs.TagList, error) {
	tags, err := s.listTags(&ms.ConditionsT{
		"ORDER": "id DESC",
	}, limit, offset)
	if err != nil {
		return nil, err
	}
	return s.tagsFormatA(userId, tags)
}

// GetFollowTags 获取用户关注的标签列表
func (s *topicSrv) GetFollowTags(userId int64, limit int, offset int) (cs.TagList, error) {
	if userId < 0 {
		return nil, nil
	}
	userTopics := []*topicInfo{}
	err := s.db.Model(&dbr.TopicUser{}).
		Where("user_id=?", userId).
		Order("is_top DESC").
		Limit(limit).
		Offset(offset).
		Find(&userTopics).Error
	if err != nil {
		return nil, err
	}
	userTopicsMap := make(map[int64]*topicInfo, len(userTopics))
	topicIds := make([]int64, 0, len(userTopics))
	topicIdsMap := make(map[int64]int, len(userTopics))
	for idx, info := range userTopics {
		userTopicsMap[info.TopicId] = info
		topicIds = append(topicIds, info.TopicId)
		topicIdsMap[info.TopicId] = idx
	}
	var tags cs.TagInfoList
	err = s.db.Model(&dbr.Tag{}).Where("quote_num > 0 and id in ?", topicIds).Order("quote_num DESC").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	formattedTags, err := s.tagsFormatB(userTopicsMap, tags)
	if err != nil {
		return nil, err
	}
	// 置顶排序后处理
	res := make(cs.TagList, len(topicIds))
	for _, tag := range formattedTags {
		res[topicIdsMap[tag.ID]] = tag
	}
	return res, nil
}

// listTags 列出标签列表
func (s *topicSrv) listTags(conditions *ms.ConditionsT, limit int, offset int) (res cs.TagList, err error) {
	var tags []*dbr.Tag
	if tags, err = (&dbr.Tag{}).List(s.db, conditions, offset, limit); err == nil {
		if len(tags) == 0 {
			return
		}
		tagMap := make(map[int64][]*cs.TagItem, len(tags))
		for _, tag := range tags {
			item := &cs.TagItem{
				ID:       tag.ID,
				UserID:   tag.UserID,
				Tag:      tag.Tag,
				QuoteNum: tag.QuoteNum,
			}
			tagMap[item.UserID] = append(tagMap[item.UserID], item)
			res = append(res, item)
		}
		ids := make([]int64, len(tagMap))
		for userID := range tagMap {
			ids = append(ids, userID)
		}
		userInfos, err := (&dbr.User{}).ListUserInfoById(s.db, ids)
		if err != nil {
			return nil, err
		}
		for _, userInfo := range userInfos {
			for _, item := range tagMap[userInfo.ID] {
				item.User = userInfo
			}
		}
	}
	return
}

// tagsFormatA 格式化标签列表A
func (s *topicSrv) tagsFormatA(userID int64, tags cs.TagList) (cs.TagList, error) {
	tagIDs := make([]int64, len(tags))
	for idx, tag := range tags {
		tagIDs[idx] = tag.ID
	}
	if userID > -1 {
		userTopics := []*topicInfo{}
		err := s.db.Model(&dbr.TopicUser{}).Where("is_del=0 and user_id=? and topic_id in ?", userID, tagIDs).Find(&userTopics).Error
		if err != nil {
			return nil, err
		}
		userTopicsMap := make(map[int64]*topicInfo, len(userTopics))
		for _, info := range userTopics {
			userTopicsMap[info.TopicId] = info
		}
		for _, tag := range tags {
			if info, exist := userTopicsMap[tag.ID]; exist {
				tag.IsFollowing, tag.IsTop = 1, info.IsTop
			}
		}
	}
	return tags, nil
}

// tagsFormatB 格式化标签列表B
func (s *topicSrv) tagsFormatB(userTopicsMap map[int64]*topicInfo, tags cs.TagInfoList) (cs.TagList, error) {
	userIDs := make([]int64, len(tags))
	for idx, tag := range tags {
		userIDs[idx] = tag.UserID
	}
	users, err := (&dbr.User{}).ListUserInfoById(s.db, userIDs)
	if err != nil {
		return nil, err
	}
	tagList := cs.TagList{}
	for _, tag := range tags {
		tagFormatted := tag.Format()
		for _, user := range users {
			if user.ID == tagFormatted.UserID {
				tagFormatted.User = user
			}
		}
		tagList = append(tagList, tagFormatted)
	}
	if len(userTopicsMap) > 0 {
		for _, tag := range tagList {
			if info, exist := userTopicsMap[tag.ID]; exist {
				tag.IsFollowing, tag.IsTop = 1, info.IsTop
			}
		}
	}
	return tagList, nil
}

// TagsByKeyword 根据关键字获取标签列表
func (s *topicSrv) TagsByKeyword(keyword string) (res cs.TagInfoList, err error) {
	keyword = "%" + strings.Trim(keyword, " ") + "%"
	tag := &dbr.Tag{}
	var tags []*dbr.Tag
	if keyword == "%%" {
		tags, err = tag.List(s.db, &dbr.ConditionsT{
			"ORDER": "quote_num DESC",
		}, 0, 6)
	} else {
		tags, err = tag.List(s.db, &dbr.ConditionsT{
			"tag LIKE ?": keyword,
			"ORDER":      "quote_num DESC",
		}, 0, 6)
	}
	if err == nil {
		for _, tag := range tags {
			res = append(res, &cs.TagInfo{
				ID:       tag.ID,
				UserID:   tag.UserID,
				Tag:      tag.Tag,
				QuoteNum: tag.QuoteNum,
			})
		}
	}
	return
}

// UpsertTags 新增或更新标签（管理员）
func (s *topicSrvA) UpsertTags(userId int64, tags []string) (_ cs.TagInfoList, err error) {
	db := s.db.Begin()
	defer func() {
		if err == nil {
			db.Commit()
		} else {
			db.Rollback()
		}
	}()
	return createTags(db, userId, tags)
}

// DecrTagsById 根据ID减少标签（管理员）
func (s *topicSrvA) DecrTagsById(ids []int64) (err error) {
	db := s.db.Begin()
	defer func() {
		if err == nil {
			db.Commit()
		} else {
			db.Rollback()
		}
	}()
	return decrTagsByIds(db, ids)
}

// ListTags 按类型列出标签列表（管理员）
func (s *topicSrvA) ListTags(typ cs.TagType, offset, limit int) (res cs.TagList, err error) {
	conditions := &ms.ConditionsT{}
	switch typ {
	case cs.TagTypeHot:
		// 热门标签
		conditions = &ms.ConditionsT{
			"ORDER": "quote_num DESC",
		}
	case cs.TagTypeNew:
		// 最新标签
		conditions = &ms.ConditionsT{
			"ORDER": "id DESC",
		}
	}
	var tags []*dbr.Tag
	if tags, err = (&dbr.Tag{}).List(s.db, conditions, offset, limit); err == nil {
		if len(tags) == 0 {
			return
		}
		tagMap := make(map[int64][]*cs.TagItem, len(tags))
		for _, tag := range tags {
			item := &cs.TagItem{
				ID:       tag.ID,
				UserID:   tag.UserID,
				Tag:      tag.Tag,
				QuoteNum: tag.QuoteNum,
			}
			tagMap[item.UserID] = append(tagMap[item.UserID], item)
			res = append(res, item)
		}
		ids := make([]int64, len(tagMap))
		for userID := range tagMap {
			ids = append(ids, userID)
		}
		userInfos, err := (&dbr.User{}).ListUserInfoById(s.db, ids)
		if err != nil {
			return nil, err
		}
		for _, userInfo := range userInfos {
			for _, item := range tagMap[userInfo.ID] {
				item.User = userInfo
			}
		}
	}
	return
}

// TagsByKeyword 根据关键字获取标签列表（管理员）
func (s *topicSrvA) TagsByKeyword(keyword string) (res cs.TagInfoList, err error) {
	keyword = "%" + strings.Trim(keyword, " ") + "%"
	tag := &dbr.Tag{}
	var tags []*dbr.Tag
	if keyword == "%%" {
		tags, err = tag.List(s.db, &dbr.ConditionsT{
			"ORDER": "quote_num DESC",
		}, 0, 6)
	} else {
		tags, err = tag.List(s.db, &dbr.ConditionsT{
			"tag LIKE ?": keyword,
			"ORDER":      "quote_num DESC",
		}, 0, 6)
	}
	if err == nil {
		for _, tag := range tags {
			res = append(res, &cs.TagInfo{
				ID:       tag.ID,
				UserID:   tag.UserID,
				Tag:      tag.Tag,
				QuoteNum: tag.QuoteNum,
			})
		}
	}
	return
}

// FollowTopic 用户关注话题
func (s *topicSrv) FollowTopic(userId int64, topicId int64) (err error) {
	return s.db.Create(&dbr.TopicUser{
		UserID:  userId,
		TopicID: topicId,
		IsTop:   0,
	}).Error
}

// UnfollowTopic 用户取消关注话题
func (s *topicSrv) UnfollowTopic(userId int64, topicId int64) error {
	return s.db.Exec("DELETE FROM "+s.tnTopicUser+" WHERE user_id=? AND topic_id=?", userId, topicId).Error
}

// StickTopic 置顶话题
func (s *topicSrv) StickTopic(userId int64, topicId int64) (status int8, err error) {
	db := s.db.Begin()
	defer db.Rollback()

	m := &dbr.TopicUser{}
	err = db.Model(m).
		Where("user_id=? and topic_id=?", userId, topicId).
		UpdateColumn("is_top", gorm.Expr("1-is_top")).Error
	if err != nil {
		return
	}
	status = -1
	err = db.Model(m).Where("user_id=? and topic_id=?", userId, topicId).Select("is_top").Scan(&status).Error
	if err != nil {
		return
	}
	if status < 0 {
		return -1, errors.New("topic not exist")
	}

	db.Commit()
	return
}
