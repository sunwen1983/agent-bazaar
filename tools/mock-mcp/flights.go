package main

import "sort"

// flights: 6 round-trip SFO -> AUS options, priced per person. The
// cheapest round trip for 2 adults is $356, which leaves room in the
// $800 demo budget for a hotel and activities. Times are realistic
// Friday-out / Sunday-back timings; booking URLs are the airlines'
// home pages (no invented deep links).

type flight struct {
	Airline         string  `json:"airline"`
	FlightNumbers   string  `json:"flight_numbers"`
	Stops           int     `json:"stops"`
	StopCity        string  `json:"stop_city,omitempty"`
	Depart          string  `json:"depart"`
	Arrive          string  `json:"arrive"`
	ReturnDepart    string  `json:"return_depart"`
	ReturnArrive    string  `json:"return_arrive"`
	Duration        string  `json:"duration_each_way"`
	PricePerPerson  float64 `json:"price_per_person_usd"`
	TotalPrice      float64 `json:"total_price_usd"`
	BookingURL      string  `json:"booking_url"`
}

var flightInventory = []flight{
	{Airline: "Southwest", FlightNumbers: "WN 1842 / WN 1977", Stops: 0, Depart: "Fri 07:10 SFO", Arrive: "Fri 12:35 AUS", ReturnDepart: "Sun 18:20 AUS", ReturnArrive: "Sun 20:05 SFO", Duration: "3h 25m", PricePerPerson: 178, BookingURL: "https://www.southwest.com"},
	{Airline: "Alaska", FlightNumbers: "AS 1120 / AS 1121", Stops: 0, Depart: "Fri 08:05 SFO", Arrive: "Fri 13:36 AUS", ReturnDepart: "Sun 16:40 AUS", ReturnArrive: "Sun 18:55 SFO", Duration: "3h 31m", PricePerPerson: 196, BookingURL: "https://www.alaskaair.com"},
	{Airline: "United", FlightNumbers: "UA 2451 / UA 2388", Stops: 1, StopCity: "DEN", Depart: "Fri 06:30 SFO", Arrive: "Fri 12:48 AUS", ReturnDepart: "Sun 14:15 AUS", ReturnArrive: "Sun 19:02 SFO", Duration: "4h 18m", PricePerPerson: 214, BookingURL: "https://www.united.com"},
	{Airline: "American", FlightNumbers: "AA 3098 / AA 3102", Stops: 1, StopCity: "DFW", Depart: "Fri 09:20 SFO", Arrive: "Fri 15:05 AUS", ReturnDepart: "Sun 13:10 AUS", ReturnArrive: "Sun 18:40 SFO", Duration: "3h 45m", PricePerPerson: 229, BookingURL: "https://www.aa.com"},
	{Airline: "Delta", FlightNumbers: "DL 1440 / DL 1441", Stops: 0, Depart: "Fri 10:15 SFO", Arrive: "Fri 15:43 AUS", ReturnDepart: "Sun 17:05 AUS", ReturnArrive: "Sun 19:30 SFO", Duration: "3h 28m", PricePerPerson: 247, BookingURL: "https://www.delta.com"},
	{Airline: "JetBlue", FlightNumbers: "B6 1204 / B6 1205", Stops: 1, StopCity: "LAX", Depart: "Fri 07:45 SFO", Arrive: "Fri 14:20 AUS", ReturnDepart: "Sun 12:30 AUS", ReturnArrive: "Sun 18:15 SFO", Duration: "4h 35m", PricePerPerson: 268, BookingURL: "https://www.jetblue.com"},
}

func searchFlights(args map[string]any) (any, *rpcErr) {
	origin := strArg(args, "origin", "SFO")
	destination := strArg(args, "destination", "AUS")
	adults := intArg(args, "adults", 2)
	if adults < 1 {
		adults = 1
	}
	maxPrice := floatArg(args, "max_price_per_person", 0)

	options := make([]flight, 0, len(flightInventory))
	for _, f := range flightInventory {
		if maxPrice > 0 && f.PricePerPerson > maxPrice {
			continue
		}
		f.TotalPrice = f.PricePerPerson * float64(adults)
		options = append(options, f)
	}
	sort.Slice(options, func(i, j int) bool { return options[i].TotalPrice < options[j].TotalPrice })

	result := map[string]any{
		"origin":      origin,
		"destination": destination,
		"adults":      adults,
		"currency":    "USD",
		"options":     options,
		"note":        "Mock data — deterministic sample inventory for the Agent Bazaar demo, not live availability. Prices are round-trip per person; total is for the whole party.",
	}
	if d := strArg(args, "depart_date", ""); d != "" {
		result["depart_date"] = d
	}
	if d := strArg(args, "return_date", ""); d != "" {
		result["return_date"] = d
	}
	return result, nil
}

func flightsDataset() dataset {
	return dataset{
		Name: "flights",
		Tools: []tool{{
			Name:        "search_flights",
			Description: "Search round-trip flights to Austin (AUS). Returns comparable options sorted by total price, cheapest first. Mock data: fixed sample inventory, not live availability.",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"origin":               map[string]any{"type": "string", "description": "Departure airport code, e.g. SFO"},
					"destination":          map[string]any{"type": "string", "description": "Arrival airport code, e.g. AUS"},
					"depart_date":          map[string]any{"type": "string", "description": "Outbound date, e.g. 2026-12-04"},
					"return_date":          map[string]any{"type": "string", "description": "Return date, e.g. 2026-12-06"},
					"adults":               map[string]any{"type": "number", "description": "Party size (default 2)"},
					"max_price_per_person": map[string]any{"type": "number", "description": "Optional ceiling on round-trip price per person in USD"},
				},
			},
			Call: searchFlights,
		}},
	}
}
