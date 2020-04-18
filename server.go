package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "mime/multipart"
    "path"
    "strings"
    "log"
    "os"
    "strconv"
)

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func HttpLogging(res http.ResponseWriter, req *http.Request, status int) {
    var message = "OK"
    if status >= 300 {
        switch status {
            case 400: message = "bad request"
            case 404: message = "not found"
            case 415: message = "unsupported media type"
            case 501: message = "not implemented"
            case 500: message = "internal server error"
            default:
                message = "internal server error"
                status = 500
        }
        http.Error(res, message, status)
    }
    log.Println(status, req.Method, "\"" + req.RequestURI + "\"")
}



func StartServer(endpoint string, debug bool) {
    if endpoint == "" {
        endpoint = "/"
    }
    http.HandleFunc(endpoint, func(res http.ResponseWriter, req *http.Request) {
        if req.URL.Path != "/" {
            HttpLogging(res, req, 404)
            return
        }
        if req.Method != "POST" {
            HttpLogging(res, req, 501)
            return
        }

        var contentType = strings.ToLower(req.Header.Get("Content-type"))
        config := MailConfig{
            Thread: true,
        }
        var redirect = ""

        if strings.HasPrefix(contentType, "multipart/form-data") {
            err := req.ParseMultipartForm(32 << 20)
            if nil == err {
                var basepath = GetTmpDir()
                form := req.MultipartForm
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
                        if err := SaveUploadedFile(file, fullpath); err != nil {
                            log.Println(err)
                        } else {
                            config.Attachments = append(config.Attachments, fullpath)
                        }
                    }
                }
                config.SkipAttachmentCheck = true
            } else {
                HttpLogging(res, req, 400)
                return
            }
        } else if strings.HasPrefix(contentType, "application/json") {
            body, _ := ioutil.ReadAll(req.Body)
            json.Unmarshal(body, &config)
            config.SkipAttachmentCheck = false
        } else {
            HttpLogging(res, req, 415)
            return
        }

        if redirect != "" {
            http.Redirect(res, req, redirect, 302)
            return
        }

        HttpLogging(res, req, 200)
        res.Header().Add("Content-type", "application/json")

        if (debug) {
            SendMail(&config)
            if redirect != "" {
                http.Redirect(res, req, redirect, 302)
            } else {
                b, _ := json.Marshal(config)
                fmt.Fprintf(res, string(b))
            }
        } else {
            go SendMail(&config)
            if redirect != "" {
                http.Redirect(res, req, redirect, 302)
            } else {
                fmt.Fprintf(res, "{\"message\":\"sent\"}")
            }
        }
    })

    var port = os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Println("Server started at port "+port)
    http.ListenAndServe(":" + port, nil)
}
