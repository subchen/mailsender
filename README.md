# mailsender

A tool which can send html email

## Usage

```
NAME:
   mailsender - Send mail using smtp server

USAGE:
   mailsender [ OPTIONS ... ] html-body-file

VERSION:
   0.0.0

AUTHORS:
   Guoqiang Chen <subchen@gmail.com>

OPTIONS:
   --from value          sender email address: "My Name <my_name@example.com>"
   --to value            recipient email address: "Your Name <your_name@example.com>". (multiple values)
   --subject value       subject of mail
   --attach value        attachment file. (multiple values)
   --smtp-server value   smtp server hostname
   --smtp-port value     smtp server port (default: 25)
   --smtp-ssl            enable ssl for smtp server (default: false)
   --smtp-user value     username for smtp server
   --smtp-pass value     password for smtp server
   --help                print this usage
   --version             print version information
```

## Examples

```bash
mailsender \
  --from="alex@gmail.com" \
  --to="bob@gmail.com" \
  --to="jack@gmail.com" \
  --subject="test message" \
  --attach="some.doc" \
  --smtp-server="smtp.gmail.com" \
  --smtp-port="465" \
  --smtp-user="alex"
  --smtp-pass="password"
  email.html
```

> Notes: Images (example: `<img src="cid:abc.jpg">`) will be automatically emebed into email.
