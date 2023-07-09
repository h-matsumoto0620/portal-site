package main

import (
	"fmt"
	"net/http"
	"os"
	"portal-site/crypto"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

type User struct {
	gorm.Model
	Username string `form:"username" binding:"required" gorm:"unique;not null"`
	Password string `form:"password" binding:"required"`
}

type Project struct {
	gorm.Model
	UserId       uint
	ProjectName  string `form:"projectName" binding:"required"`
	ProjectStart string `form:"projectStart" binding:"required"`
	ProjectEnd   string `form:"projectEnd" binding:"required"`
	Content      string `form:"content"`
	ProgramLang  string `form:"programLang" gorm:"type:text"`
	OS           string `form:"os"`
	Env          string `form:"env" gorm:"type:text"`
	WorkInCharge string `form:"workInCharge"`
	Occupation   string `form:"occupation"`
	Memo         string `form:"memo" gorm:"type:text"`
}

func gormConnect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	DBMS := os.Getenv("DBMS")
	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	DB := os.Getenv("DB")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")

	CONNECT := DBUser + ":" + DBPass + "@tcp(" + HOST + ":" + PORT + ")/" + DB + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func dbInit() {
	db := gormConnect()
	defer db.Close()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Project{})
}

// ユーザー登録処理
func createUser(username string, password string) []error {
	passwordEncrypt, _ := crypto.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&User{Username: username, Password: passwordEncrypt}).GetErrors(); err != nil {
		return err
	}
	return nil
}

// ユーザーを一件取得
func getUser(username string) User {
	db := gormConnect()
	var user User
	db.First(&user, "username = ?", username)
	db.Close()
	return user
}

// ログインセッション登録
func setLoginSessions(ctx *gin.Context, UserData User) {
	session := sessions.Default(ctx)
	session.Set("UserID", UserData.ID)
	session.Save()
}

// ログインセッション確認
func sessionCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userId := session.Get("UserID")

		//　セッションがない場合、ログインフォームに移動
		if userId == nil {
			fmt.Println("ログインされていません")
			ctx.Redirect(http.StatusMovedPermanently, "/login")
			ctx.Abort()
		} else {
			ctx.Set("UserID", userId)
			ctx.Next()
		}
		fmt.Println("ログインチェック終了")
	}
}

// 案件情報登録
func createProject(ctx *gin.Context, project Project) []error {
	db := gormConnect()
	defer db.Close()

	session := sessions.Default(ctx)
	userId := session.Get("UserID")

	if err := db.Create(&Project{UserId: userId.(uint), ProjectName: project.ProjectName, ProjectStart: project.ProjectStart, ProjectEnd: project.ProjectEnd, Content: project.Content, ProgramLang: project.ProgramLang, OS: project.OS, Env: project.Env, WorkInCharge: project.WorkInCharge, Occupation: project.Occupation, Memo: project.Memo}).GetErrors(); err != nil {
		return err
	}
	return nil
}

// 案件情報取得
func getProjectData(ctx *gin.Context) []Project {
	session := sessions.Default(ctx)
	userId := session.Get("UserID")

	db := gormConnect()
	var projectList []Project
	db.Find(&projectList, "user_id = ?", userId)

	return projectList
}

func main() {
	// DBの初期設定
	dbInit()

	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")
	router.Static("/assets", "./assets")

	// セッションの設定
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// ユーザー登録画面
	router.GET("/signup", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "signup.html", gin.H{})
	})

	// ユーザー登録
	router.POST("/signup", func(ctx *gin.Context) {
		var form User
		// バリデーション処理
		if err := ctx.Bind(&form); err != nil {
			ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
			ctx.Abort()
		} else {
			username := ctx.PostForm("username")
			password := ctx.PostForm("password")
			// 登録ユーザーが重複していた場合にはじく処理
			if err := createUser(username, password); len(err) != 0 {
				ctx.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
			}
			ctx.Redirect(302, "/login")
		}
	})

	// ユーザーログイン画面
	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})

	// ユーザーログイン
	router.POST("/login", func(ctx *gin.Context) {

		// DBから取得したユーザーパスワード(Hash)
		userData := getUser(ctx.PostForm("username"))
		dbPassword := userData.Password
		// フォームから取得したユーザーパスワード
		formPassword := ctx.PostForm("password")

		// ユーザーパスワードの比較
		if err := crypto.CompareHashAndPassword(dbPassword, formPassword); err != nil {
			fmt.Println("ログインできませんでした")
			ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
			ctx.Abort()
		} else {
			fmt.Println("ログインできました")
			setLoginSessions(ctx, userData)
			ctx.Redirect(302, "/dashboard")
		}
	})

	dashboard := router.Group("/dashboard")
	dashboard.Use(sessionCheck())
	{
		dashboard.GET("", func(ctx *gin.Context) {
			projectData := getProjectData(ctx)
			ctx.HTML(http.StatusOK, "dashboard.html", gin.H{"projectData": projectData})
		})
	}

	register := router.Group("/register")
	register.Use(sessionCheck())
	{
		register.GET("", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "projectRegister.html", gin.H{})
		})
		register.POST("", func(ctx *gin.Context) {
			var pro Project
			if err := ctx.Bind(&pro); err != nil {
				ctx.HTML(http.StatusBadRequest, "projectRegister.html", gin.H{"err": err})
				ctx.Abort()
			} else {
				if err := createProject(ctx, pro); len(err) != 0 {
					ctx.HTML(http.StatusBadRequest, "projectRegister.html", gin.H{"err": err})
				}
				ctx.Redirect(302, "/dashboard")
			}
		})
	}

	// router.GET("/icons", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "icons.html", gin.H{})
	// })

	router.Run(":8080")
}
