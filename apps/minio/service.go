package minio

import (
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"kowhai/global"
	"log"
)

var minioClient *minio.Client

func InitMinioClient() (err error) {
	minioUrl := global.Config.Minio.Host
	minioPort := global.Config.Minio.Port
	minioAccessKey := global.Config.Minio.AccessKey
	minioSecretKey := global.Config.Minio.SecretKey
	minioClient, err = minio.New(fmt.Sprintf("%s:%s", minioUrl, minioPort), minioAccessKey, minioSecretKey, false)
	if err != nil {
		global.Logger.Error("minio服务连接失败", err.Error())
		return
	}
	log.Println("minio服务初始化完成!")
	return nil
}

func InitStorageBuckets() error {
	// 初始化存储桶
	buckets := []string{VEDIO_BUCKET, AVATAR_BUCKET}
	for _, bucketName := range buckets {
		exists, err := minioClient.BucketExists(bucketName)
		if err != nil {
			return fmt.Errorf("检查存储桶是否存在失败: %w", err)
		}
		if !exists {
			err = minioClient.MakeBucket(bucketName, LOCATION)
			if err != nil {
				return fmt.Errorf("创建存储桶 %s 失败: %w", bucketName, err)
			}
			log.Printf("存储桶 %s 创建成功！\n", bucketName)
		} else {
			log.Printf("存储桶 %s 已存在\n", bucketName)
		}
	}
	return nil
}

// 上传视频
func UploadVideo(userId int, fileDir, fileName string, data io.Reader) error {
	bucketName := VEDIO_BUCKET
	objectName := fmt.Sprintf("%d/%s/%s", userId, fileDir, fileName)
	_, err := minioClient.PutObject(bucketName, objectName, data, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("视频上传失败: %w", err)
	}
	return nil
}

// 上传avatar
func Uploadavatar(userId int, filename string, data io.Reader) error {
	bucketName := AVATAR_BUCKET
	objectName := fmt.Sprintf("%d/%s", userId, filename)
	_, err := minioClient.PutObject(bucketName, objectName, data, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("avatar上传失败: %w", err)
	}
	return nil
}

func GetAvatarUrl(userId int, filename string) string {
	return fmt.Sprintf("http://%s:%s/%s/%d/%s", global.Config.Minio.Host, global.Config.Minio.Port, AVATAR_BUCKET, userId, filename)
}
