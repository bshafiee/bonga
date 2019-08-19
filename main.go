package main

import (
	"log"

	"github.com/bshafiee/bonga/craiglist"
	"github.com/bshafiee/bonga/notification"
	"github.com/bshafiee/bonga/scraping"
)

const sendmail = "/usr/sbin/sendmail"

func main() {
	notifEngine := notification.NewNotificationEngine("list.txt", []notification.Channel{notification.NewEmailNotification()})
	if err := notifEngine.Initialize(); err != nil {
		log.Fatal(err)
	}

	c := craiglist.NewCraiglistScraper()
	params := []scraping.Parameter{
		craiglist.HasPictureParam{true},
		craiglist.MinBathParam{2},
		craiglist.MinBedsParam{2},
		craiglist.MinPriceParam{2500},
		craiglist.MaxPriceParam{3200},
	}
	res, err := c.Query(params)
	if err != nil {
		log.Fatal(err)
	}

	if err := notifEngine.Notify(res); err != nil {
		log.Fatal(err)
	}
}
