package mail_sender

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

func Mail_sender(email, url_accepter, code, mail_user, mail_pass string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", mail_user)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Подвержение почты")
	url_code := url_accepter + "/" + code
	body_text := fmt.Sprintf(`<h1>Привет!</h1>
			<p>Это письмо отправленно для подверждения вашего аккаунта:<a href="%s">Подвердить</a> </p>
			<p>С уважением,<br/>Ваше Sso_service</p>`, url_code)
	m.SetBody("text/html", body_text)

	d := gomail.NewDialer("smtp.gmail.com", 587, mail_user, mail_pass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
