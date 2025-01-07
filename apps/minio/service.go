package minio

import (
	"fmt"
	"github.com/minio/minio-go"
	"log"
	"vidspark/configs"
)

var minioClient *minio.Client

func InitMinioClient() (err error) {
	config := configs.InitConfig()
	minioUrl := fmt.Sprintf("http://%s", config.Minio.Host)
	minioPort := config.Minio.Port
	minioAccessKey := config.Minio.AccessKey
	minioSecretKey := config.Minio.SecretKey
	minioClient, err = minio.New(fmt.Sprintf("%s:%s", minioUrl, minioPort), minioAccessKey, minioSecretKey, false)
	if err != nil {
		log.Fatalln("minio服务连接失败", err.Error())
		return
	}
	log.Println("minio服务初始化完成!")
	return nil
}
