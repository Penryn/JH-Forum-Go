// 该代码定义了对象存储服务的接口，用于实现基于阿里云OSS、MINIO或其他存储服务的基本功能。

package core

import (
	"io"
)

// ObjectStorageService 对象存储服务接口，继承了创建和删除服务接口
type ObjectStorageService interface {
	OssCreateService
	OssDeleteService

	SignURL(objectKey string, expiredInSec int64) (string, error) // 生成签名URL
	ObjectURL(objectKey string) string                            // 获取对象URL
	ObjectKey(cUrl string) string                                 // 获取对象Key
}

// OssCreateService 对象存储系统的对象创建服务接口
type OssCreateService interface {
	PutObject(objectKey string, reader io.Reader, objectSize int64, contentType string, persistance bool) (string, error) // 上传对象
	PersistObject(objectKey string) error                                                                                 // 持久化对象
}

// OssDeleteService 对象存储系统的对象删除服务接口
type OssDeleteService interface {
	DeleteObject(objectKey string) error          // 删除单个对象
	DeleteObjects(objectKeys []string) error      // 批量删除对象
	IsObjectExist(objectKey string) (bool, error) // 判断对象是否存在
}
