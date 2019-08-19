package notification

import (
	"fmt"
	"log"

	"github.com/bshafiee/bonga/scraping"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const sendmail = "/usr/sbin/sendmail"

type email struct {
}

func NewEmailNotification() Channel {
	return &email{}
}

func (e *email) Notify(res []scraping.Result) error {
	body := "Hello <b>Behrooz & Hannah</b>😍<br> "
	for _, r := range res {
		body += r.Price + " <a href='" + r.URL + "'>" + r.Title + "</a><br>"
	}
	return submitMail(body)
}

func submitMail(body string) (err error) {
	m := mail.NewV3Mail()
	address := "notify@listings.com"
	name := "Behrooz"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)
	m.Subject = "Found these new listings"
	p := mail.NewPersonalization()
	p.AddTos(mail.NewEmail("Behrooz", "shafiee01@gmail.com"), mail.NewEmail("Hannah", "champ.hannah@gmail.com"))
	m.AddPersonalizations(p)
	m.AddContent(mail.NewContent("text/html", body))
	client := sendgrid.NewSendClient("")
	if _, err := client.Send(m); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("email successfully sent")
	return nil

}
