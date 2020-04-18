package main

import (
    "path/filepath"
    "path"
    "encoding/base64"
    "regexp"
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

func RandomNumber() string {
    return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

func GetTmpDir() string {
    basepath, e := ioutil.TempDir("", "go_mailer_" + RandomNumber())
    if e != nil {
        return ""
    }
    return basepath
}

func NormalizeAttachment(attachments []interface{}, skip bool) []interface{} {
    normalized := []interface{}{}
    var basepath = GetTmpDir()
    if basepath == "" {
        return normalized
    }

    if nil != attachments {
        for index, o := range attachments {
            if reflect.TypeOf(o) == reflect.TypeOf("") {
                var value string = o.(string)
                if skip && len(value) < 1000 {
                    fullpath, er := filepath.Abs(value)
                    if nil == er {
                        if _, e := os.Stat(fullpath); !os.IsNotExist(e) {
                            normalized = append(normalized, fullpath)
                        }
                    }
                    continue
                }

                var fullpath = ConvertDataToFilePath(index, basepath, map[string]interface{}{
                    "data": value,
                })
                if fullpath != "" {
                    if _, e := os.Stat(fullpath); !os.IsNotExist(e) {
                        normalized = append(normalized, fullpath)
                    }
                }
            } else if reflect.TypeOf(o) == reflect.TypeOf(map[string]interface{}{}) {
                attachment, _ := o.(map[string]interface{})
                var fullpath = ConvertDataToFilePath(index, basepath, attachment)

                if fullpath != "" {
                    if _, e := os.Stat(fullpath); !os.IsNotExist(e) {
                        normalized = append(normalized, fullpath)
                    }
                }
            }
        }
    }

    return normalized
}

func ConvertDataToFilePath(index int, basepath string, attachment map[string]interface{}) string {
    var fullpath = ""

    if nil != attachment["data"] {
        var data = attachment["data"].(string)

        if data != "" {
            var name = ""
            var extname = "txt"
            pattern, errp := regexp.Compile(`^data:([^/]+)/([^;]+);base64,(.*)`)
            var decoded = false
            var bytedata = []byte(data)
            if nil == errp {
                matched := pattern.FindStringSubmatch(data)
                if len(matched) == 4 {
                    extname = matched[2]
                    dec, errd := base64.StdEncoding.DecodeString(matched[3])
                    if nil == errd {
                        bytedata = dec
                        decoded = true
                    }
                }
            }
            if !decoded {
                dec, errd := base64.StdEncoding.DecodeString(data)
                if nil == errd {
                    bytedata = dec
                    decoded = true
                }
            }
            if !decoded {
                return ""
            }
            if nil != attachment["name"] && "" != attachment["name"].(string) {
                name = attachment["name"].(string)
            } else {
                name = "file-" + strconv.Itoa(index) + "." + extname
            }
            fullpath = path.Join(basepath, name)
            err := ioutil.WriteFile(fullpath, bytedata, 0777)
            if nil != err {
                return ""
            }
        }
    }

    return fullpath
}

func Normalize(config *MailConfig) {
    config.To = NormalizeAddress(config.To)
    config.Cc = NormalizeAddress(config.Cc)
    config.Bcc = NormalizeAddress(config.Bcc)
    config.Attachments = NormalizeAttachment(config.Attachments, config.SkipAttachmentCheck)
    if nil == config.Attachments {
        config.Attachments = []interface{}{}
    }
    Provider(config)
}
