package main

type MailConfig struct {
    Provider string                 `json:"provider,omitempty"`
    Host string                     `json:"host,omitempty"`
    Port int                        `json:"port,omitempty"`
    From string                     `json:"from,omitempty"`
    To []interface{}                `json:"to"`
    Cc []interface{}                `json:"cc"`
    Bcc []interface{}               `json:"bcc"`
    Subject string                  `json:"subject,omitempty"`
    Body string                     `json:"body,omitempty"`
    Attachments []interface{}       `json:"attachments"`
    Username string                 `json:"username,omitempty"`
    Password string                 `json:"password,omitempty"`
    PlainText bool                  `json:"plain_text"`
    Single bool                     `json:"single"`
    Thread bool                     `json:"-"`
    // Cli bool                        `json:"-"`
    SkipAttachmentCheck bool        `json:"-"`
}
