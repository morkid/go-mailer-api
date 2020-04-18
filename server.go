package main

import (
    "github.com/gin-gonic/gin"
    "strconv"
    "fmt"
    "path"
)

func StartServer(debug bool) {
    if (!debug) {
        gin.SetMode(gin.ReleaseMode)
    }
    r := gin.Default()

    r.GET("/", func (c *gin.Context) {
        c.JSON(501, gin.H{
            "message": "Not implemented",
        })
    })

    r.POST("/", func (c *gin.Context) {
        var config MailConfig
        var redirect = ""
        config.Thread = true
        form, e := c.MultipartForm()
        if e == nil {
            var basepath = GetTmpDir()

            if len(form.Value["redirect"]) > 0 {
                redirect = form.Value["redirect"][0]
            }

            if len(form.Value["provider"]) > 0 {
                config.Provider = form.Value["provider"][0]
            }

            if len(form.Value["host"]) > 0 {
                config.Host = form.Value["host"][0]
            }

            if len(form.Value["port"]) > 0 {
                a, _:= strconv.Atoi(form.Value["port"][0])
                config.Port = a
            }

            if len(form.Value["from"]) > 0 {
                config.From = form.Value["from"][0]
            }

            if len(form.Value["to"]) > 0 {
                to, _ := form.Value["to"]
                config.To = make([]interface{}, len(to))
                for i, v := range to {
                    config.To[i] = v
                }
            }

            if len(form.Value["cc"]) > 0 {
                cc, _ := form.Value["cc"]
                config.Cc = make([]interface{}, len(cc))
                for i, v := range cc {
                    config.Cc[i] = v
                }
            }

            if len(form.Value["bcc"]) > 0 {
                bcc, _ := form.Value["bcc"]
                config.Bcc = make([]interface{}, len(bcc))
                for i, v := range bcc {
                    config.Bcc[i] = v
                }
            }

            if len(form.Value["subject"]) > 0 {
                config.Subject = form.Value["subject"][0]
            }

            if len(form.Value["body"]) > 0 {
                config.Body = form.Value["body"][0]
            }

            if len(form.Value["username"]) > 0 {
                config.Username = form.Value["username"][0]
            }

            if len(form.Value["password"]) > 0 {
                config.Password = form.Value["password"][0]
            }

            if len(form.Value["plain_text"]) > 0 {
                var t = form.Value["plain_text"][0]
                if t == "1" || t == "true" || t == "TRUE" || t == "on" {
                    config.PlainText = true
                }
            }

            if len(form.Value["single"]) > 0 {
                var t = form.Value["single"][0]
                if t == "1" || t == "true" || t == "TRUE" || t == "on" {
                    config.Single = true
                }
            }

            config.Attachments = []interface{}{}
            if len(form.File["attachments"]) > 0 {
                for _, file := range form.File["attachments"] {
                    var fullpath = path.Join(basepath, file.Filename)
                    if err := c.SaveUploadedFile(file, fullpath); err != nil {
                        fmt.Println(err)
                    } else {
                        config.Attachments = append(config.Attachments, fullpath)
                    }
                }
            }

            config.SkipAttachmentCheck = true
        } else {
            config.SkipAttachmentCheck = false
            c.BindJSON(&config)
        }

        if (debug) {
            SendMail(&config)
            if redirect != "" {
                c.Redirect(302, redirect)
            } else {
                c.JSON(200, &config)
            }
        } else {
            go SendMail(&config)
            if redirect != "" {
                c.Redirect(302, redirect)
            } else {
                c.JSON(200, gin.H{
                    "message": "sent",
                })
            }
        }
    })

    r.Run()
}
