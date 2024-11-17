package migrations

import (
	"gorm.io/gorm"
	"log"
	"vidspark/apps/comment"
	"vidspark/apps/user"
	"vidspark/apps/vedio"
)

// Migrate 数据库迁移，不迁移base的model，后续有别的model迁移需要可以往代码加
func Migrate(db *gorm.DB) {
	log.Println("Migrating database ...")
	err := db.AutoMigrate(&user.User{}, &vedio.Video{}, &comment.Comment{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration complete")
}
