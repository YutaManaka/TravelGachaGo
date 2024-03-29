package main

import (
	"fmt"
  	"time"
	"math/rand"
	"os"

  	"github.com/gin-gonic/gin"
  	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
  	_ "github.com/go-sql-driver/mysql"
  )

  // DB接続
  func sqlConnect() (database *gorm.DB) {
	// env読み込み
	godotenv.Load(".env")
	DBMS := os.Getenv("DBMS")
	USER := os.Getenv("USER")
	PASS := os.Getenv("PASS")
	PROTOCOL := os.Getenv("PROTOCOL")
	DBNAME := os.Getenv("DBNAME")

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	count := 0
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
	  for {
		if err == nil {
		  fmt.Println("エラーなし")
		  break
		}
		fmt.Print(".")
		time.Sleep(time.Second)
		count++
		if count > 180 {
		  fmt.Println("")
		  fmt.Println("DB接続失敗")
		  panic(err)
		}
		db, err = gorm.Open(DBMS, CONNECT)
	  }
	}
	fmt.Println("DB接続成功")

	return db
  }

  // モデル定義
  type Gacha struct {
	gorm.Model
	Name string
	Used int
  }

  // ルーティング
  func main() {
	// マイグレーション
	db := sqlConnect()
	db.AutoMigrate(&Gacha{})
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/templates", "./templates")

	// form表示
	router.GET("/", func(ctx *gin.Context) {
	  ctx.HTML(200, "form.html", gin.H{})
	})

	// レコード作成
	router.POST("/create", func(ctx *gin.Context) {
		db := sqlConnect()
		defer db.Close()

		name := ctx.PostForm("name")
		db.Create(&Gacha{Name: name, Used: 0})
		defer db.Close()

		ctx.Redirect(302, "/")
	})

	// 目的地を取得
	router.GET("/destination", func(ctx *gin.Context) {
		db := sqlConnect()
		defer db.Close()

		// 未使用のレコードを取得
		destinations := []Gacha{}
		db.Where("used = ?", 0).Find(&destinations)

		// 未使用レコードがない場合、エラーページへ遷移
		if len(destinations) == 0 {
			ctx.HTML(200, "error.html", gin.H{})
		} else {
			// ランダムに選択
			index := rand.Intn(len(destinations))

			// 利用済に更新
			db.Model(&Gacha{}).Where("id = ?", destinations[index].ID).Update("used", 1)

			// 都市名を表示
			ctx.HTML(200, "destination.html", gin.H{"destination": destinations[index].Name})
		}
	})

	// 再度ガチャ
	router.GET("/retry", func(ctx *gin.Context) {
		ctx.HTML(200, "form.html", gin.H{})
	})

	router.Run()
  }
