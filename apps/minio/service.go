package minio

import (
	"fmt"
	"github.com/minio/minio-go"
	"io"
	"log"
	"vidspark/configs"
	"vidspark/global"
)

var minioClient *minio.Client

func InitMinioClient() (err error) {
	config := configs.InitConfig()
	minioUrl := config.Minio.Host
	minioPort := config.Minio.Port
	minioAccessKey := config.Minio.AccessKey
	minioSecretKey := config.Minio.SecretKey
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
func UploadVideo(userID int, file io.Reader, fileName string) (string, error) {
	bucketName := VEDIO_BUCKET

	config := configs.InitConfig()

	endpoint := fmt.Sprintf("http://%s:%s", config.Minio.Host, config.Minio.Port)
	objectName := fmt.Sprintf("%v/videos/%s", userID, fileName)

	_, err := minioClient.PutObject(bucketName, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("视频上传失败: %w", err)
	}

	// 返回文件的访问 URL
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectName), nil
}
