package main

import (
    "github.com/gin-gonic/gin"
)

func StartServer(debug bool) {
    if (!debug) {
        gin.SetMode(gin.ReleaseMode)
    }
    r := gin.Default()

    r.GET("/", func (c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello world",
        })
    })

    r.POST("/", func (c *gin.Context) {
        var config MailConfig
        c.BindJSON(&config)
        SendMail(&config, false)
        c.JSON(200, config)
    })

    r.Run()
}
