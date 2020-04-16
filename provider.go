package main

func Provider(c *MailConfig) {
    switch c.Provider {
        case "gmail":
            GmailProvider(c)
    }
}

func GmailProvider(config *MailConfig) {
    config.Host = "smtp.gmail.com"
    config.Port = 587
}
