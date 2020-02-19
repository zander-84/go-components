package CLoggerZapEmail

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

type EmailHook struct {
	Host     string
	Port     int
	User     string
	Password string
	To       []string
	Subject  string
}

func (l *EmailHook) Write(p []byte) (n int, err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", l.User)
	m.SetHeader("To", l.To...)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", l.Subject+"-log")
	//m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	m.SetBody("text/html", string(p))
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(l.Host, l.Port, l.User, l.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		//panic(err)
	}

	n = len(p)
	return n, err
}
