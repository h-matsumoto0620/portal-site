package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
        router.LoadHTMLGlob("views/*.html")
	    router.Static("/assets", "./assets")

    router.GET("/dashboard", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "dashboard.html", gin.H{})
    })
    
    router.GET("/icons", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "icons.html", gin.H{})
    })

    router.GET("/map", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "map.html", gin.H{})
    })

    router.GET("/notifications", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "notifications.html", gin.H{})
    })

    router.GET("/tables", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "tables.html", gin.H{})
    })

    router.GET("/upgrade", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "upgrade.html", gin.H{})
    })

    router.GET("/user", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "user.html", gin.H{})
    })

    router.Run(":8080")
}
