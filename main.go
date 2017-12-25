package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/subchen/go-cli"
	"github.com/ungerik/go-dry"
)

// version
var (
	BuildVersion   string
	BuildGitRev    string
	BuildGitCommit string
	BuildDate      string
)

// flags
var (
	from       string
	recipients []string
	subject    string
	embedes    []string
	attaches   []string

	smtpServerHost string
	smtpServerPort int
	smtpUser       string
	smtpPass       string
	smtpSSL        bool
)

func main() {
	app := cli.NewApp()
	app.Name = "mailsender"
	app.Usage = "Send mail using smtp server"
	app.UsageText = "[ OPTIONS ... ] html-body-file"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"
	app.FlagsAlign = false

	app.Flags = []*cli.Flag{
		{
			Name:   "from",
			Usage:  `sender email address: "My Name <my_name@example.com>"`,
			EnvVar: "MAIL_FROM",
			Value:  &from,
		},
		{
			Name:   "to",
			Usage:  `recipient email address: "Your Name <your_name@example.com>". (multiple values)`,
			EnvVar: "MAIL_TO",
			Value:  &recipients,
		},
		{
			Name:  "subject",
			Usage: "subject of mail",
			Value: &subject,
		},
		{
			Name:  "attach",
			Usage: "attachment file. (multiple values)",
			Value: &attaches,
		},
		{
			Name:  "embed",
			Usage: `embed file for html body, example: <img src="cid:image.jpg">. (multiple values)`,
			Value: &embedes,
		},
		{
			Name:   "smtp-server",
			Usage:  "smtp server hostname",
			EnvVar: "SMTP_SERVER",
			Value:  &smtpServerHost,
		},
		{
			Name:     "smtp-port",
			Usage:    "smtp server port",
			DefValue: "25",
			EnvVar:   "SMTP_PORT",
			Value:    &smtpServerPort,
		},
		{
			Name:     "smtp-ssl",
			Usage:    "enable ssl for smtp server",
			DefValue: "false",
			EnvVar:   "SMTP_SSL",
			Value:    &smtpSSL,
		},
		{
			Name:   "smtp-user",
			Usage:  "username for smtp server",
			EnvVar: "SMTP_USER",
			Value:  &smtpUser,
		},
		{
			Name:   "smtp-pass",
			Usage:  "password for smtp server",
			EnvVar: "SMTP_PASS",
			Value:  &smtpPass,
		},
	}

	if BuildVersion != "" {
		app.Version = BuildVersion + "-" + BuildGitRev
	}
	app.BuildGitCommit = BuildGitCommit
	app.BuildDate = BuildDate

	app.Action = func(c *cli.Context) {
		if c.NArg() != 1 {
			c.ShowHelp()
			return
		}

		defer func() {
			if err := recover(); err != nil {
				os.Stderr.WriteString(fmt.Sprintf("fatal: %v\n", err))
				os.Exit(1)
			}
		}()

		bodyfile := c.Args()[0]
		if !dry.FileExists(bodyfile) {
			panic("body file not found: " + bodyfile)
		}
		body, err := ioutil.ReadFile(bodyfile)
		if err != nil {
			panic(err)
		}

		m := gomail.NewMessage()
		m.SetHeader("From", from)
		m.SetHeader("To", recipients...)
		m.SetHeader("Subject", subject)
		m.SetDateHeader("X-Date", time.Now())
		m.SetBody("text/html; charset=utf-8", string(body))

		for _, file := range embedes {
			if !dry.FileExists(file) {
				panic("embeded file not found: " + file)
			}
			m.Embed(file)
		}
		for _, file := range attaches {
			if !dry.FileExists(file) {
				panic("attached file not found: " + file)
			}
			m.Attach(file)
		}

		d := gomail.Dialer{
			Host:     smtpServerHost,
			Port:     smtpServerPort,
			Username: smtpUser,
			Password: smtpPass,
			SSL:      smtpServerPort == 465 || smtpSSL,
			TLSConfig: &tls.Config{
				ServerName:         smtpServerHost,
				InsecureSkipVerify: true,
			},
		}

		fmt.Printf("Sending mail to %s ...\n", strings.Join(recipients, ", "))

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}

	app.Run(os.Args)
}
