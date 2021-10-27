package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

// 自定义邮件类
type Email struct {
	*SMTPInfo
}

//参数
type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	PassWord string
	From     string
}

// 创建自定义邮件
func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

// SendEmail 发送邮件
func (e *Email) SendEmail(subject, body string, to []string, Cc []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)        // 收件人
	m.SetHeader("Subject", subject) // 主题
	m.SetHeader("Cc", Cc...)        //抄送
	m.SetBody("text/html", body)    // 正文

	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.PassWord)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m)
}
