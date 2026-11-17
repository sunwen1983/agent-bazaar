package main

import "sort"

// stays: 8 Austin hotels, priced per night. Two nights at the cheapest
// ($118) is $236; two nights at the mid-range Carpenter ($149) is $298.
// Combined with the cheapest flights ($356 for 2) and activities, the
// demo weekend stays under the $800 budget.

type stay struct {
	Name          string  `json:"name"`
	Neighborhood  string  `json:"neighborhood"`
	Rating        float64 `json:"rating_out_of_5"`
	PricePerNight float64 `json:"price_per_night_usd"`
	TotalPrice    float64 `json:"total_price_usd"`
	Cancellation  string  `json:"cancellation"`
	BookingURL    string  `json:"booking_url"`
}

var stayInventory = []stay{
	{Name: "Austin Motel", Neighborhood: "South Congress", Rating: 4.4, PricePerNight: 118, Cancellation: "Free cancellation until check-in", BookingURL: "https://austinmotel.com"},
	{Name: "Hotel San José", Neighborhood: "South Congress", Rating: 4.6, PricePerNight: 132, Cancellation: "Free cancellation until check-in", BookingURL: "https://sanjosehotel.com"},
	{Name: "East Austin Hotel", Neighborhood: "East Austin", Rating: 4.5, PricePerNight: 139, Cancellation: "Free cancellation until 48h before check-in", BookingURL: "https://www.eastaustinhotel.com"},
	{Name: "Carpenter Hotel", Neighborhood: "Zilker", Rating: 4.6, PricePerNight: 149, Cancellation: "Free cancellation until check-in", BookingURL: "https://www.carpenterhotel.com"},
	{Name: "The LINE Austin", Neighborhood: "Downtown", Rating: 4.5, PricePerNight: 164, Cancellation: "Free cancellation until 24h before check-in", BookingURL: "https://www.thelinehotel.com"},
	{Name: "Hotel Van Zandt", Neighborhood: "Rainey Street", Rating: 4.7, PricePerNight: 187, Cancellation: "Free cancellation until check-in", BookingURL: "https://www.hotelvanzandt.com"},
	{Name: "JW Marriott Austin", Neighborhood: "Downtown", Rating: 4.6, PricePerNight: 198, Cancellation: "Free cancellation until 48h before check-in", BookingURL: "https://www.marriott.com"},
	{Name: "Fairmont Austin", Neighborhood: "Downtown", Rating: 4.7, PricePerNight: 210, Cancellation: "Free cancellation until 72h before check-in", BookingURL: "https://www.fairmont.com"},
}

func searchStays(args map[string]any) (any, *rpcErr) {
	nights := intArg(args, "nights", 2)
	if nights < 1 {
		nights = 1
	}
	maxPrice := floatArg(args, "max_price_per_night", 0)
	neighborhood := strArg(args, "neighborhood", "")

	options := make([]stay, 0, len(stayInventory))
	for _, s := range stayInventory {
		if maxPrice > 0 && s.PricePerNight > maxPrice {
			continue
		}
		if neighborhood != "" && !containsFold(s.Neighborhood, neighborhood) {
			continue
		}
		s.TotalPrice = s.PricePerNight * float64(nights)
		options = append(options, s)
	}
	sort.Slice(options, func(i, j int) bool { return options[i].TotalPrice < options[j].TotalPrice })

	result := map[string]any{
		"city":     "Austin",
		"nights":   nights,
		"currency": "USD",
		"options":  options,
		"note":     "Mock data — deterministic sample inventory for the Agent Bazaar demo, not live availability. Total is for the whole stay (all nights, one room).",
	}
	if d := strArg(args, "check_in", ""); d != "" {
		result["check_in"] = d
	}
	if d := strArg(args, "check_out", ""); d != "" {
		result["check_out"] = d
	}
	return result, nil
}

func staysDataset() dataset {
	return dataset{
		Name: "stays",
		Tools: []tool{{
			Name:        "search_stays",
			Description: "Search hotels in Austin. Returns options sorted by total price, cheapest first. Mock data: fixed sample inventory, not live availability.",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"check_in":             map[string]any{"type": "string", "description": "Check-in date, e.g. 2026-12-04"},
					"check_out":            map[string]any{"type": "string", "description": "Check-out date, e.g. 2026-12-06"},
					"nights":               map[string]any{"type": "number", "description": "Number of nights (default 2)"},
					"adults":               map[string]any{"type": "number", "description": "Party size (default 2)"},
					"max_price_per_night":  map[string]any{"type": "number", "description": "Optional ceiling on nightly price in USD"},
					"neighborhood":         map[string]any{"type": "string", "description": "Optional neighborhood filter, e.g. Downtown"},
				},
			},
			Call: searchStays,
		}},
	}
}
