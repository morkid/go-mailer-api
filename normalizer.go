package main

import (
    "reflect"
)

func NormalizeAddress(mails []interface{}) []interface{} {
    normalized := []interface{}{}
    var stringMap map[string]interface{}
    if nil != mails {
        for _, o := range mails {
            if reflect.TypeOf(o) == reflect.TypeOf("") {
                normalized = append(normalized, map[string]interface{}{
                    "address": o,
                    "name": "",
                })
            } else if reflect.TypeOf(o) == reflect.TypeOf(stringMap) {
                normalized = append(normalized, o.(interface{}))
            }
        }
    }

    return normalized
}

func Normalize(config *MailConfig) {
    config.To = NormalizeAddress(config.To)
    config.Cc = NormalizeAddress(config.Cc)
    config.Bcc = NormalizeAddress(config.Bcc)
    if nil == config.Attachments {
        config.Attachments = []interface{}{}
    }
    Provider(config)
}
