package cache

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"github.com/Masterminds/semver/v3"
	"github.com/redis/rueidis"
	"JH-Forum/internal/core"
)

var (
	_ core.RedisCache = (*redisCache)(nil)            // 确保实现了 Redis 缓存接口
	_ tweetsCache     = (*redisCacheTweetsCache)(nil) // 确保实现了推文缓存接口
)

const (
	_cacheIndexKeyPattern = _cacheIndexKey + "*"
	_pushToSearchJobKey   = "paopao_push_to_search_job"
	_countLoginErrKey     = "paopao_count_login_err"
	_imgCaptchaKey        = "paopao_img_captcha:"

	_countWhisperKey      = "paopao_whisper_key"
	_rechargeStatusKey    = "paopao_recharge_status:"
)

// redisCache 实现了 Redis 缓存接口
type redisCache struct {
	c rueidis.Client
}

// redisCacheTweetsCache 实现了与推文相关的 Redis 缓存接口
type redisCacheTweetsCache struct {
	expireDuration time.Duration
	expireInSecond int64
	c              rueidis.Client
}

// getTweetsBytes 从 Redis 中获取推文数据的字节表示
func (s *redisCacheTweetsCache) getTweetsBytes(key string) ([]byte, error) {
	res, err := rueidis.MGetCache(s.c, context.Background(), s.expireDuration, []string{key})
	if err != nil {
		return nil, err
	}
	message := res[key]
	return message.AsBytes()
}

// setTweetsBytes 将推文数据的字节表示存入 Redis
func (s *redisCacheTweetsCache) setTweetsBytes(key string, bs []byte) error {
	cmd := s.c.B().Set().Key(key).Value(rueidis.BinaryString(bs)).ExSeconds(s.expireInSecond).Build()
	return s.c.Do(context.Background(), cmd).Error()
}

// delTweets 删除 Redis 中的推文数据
func (s *redisCacheTweetsCache) delTweets(keys []string) error {
	cmd := s.c.B().Del().Key(keys...).Build()
	return s.c.Do(context.Background(), cmd).Error()
}

// allKeys 返回所有与推文相关的缓存键
func (s *redisCacheTweetsCache) allKeys() (res []string, err error) {
	ctx, cursor := context.Background(), uint64(0)
	for {
		cmd := s.c.B().Scan().Cursor(cursor).Match(_cacheIndexKeyPattern).Count(50).Build()
		entry, err := s.c.Do(ctx, cmd).AsScanEntry()
		if err != nil {
			return nil, err
		}
		res = append(res, entry.Elements...)
		if entry.Cursor != 0 {
			cursor = entry.Cursor
			continue
		}
		break
	}
	return
}

// Name 返回缓存实例的名称
func (s *redisCacheTweetsCache) Name() string {
	return "RedisCacheIndex"
}

// Version 返回缓存实例的版本信息
func (s *redisCacheTweetsCache) Version() *semver.Version {
	return semver.MustParse("v0.1.0")
}

// SetPushToSearchJob 将推送到搜索作业设置到 Redis
func (r *redisCache) SetPushToSearchJob(ctx context.Context) error {
	return r.c.Do(ctx, r.c.B().Set().
		Key(_pushToSearchJobKey).Value("1").
		Nx().ExSeconds(3600).
		Build()).Error()
}

// DelPushToSearchJob 从 Redis 删除推送到搜索作业
func (r *redisCache) DelPushToSearchJob(ctx context.Context) error {
	return r.c.Do(ctx, r.c.B().Del().Key(_pushToSearchJobKey).Build()).Error()
}

// SetImgCaptcha 将图片验证码设置到 Redis
func (r *redisCache) SetImgCaptcha(ctx context.Context, id string, value string) error {
	return r.c.Do(ctx, r.c.B().Set().
		Key(_imgCaptchaKey+id).Value(value).
		ExSeconds(300).
		Build()).Error()
}

// GetImgCaptcha 从 Redis 获取图片验证码
func (r *redisCache) GetImgCaptcha(ctx context.Context, id string) (string, error) {
	res, err := r.c.Do(ctx, r.c.B().Get().Key(_imgCaptchaKey+id).Build()).AsBytes()
	return unsafe.String(&res[0], len(res)), err
}

// DelImgCaptcha 从 Redis 删除图片验证码
func (r *redisCache) DelImgCaptcha(ctx context.Context, id string) error {
	return r.c.Do(ctx, r.c.B().Del().Key(_imgCaptchaKey+id).Build()).Error()
}


// GetCountLoginErr 从 Redis 获取登录错误计数
func (r *redisCache) GetCountLoginErr(ctx context.Context, id int64) (int64, error) {
	return r.c.Do(ctx, r.c.B().Get().Key(fmt.Sprintf("%s:%d", _countLoginErrKey, id)).Build()).AsInt64()
}

// DelCountLoginErr 从 Redis 删除登录错误计数
func (r *redisCache) DelCountLoginErr(ctx context.Context, id int64) error {
	return r.c.Do(ctx, r.c.B().Del().Key(fmt.Sprintf("%s:%d", _countLoginErrKey, id)).Build()).Error()
}

// IncrCountLoginErr 增加登录错误计数，并设置过期时间
func (r *redisCache) IncrCountLoginErr(ctx context.Context, id int64) error {
	err := r.c.Do(ctx, r.c.B().Incr().Key(fmt.Sprintf("%s:%d", _countLoginErrKey, id)).Build()).Error()
	if err == nil {
		err = r.c.Do(ctx, r.c.B().Expire().Key(fmt.Sprintf("%s:%d", _countLoginErrKey, id)).Seconds(3600).Build()).Error()
	}
	return err
}

// GetCountWhisper 从 Redis 获取私信计数
func (r *redisCache) GetCountWhisper(ctx context.Context, uid int64) (int64, error) {
	return r.c.Do(ctx, r.c.B().Get().Key(fmt.Sprintf("%s:%d", _countWhisperKey, uid)).Build()).AsInt64()
}

// IncrCountWhisper 增加私信计数，并设置过期时间
func (r *redisCache) IncrCountWhisper(ctx context.Context, uid int64) (err error) {
	key := fmt.Sprintf("%s:%d", _countWhisperKey, uid)
	if err = r.c.Do(ctx, r.c.B().Incr().Key(key).Build()).Error(); err == nil {
		currentTime := time.Now()
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())
		err = r.c.Do(ctx, r.c.B().Expire().Key(key).Seconds(int64(endTime.Sub(currentTime)/time.Second)).Build()).Error()
	}
	return
}

// SetRechargeStatus 将充值状态设置到 Redis
func (r *redisCache) SetRechargeStatus(ctx context.Context, tradeNo string) error {
	return r.c.Do(ctx, r.c.B().Set().
		Key(_rechargeStatusKey+tradeNo).Value("1").
		Nx().ExSeconds(5).Build()).Error()
}

// DelRechargeStatus 从 Redis 删除充值状态
func (r *redisCache) DelRechargeStatus(ctx context.Context, tradeNo string) error {
	return r.c.Do(ctx, r.c.B().Del().Key(_rechargeStatusKey+tradeNo).Build()).Error()
}
