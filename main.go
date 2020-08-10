package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
        router.LoadHTMLGlob("views/*.html")
	    router.Static("/assets", "./assets")

    router.GET("/", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "dashboard.html", gin.H{})
    })
    
    router.GET("/user.html", func(ctx *gin.Context){
        ctx.HTML(http.StatusOK, "user.html", gin.H{})
    })

    router.Run(":8080")
}
