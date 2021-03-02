package smtp

import (
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"github.com/zhashkevych/courses-backend/pkg/email"
)

type SMTPSender struct {
	from string
	pass string
	host string
	port int
}

func NewSMTPSender(from string, pass string, host string, port int) *SMTPSender {
	return &SMTPSender{from: from, pass: pass, host: host, port: port}
}

func (s *SMTPSender) Send(input email.SendEmailInput) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.from)
	msg.SetHeader("To", input.To)
	msg.SetHeader("Subject", input.Subject)
	msg.SetBody("text/html", input.Body)

	dialer := gomail.NewDialer(s.host, s.port, s.from, s.pass)
	if err := dialer.DialAndSend(msg); err != nil {
		return errors.Wrap(err, "failed to sent email via smtp")
	}

	// Send the email to Bob
	//err := gomail.Send(d, msg)
	//if err := d.DialAndSend(msg); err != nil {
	//	panic(err)
	//}
	//
	//addr := fmt.Sprintf("%s:%d", s.host, s.port)
	//auth := smtp.PlainAuth("", s.from, s.pass, s.host)
	//if err := smtp.SendMail(addr, auth, s.from, []string{input.To}, msg); err != nil {
	//	return errors.Wrap(err, "failed to sent email via smtp")
	//}

	return nil
}
