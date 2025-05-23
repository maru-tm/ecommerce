package email

import (
	"api-gateway/config"
	"fmt"
	"net/smtp"
)

// Mailer структура с конфигурацией
type Mailer struct {
	From     string
	Password string
	Host     string
	Port     string
}

// NewMailer создаёт новый Mailer на основе конфигурации
func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{
		From:     cfg.EmailFrom,
		Password: cfg.EmailPassword,
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
	}
}

// SendWelcomeEmail отправляет приветственное письмо
func (m *Mailer) SendWelcomeEmail(toEmail, fullName string) error {
	auth := smtp.PlainAuth("", m.From, m.Password, m.Host)

	subject := "Добро пожаловать в Bookstore!"
	body := fmt.Sprintf(`
Добро пожаловать, %s!

Спасибо за регистрацию на нашей платформе.

С любовью, команда Bookstore 💙
`, fullName)

	msg := []byte("Subject: " + subject + "\r\n" +
		"From: Bookstore <" + m.From + ">\r\n" +
		"To: " + toEmail + "\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
		body)

	return smtp.SendMail(m.Host+":"+m.Port, auth, m.From, []string{toEmail}, msg)
}
