package main

import (
    "fmt"
    "os"
    "flag"
    "strconv"
    "encoding/json"
)

type StringFlags []string

func (i *StringFlags) String() string {
    return "my string representation"
}

func (i *StringFlags) Set(value string) error {
    *i = append(*i, value)
    return nil
}

func main() {
    var tos StringFlags
    var ccs StringFlags
    var bccs StringFlags
    var attachments StringFlags

    httpEnabled := flag.Bool("http", false, "Enable http server")
    port := flag.Int("port", 8080, "Start http with port")
    help := flag.Bool("help", false, "Display this help")
    debug := flag.Bool("debug", false, "enable debug (default false)")

    provider := flag.String("provider", "", "Set email provider, eg: gmail")
    host := flag.String("host", "", "Set email host")
    from := flag.String("from", "", "Set Sender email")
    subject := flag.String("subject", "", "Set email subject")
    body := flag.String("body", "", "Set email body")
    username := flag.String("username", "", "Set email username")
    password := flag.String("password", "", "Set email password")
    plaintext := flag.Bool("plain-text", false, "Send email as plaintext. (default false)")
    single := flag.Bool("single", false, "single email per receiver. (default false)")
    flag.Var(&tos, "to", "Receiver Email, define multiple to send more than one receiver")
    flag.Var(&ccs, "cc", "CC, define multiple to send more than one cc")
    flag.Var(&bccs, "bcc", "BCC, define multiple to send more than one bcc")
    flag.Var(&attachments, "attachment", "Set email attachments")

    flag.Parse()

    if *help {
        flag.PrintDefaults()
        os.Exit(0)
    }

    if *httpEnabled {
        if *port != 8080 {
            os.Setenv("PORT", strconv.Itoa(*port))
        }
        StartServer(*debug)
    } else {
        if tos == nil {
            flag.PrintDefaults()
        }
        var tosStrings = []interface{}{}
        var ccsStrings = []interface{}{}
        var bccsStrings = []interface{}{}
        var attachmentsStrings = []interface{}{}
        for _, o := range tos {
            tosStrings = append(tosStrings, o)
        }
        for _, o := range ccs {
            ccsStrings = append(ccsStrings, o)
        }
        for _, o := range bccs {
            bccsStrings = append(bccsStrings, o)
        }
        for _, o := range attachments {
            attachmentsStrings = append(attachmentsStrings, o)
        }

        config := MailConfig{
            Provider: *provider,
            Host: *host,
            Port: *port,
            From: *from,
            To: tosStrings,
            Cc: ccsStrings,
            Attachments: attachmentsStrings,
            Bcc: bccsStrings,
            Subject: *subject,
            Body: *body,
            Username: *username,
            Password: *password,
            PlainText: *plaintext,
            Single: *single,
            SkipAttachmentCheck: true,
            Thread: false,
        }
        if (*debug) {
            Normalize(&config)
            j, _ := json.MarshalIndent(&config, "", "    ")
            fmt.Println(string(j))
        } else {
            SendMail(&config)
        }
    }
}
