
# Send email with Rest API or CLI

Simple way to send email over **Rest API** and **CLI**.

Build app (golang required):
```bash
go get -d
go build -o mailer .
mailer --help
```

# CLI options
type `mailer --help` to display all available options
```
-attachment value
      Set email attachments
-bcc value
      BCC, define multiple to send more than one bcc
-body string
      Set email body
-cc value
      CC, define multiple to send more than one cc
-debug
      enable debug (default false)
-endpoint string
      Start http with endpoint (default "/")
-from string
      Set Sender email
-help
      Display this help
-host string
      Set email host
-http
      Enable http server
-password string
      Set email password
-plain-text
      Send email as plaintext. (default false)
-port int
      Start http with port (default 8080)
-provider string
      Set email provider, eg: gmail
-single
      single email per receiver. (default false)
-subject string
      Set email subject
-to value
      Receiver Email, define multiple to send more than one receiver
-username string
      Set email username
```

# Send email using CLI
Example:
```bash
mailer --provider gmail \
  --from my@gmail.com \
  --password s3cr3t \
  --to receiver1@mail.com \
  --to receiver2@mail.com \
  --subject "Send from cli" \
  --body "Hello <b>world</b>" \
  --attachment file1.txt \
  --attachment file2.txt
```
Don't forget to enable less secure apps to your gmail as a sender
[Enable less secure apps](https://support.google.com/a/answer/6260879)

# Send email using Rest API
You must start http server before send email using Rest API:
```
mailer --http --port 8765 --endpoint /api/v1/mailer
```

Send email using curl:
```
curl -X POST \
  -H 'Content-type: application/json' \
  -d '{
        "from": "my@gmail.com",
        "password": "s3cr3t",
        "to": [
            "receiver@gmail.com"
        ],
        "subject": "Send from Rest API",
        "body": "Hello <b>world</b>",
        "attachments": [
            {
                "name": "image.png",
                "data": "data:image/png,base64,TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQ="
            }
        ]
    }' \
  http://localhost:8765/api/v1/mailer
```
