package tripinfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/qpxexpress/v1"
)

type TripInfo struct {
	Currency  string
	Price     float32
	TripID    string
	Duration  int64
	TripDate  string
	QueryTime time.Time
}

func apiClient() (*qpxexpress.Service, error) {
	client := &http.Client{Transport: &transport.APIKey{Key: os.Getenv("QPX_EXPRESS_KEY")}}
	return qpxexpress.New(client)
}

func searchRequest(dateString string) *qpxexpress.TripsSearchRequest {
	// Trave information
	itenary := []*qpxexpress.SliceInput{&qpxexpress.SliceInput{Date: dateString, Destination: "ATH", Origin: "HAM", MaxStops: 2}}
	passengerCounts := &qpxexpress.PassengerCounts{AdultCount: 1}
	tripOptionsRequest := &qpxexpress.TripOptionsRequest{Passengers: passengerCounts, SaleCountry: "GR", Slice: itenary, MaxPrice: "EUR400"}
	tripsSearchRequest := &qpxexpress.TripsSearchRequest{Request: tripOptionsRequest}

	return tripsSearchRequest
}

func GetDataFromApi(dateString string) (*qpxexpress.TripsSearchResponse, error) {
	client, e := apiClient()
	if e != nil {
		return nil, e
	}
	searchRequest := searchRequest(dateString)
	return client.Trips.Search(searchRequest).Do()
}

func WriteToFile(response *qpxexpress.TripsSearchResponse, filename string) error {
	jsonResponse, e := json.Marshal(response)
	if e != nil {
		return e
	}
	return ioutil.WriteFile(filename, jsonResponse, 0644)
}

func ReadFromFile(filename string) (*qpxexpress.TripsSearchResponse, error) {
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}

	var response qpxexpress.TripsSearchResponse
	json.Unmarshal(data, &response)

	return &response, nil
}

func GetDataAndWriteToFile(date time.Time) (*qpxexpress.TripsSearchResponse, error) {
	dateString := date.Format("2006-01-02")
	response, e := GetDataFromApi(dateString)
	if e != nil {
		return nil, e
	}

	filename := fmt.Sprintf("results/response_%v_%v.json", dateString, time.Now().Unix())
	return response, WriteToFile(response, filename)
}

func MakeTripInfo(t *qpxexpress.TripOption) (*TripInfo, error) {
	var tripInfo TripInfo

	tripInfo.Currency = t.SaleTotal[0:3]
	price, e := strconv.ParseFloat(t.SaleTotal[3:], 32)
	if e != nil {
		return nil, e
	}

	tripInfo.Price = float32(price)
	tripInfo.QueryTime = time.Now()

	slice := t.Slice[0]
	tripInfo.Duration = slice.Duration

	tripIDs := []string{}
	tripInfo.TripDate = slice.Segment[0].Leg[0].DepartureTime

	for _, segment := range slice.Segment {
		tripIDs = append(tripIDs, segment.Flight.Carrier+segment.Flight.Number)
	}

	tripInfo.TripID = strings.Join(tripIDs, ";")
	return &tripInfo, nil
}
