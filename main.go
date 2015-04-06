package main

import (
	"fmt"
	"time"

	"github.com/tdgs/airline_fair_tracker/tripinfo"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//response, e := tripinfo.ReadFromFile("response.json")
	//check(e)
	travelDate, e := time.Parse("2006-01-02", "2015-06-19")
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
