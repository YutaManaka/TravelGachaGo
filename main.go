package main

import (
	"fmt"
  	"time"

  	"github.com/gin-gonic/gin"
  	"github.com/jinzhu/gorm"
  	_ "github.com/go-sql-driver/mysql"
  )

  // TODO: env読み込み

  // DB接続
  func sqlConnect() (database *gorm.DB) {
	DBMS := "mysql"
	USER := "go_test"
	PASS := "password"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "go_database"

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

  // レコード作成
//   func CreateGacha(name) {
// 	db := sqlConnect()
// 	db.Create(&Gacha{Name: name, Used: 0})
// 	defer db.Close()
//   }

  // ルーティング
  func main() {
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
		// CreateGacha(ctx.PostForm("name"))
		name := ctx.PostForm("name")
		db.Create(&Gacha{Name: name, Used: 0})
		defer db.Close()

		ctx.Redirect(302, "/")
	})

	// 目的地を取得
	// TODO: 目的地名の取得ロジックは後程実装
	router.GET("/destination", func(ctx *gin.Context) {
		destination_name:= ""
		if destination_name == "" {
			ctx.HTML(200, "error.html", gin.H{})
		} else {
			ctx.HTML(200, "destination.html", gin.H{"destination": destination_name})
		}
	})

	// 再度ガチャ
	router.GET("/retry", func(ctx *gin.Context) {
		ctx.HTML(200, "form.html", gin.H{})
	})

	router.Run()
  }
