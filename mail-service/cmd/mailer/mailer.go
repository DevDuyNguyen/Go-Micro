package mailer

import (
	"github.com/wneessen/go-mail"
	gomail "github.com/wneessen/go-mail"
)

// type mailerConfig struct {
// 	SMTPAuthType gomail.SMTPAuthType
// 	TLSPolicy gomail.TLSPolicy
// 	UserName string
// 	Password string
// }
type Mailer struct{
	mailUserAgent *gomail.Client
}

type SMTPMessage struct{
	MailFrom string
	MailTo string
	Headers map[string]string
	MessageBody string
}

func New(host string, port int, hasSMTPAuth bool, hasTLSPolicy bool, UserName string, Password string) (*Mailer,error){
	var smtpAuthType gomail.SMTPAuthType
	var tlsPolicy gomail.TLSPolicy
	if hasSMTPAuth==true{
		smtpAuthType=gomail.SMTPAuthAutoDiscover
	}else{
		smtpAuthType=gomail.SMTPAuthNoAuth
	}
	if hasTLSPolicy==true{
		tlsPolicy=gomail.TLSOpportunistic
	}else{
		tlsPolicy=gomail.NoTLS
	}

	client,err:=gomail.NewClient(host, mail.WithPort(port),
		mail.WithSMTPAuth(smtpAuthType), mail.WithTLSPortPolicy(tlsPolicy),
		mail.WithUsername(UserName), mail.WithPassword(Password),
	)
	if err!=nil{
		return nil, err
	}

	return &Mailer{
		mailUserAgent: client,
	}, nil
}

func (mailer *Mailer) SendSMTPMAil(mess SMTPMessage) error{
	smtpMessage:=gomail.NewMsg()
	if err:= smtpMessage.From(mess.MailFrom); err!=nil{
		return err
	}	
	if err:=smtpMessage.To(mess.MailFrom); err!=nil{
		return err
	}
	smtpMessage.Subject(mess.Headers["Subject"])
	smtpMessage.SetBodyString(gomail.TypeTextPlain, mess.MessageBody)

	err:=mailer.mailUserAgent.DialAndSend(smtpMessage)
	if err!=nil{
		return err
	}
	return nil
}