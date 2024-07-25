package web

import (
	"github.com/alimy/mir/v4"
	"JH-Forum/pkg/xerror"
)

// fileCheck 根据上传类型和文件大小进行检查。
func fileCheck(uploadType string, size int64) mir.Error {
	if uploadType != "public/video" &&
		uploadType != "public/image" &&
		uploadType != "public/avatar" &&
		uploadType != "attachment" {
		return xerror.InvalidParams
	}
	if size > 1024*1024*100 {
		return ErrFileInvalidSize.WithDetails("最大允许100MB")
	}
	return nil
}

// getFileExt 根据文件类型获取文件扩展名。
func getFileExt(s string) (string, mir.Error) {
	switch s {
	case "image/png":
		return ".png", nil
	case "image/jpg":
		return ".jpg", nil
	case "image/jpeg":
		return ".jpeg", nil
	case "image/gif":
		return ".gif", nil
	case "video/mp4":
		return ".mp4", nil
	case "video/quicktime":
		return ".mov", nil
	case "application/zip",
		"application/x-zip",
		"application/octet-stream",
		"application/x-zip-compressed":
		return ".zip", nil
	default:
		return "", ErrFileInvalidExt.WithDetails("仅允许 png/jpg/gif/mp4/mov/zip 类型")
	}
}
