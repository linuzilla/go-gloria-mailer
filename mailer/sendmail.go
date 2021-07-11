package mailer

import (
	"fmt"
	"github.com/linuzilla/go-gloria-mailer/config"
	"github.com/linuzilla/go-gloria-mailer/flags"
	"github.com/linuzilla/go-gloria-mailer/mime-composer"
	"net/smtp"
)

type Client interface {
	SendMailTo(recipient string, displayName string, subject string,
		callback func(composer mime_composer.MimeComposer)) error
	SetFrom(email, displayName string)
}

type clientImpl struct {
	settings            config.SmtpSection
	smtpHost            string
	mailFrom            string
	mailFromDisplayName string
	auth                smtp.Auth
}

func New(settings config.SmtpSection) Client {
	var auth smtp.Auth = nil

	if settings.Auth {
		auth = smtp.PlainAuth("", settings.User, settings.Password, settings.Host)
	}
	return &clientImpl{
		smtpHost: fmt.Sprintf("%s:%d", settings.Host, settings.Port),
		auth:     auth,
		settings: settings,
	}
}

func (impl *clientImpl) SetFrom(email, displayName string) {
	impl.mailFrom = email
	impl.mailFromDisplayName = displayName
}

func (impl *clientImpl) SendMailTo(recipient string, displayName string, subject string,
	callback func(composer mime_composer.MimeComposer)) error {
	fmt.Printf("Sendmail to : [ %s ] (%s)\n", recipient, displayName)

	composer := mime_composer.New().
		From(impl.mailFrom, impl.mailFromDisplayName).
		To(recipient, displayName).
		Subject(subject)
	callback(composer)

	if flags.SendEmail {
		if err := smtp.SendMail(impl.smtpHost, impl.auth, impl.mailFrom, []string{recipient}, []byte(composer.Compose())); err != nil {
			return err
		} else {
			fmt.Println("Email sent to:", recipient, " via ", impl.smtpHost)
		}
	} else {
		fmt.Println("Just testing, email not really sent to", recipient)
	}
	return nil
}
