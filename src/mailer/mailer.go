package main

import (
	"bytes"
	"github.com/vanng822/go-premailer/premailer"
	"github.com/xhit/go-simple-mail/v2"
	"html/template"
	"log"
	"time"
)

const (
	connTimeout       = 10 * time.Second
	sendTimeout       = 10 * time.Second
	keepAlive         = false
	htmlTemplatePath  = "./templates/mail.html.gohtml"
	plainTemplatePath = "./templates/mail.plain.gohtml"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	// Create SMTP server for sending email
	smtpServer := mail.NewSMTPClient()
	smtpServer.Host = m.Host
	smtpServer.Port = m.Port
	smtpServer.Username = m.Username
	smtpServer.Password = m.Password
	smtpServer.Encryption = m.GetEncryption(m.Encryption)
	smtpServer.KeepAlive = keepAlive
	smtpServer.ConnectTimeout = connTimeout
	smtpServer.SendTimeout = sendTimeout

	log.Println("encryption ", m.Encryption)
	log.Println("smtp encryption ", smtpServer.Encryption)

	htmlMessage, err := m.createHtmlMessage(msg)
	if err != nil {
		log.Println("failed to create html message...", err)
		return err
	}

	plaintextMessage, err := m.createPlaintextMessage(msg)
	if err != nil {
		log.Println("failed to create plaintext message...", err)
		return err
	}

	// Connect to SMTP server
	// Create new SMTP client
	smtpClient, err := smtpServer.Connect()
	if err != nil {
		log.Println("failed to establish conn with smtp server...", err)
		return err
	}

	// Create new Email
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plaintextMessage)
	email.AddAlternative(mail.TextHTML, htmlMessage)

	if len(msg.Attachments) > 0 {
		for _, a := range msg.Attachments {
			email.AddAttachment(a)
		}
	}

	// Send email over SMTP client
	err = email.Send(smtpClient)
	if err != nil {
		log.Println("failed to send email...", err)
		return err
	}

	return nil
}

func (m *Mail) GetEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

func (m *Mail) createHtmlMessage(msg Message) (string, error) {
	templ, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		log.Println("failed to parse email html template...", err)
		return "", err
	}

	var (
		// Create a bytes buffer out of template
		templBytesBuffer bytes.Buffer
	)
	// Apply template and write output to writer
	if err = templ.ExecuteTemplate(&templBytesBuffer, "body", msg.DataMap); err != nil {
		log.Println("failed to execute html template...", err)
		return "", err
	}

	htmlMessage := templBytesBuffer.String()
	htmlMessage, err = m.inlineCss(htmlMessage)
	if err != nil {
		log.Println("failed to inline css to html email...", err)
		return "", err
	}

	return htmlMessage, nil
}

func (m *Mail) createPlaintextMessage(msg Message) (string, error) {
	templ, err := template.ParseFiles(plainTemplatePath)
	if err != nil {
		log.Println("failed to parse email plaintext template...", err)
		return "", nil
	}

	var templBytesBuffer bytes.Buffer
	if err = templ.ExecuteTemplate(&templBytesBuffer, "body", msg.DataMap); err != nil {
		log.Println("failed to execute plaintext template...", err)
		return "", err
	}

	plaintextMessage := templBytesBuffer.String()

	return plaintextMessage, nil
}

// Css should be inlined when dealing with emails
// See github.com/premailer
func (m *Mail) inlineCss(s string) (string, error) {
	opts := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	p, err := premailer.NewPremailerFromString(s, &opts)
	if err != nil {
		log.Println("failed to create premailer from string...", err)
		return "", err
	}

	res, err := p.Transform()
	if err != nil {
		log.Println("failed to transform premailer to result string...", err)
		return "", err
	}

	return res, nil
}
