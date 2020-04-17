package main

import (
    // "fmt"
    "path/filepath"
    "path"
    "reflect"
    "os"
    "strconv"
    "time"
    "io/ioutil"
)

func NormalizeAddress(mails []interface{}) []interface{} {
    normalized := []interface{}{}
    if nil != mails {
        for _, o := range mails {
            if reflect.TypeOf(o) == reflect.TypeOf("") {
                normalized = append(normalized, map[string]interface{}{
                    "address": o,
                    "name": "",
                })
            } else if reflect.TypeOf(o) == reflect.TypeOf(map[string]interface{}{}) {
                normalized = append(normalized, o.(interface{}))
            }
        }
    }

    return normalized
}

func NormalizeAttachment(attachments []interface{}, cli bool) []interface{} {
    normalized := []interface{}{}
    basepath, e := ioutil.TempDir("", "go_mailer_" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
    if e != nil {
        return normalized
    }

    if nil != attachments {
        for i, o := range attachments {
            if reflect.TypeOf(o) == reflect.TypeOf("") {
                var value string = o.(string)
                if cli && len(value) < 1000 {
                    fullpath, er := filepath.Abs(value)
                    if nil == er {
                        if _, e := os.Stat(fullpath); os.IsNotExist(e) {
                            normalized = append(normalized, fullpath)
                            continue
                        }
                    }
                }
            } else if reflect.TypeOf(o) == reflect.TypeOf(map[string]interface{}{}) {
                attachment, _ := o.(map[string]interface{})
                var fullpath = ""

                if nil != attachment["data"] {
                    var data = attachment["data"].(string)

                    if data != "" {
                        var name = ""
                        if nil != attachment["name"] && "" != attachment["name"].(string) {
                            name = attachment["name"].(string)
                        } else {
                            name = "file-" + strconv.Itoa(i) + ".txt"
                        }
                        fullpath = path.Join(basepath, name)
                        err := ioutil.WriteFile(fullpath, []byte(data), 0644)
                        if nil != err {
                            continue
                        }
                    }
                }

                if fullpath != "" {
                    if _, e := os.Stat(fullpath); os.IsNotExist(e) {
                        normalized = append(normalized, fullpath)
                    }
                }
            }
        }
    }

    return normalized
}

func Normalize(config *MailConfig) {
    config.To = NormalizeAddress(config.To)
    config.Cc = NormalizeAddress(config.Cc)
    config.Bcc = NormalizeAddress(config.Bcc)
    config.Attachments = NormalizeAttachment(config.Attachments, config.Cli)
    if nil == config.Attachments {
        config.Attachments = []interface{}{}
    }
    Provider(config)
}
