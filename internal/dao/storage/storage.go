// 该包实现了与多个云对象存储服务交互的功能，包括Aliyun OSS、Tencent COS、Huawei OBS、MinIO和S3。

package storage

import (
	"time"

	"github.com/alimy/tryst/cfg"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"JH-Forum/internal/conf"
	"JH-Forum/internal/core"
	"github.com/sirupsen/logrus"
)

// MustMinioService 创建并返回 MinIO 服务客户端
func MustMinioService() (core.ObjectStorageService, core.VersionInfo) {
	// 初始化 minio 客户端对象。
	client, err := minio.New(conf.MinIOSetting.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.MinIOSetting.AccessKey, conf.MinIOSetting.SecretKey, ""),
		Secure: conf.MinIOSetting.Secure,
	})
	if err != nil {
		logrus.Fatalf("storage.MustMinioService 创建客户端失败: %s", err)
	}

	domain := conf.GetOssDomain()
	var cs core.OssCreateService
	if cfg.If("OSS:TempDir") {
		cs = &minioCreateTempDirServant{
			client:  client,
			bucket:  conf.MinIOSetting.Bucket,
			domain:  domain,
			tempDir: conf.ObjectStorage.TempDirSlash(),
		}
		logrus.Debugln("使用 OSS:TempDir 功能")
	} else if cfg.If("OSS:Retention") {
		cs = &minioCreateRetentionServant{
			client:          client,
			bucket:          conf.MinIOSetting.Bucket,
			domain:          domain,
			retainInDays:    time.Duration(conf.ObjectStorage.RetainInDays) * time.Hour * 24,
			retainUntilDate: time.Date(2049, time.December, 1, 12, 0, 0, 0, time.UTC),
		}
		logrus.Debugln("使用 OSS:Retention 功能")
	} else {
		cs = &minioCreateServant{
			client: client,
			bucket: conf.MinIOSetting.Bucket,
			domain: domain,
		}
		logrus.Debugln("使用 OSS:Direct 功能")
	}

	obj := &minioServant{
		OssCreateService: cs,
		client:           client,
		bucket:           conf.MinIOSetting.Bucket,
		domain:           domain,
	}
	return obj, obj
}
