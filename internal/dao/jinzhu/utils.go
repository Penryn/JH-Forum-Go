package jinzhu

import (
	"JH-Forum/internal/core/cs"
	"JH-Forum/internal/dao/jinzhu/dbr"
	"gorm.io/gorm"
)

// createTags 根据提供的标签名称列表创建标签。
func createTags(db *gorm.DB, userId int64, tags []string) (res cs.TagInfoList, err error) {
	for _, name := range tags {
		tag := &dbr.Tag{Tag: name}
		if tag, err = tag.Get(db); err == nil {
			// 更新已存在的标签
			tag.QuoteNum++
			if err = tag.Update(db); err != nil {
				return
			}
		} else {
			// 创建新标签
			if tag, err = (&dbr.Tag{
				UserID:   userId,
				QuoteNum: 1,
				Tag:      name,
			}).Create(db); err != nil {
				return
			}
		}
		// 将结果添加到返回列表中
		res = append(res, &cs.TagInfo{
			ID:       tag.ID,
			UserID:   tag.UserID,
			Tag:      tag.Tag,
			QuoteNum: tag.QuoteNum,
		})
	}
	return
}

// decrTagsByIds 根据提供的标签ID列表递减标签引用计数。
func decrTagsByIds(db *gorm.DB, ids []int64) (err error) {
	for _, id := range ids {
		tag := &dbr.Tag{Model: &dbr.Model{ID: id}}
		if tag, err = tag.Get(db); err == nil {
			tag.QuoteNum--
			if err = tag.Update(db); err != nil {
				return
			}
		} else {
			continue
		}
	}
	return nil
}

// deleteTags 根据提供的标签名称列表删除标签及递减其引用计数。
func deleteTags(db *gorm.DB, tags []string) error {
	allTags, err := (&dbr.Tag{}).TagsFrom(db, tags)
	if err != nil {
		return err
	}
	for _, tag := range allTags {
		tag.QuoteNum--
		if tag.QuoteNum < 0 {
			tag.QuoteNum = 0
		}
		// 宽松处理错误，尽可能更新标签记录，只记录最后一次错误
		if e := tag.Update(db); e != nil {
			err = e
		}
	}
	return err
}

// getUsersByIDs 根据提供的用户ID列表获取用户列表。
func getUsersByIDs(db *gorm.DB, ids []int64) ([]*dbr.User, error) {
	user := &dbr.User{}
	return user.List(db, &dbr.ConditionsT{
		"id IN ?": ids,
	}, 0, 0)
}
