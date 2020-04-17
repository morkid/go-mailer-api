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
            "message": "not implemented",
        })
    })

    r.POST("/", func (c *gin.Context) {
        var config MailConfig
        c.BindJSON(&config)
        config.Thread = true
        config.SkipAttachmentCheck = false
        if (debug) {
            SendMail(&config)
            c.JSON(200, &config)
        } else {
            go SendMail(&config)
            c.JSON(200, gin.H{
                "message": "sent",
            })
        }
    })

    r.Run()
}
