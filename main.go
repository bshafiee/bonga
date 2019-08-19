package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bshafiee/bonga/craiglist"
	"github.com/bshafiee/bonga/notification"
	"github.com/bshafiee/bonga/scraping"
)

const sendmail = "/usr/sbin/sendmail"

func checkWritable() {
	file, err := os.Create("gs:///bonga-250321.appspot.com/hi.txt")
	if err != nil {
		log.Fatal("failed to open writable")
	}
	fmt.Println("file success")
	defer file.Close()
}

func main() {
	checkWritable()
	key := os.Getenv("SENDGRID_API_KEY")
	if len(key) == 0 {
		log.Fatal("could not load the SENDGRID key")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT must be set")
		port = "8080"
	}
	// implement /services/ping for service health
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Beroo")
	})
	go http.ListenAndServe(":"+port, nil)

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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ticker := time.NewTicker(300 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println("here tick")
			res, err := c.Query(params)
			if err != nil {
				log.Fatal(err)
			}
			if len(res) > 0 {
				if err := notifEngine.Notify(res); err != nil {
					log.Fatal(err)
				}
			}
		case <-sigCh:
			return
		}
	}

}
