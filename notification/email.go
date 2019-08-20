package notification

import (
	"fmt"
	"log"
	"os"

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
		if len(r.Date) > 0 {
			body += r.Price + " (" + r.Date + ") <a href='" + r.URL + "'>" + getShortenTitle(r.Title) + "</a><br>"
		} else {
			body += r.Price + " <a href='" + r.URL + "'>" + getShortenTitle(r.Title) + "</a><br>"
		}
	}
	return submitMail(body)
}

func getShortenTitle(t string) string {
	if len(t) < 100 {
		return t
	}
	return t[0:100]
}

func submitMail(body string) (err error) {
	key := os.Getenv("SENDGRID_API_KEY")
	if len(key) == 0 {
		return fmt.Errorf("could not load the key")
	}
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
	client := sendgrid.NewSendClient(key)
	if _, err := client.Send(m); err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("email successfully sent")
	return nil

}
