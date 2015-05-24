package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tdgs/airline_fair_tracker/tripinfo"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	//response, e := tripinfo.ReadFromFile("response.json")
	//check(e)

	dateString := flag.String("date", "", "date of travel")
	flag.Parse()
	if *dateString == "" {
		log.Fatal("You must provide a date string in the form of 2015-06-19")
	}

	travelDate, e := time.Parse("2006-01-02", *dateString)
	check(e)
	response, e := tripinfo.GetDataAndWriteToFile(travelDate)
	check(e)

	tripInfos := []*tripinfo.TripInfo{}

	for _, tripOption := range response.Trips.TripOption {
		info, e := tripinfo.MakeTripInfo(tripOption)
		check(e)

		fmt.Println(info)
		tripInfos = append(tripInfos, info)
	}
}
