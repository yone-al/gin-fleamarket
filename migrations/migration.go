// アプリケーションの起動とDBのマイグレーション処理を分離
package main

import (
	"gin-fleamarket/infra"
	"gin-fleamarket/models"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.Item{}, &models.User{}); err != nil {
		panic("Failed to migrate database")
	}
}
