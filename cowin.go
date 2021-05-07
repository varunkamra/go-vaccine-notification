package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gen2brain/beeep"
)

func main() {
	date := time.Now()
	const SOUTH_DELHI_DISTRICT_ID = 150
	url := fmt.Sprintf("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?district_id=%d&date=%s", SOUTH_DELHI_DISTRICT_ID, date.Format("02-01-2006"))
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	type Session struct {
		Session_id         string
		Date               string
		Available_capacity int
		Min_age_limit      int
		Vaccine            string
		Slots              []string
	}
	type Center struct {
		Center_id     int
		Name          string
		Address       string
		State_name    string
		District_name string
		Block_name    string
		Pincode       int
		Lat           int
		Long          int
		From          string
		To            string
		Fee_type      string
		Sessions      []Session
	}
	type Centers struct {
		Centers []Center
	}
	center1 := Centers{}

	err = json.NewDecoder(resp.Body).Decode(&center1)
	if err != nil {
		panic(err)
	}
	found := false
	for _, cent := range center1.Centers {
		for _, sess := range cent.Sessions {
			if sess.Min_age_limit == 18 && sess.Available_capacity > 0 {
				found = true
				title := "Appointment found!"
				body := cent.Name + "\n" + cent.Address
				notify(title, body)
			}
		}
	}
	if !found {
		title := "Appointment not found!"
		body := "No appointments found."
		notify(title, body)
	}
}

func notify(title string, body string) {
	err := beeep.Notify(title, body, "")
	if err != nil {
		panic(err)
	}
}
