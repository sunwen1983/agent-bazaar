package main

import "sort"
import "strings"

// activities: 10 Austin things to do, priced per person. Three are free,
// so a full weekend of activities for 2 adults can cost anywhere from $0
// to about $180 — the demo budget has room either way.

type activity struct {
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Neighborhood    string  `json:"neighborhood"`
	PricePerPerson  float64 `json:"price_per_person_usd"`
	DurationHours   float64 `json:"duration_hours"`
	BestDay         string  `json:"best_day"`
	Description     string  `json:"description"`
}

var activityInventory = []activity{
	{Name: "Texas State Capitol tour", Category: "sightseeing", Neighborhood: "Downtown", PricePerPerson: 0, DurationHours: 2, BestDay: "Saturday", Description: "Free guided tour of the 1888 capitol building and grounds."},
	{Name: "Zilker Park & Barton Springs Pool", Category: "outdoor", Neighborhood: "Zilker", PricePerPerson: 9, DurationHours: 3, BestDay: "Saturday", Description: "Picnic on the lawn, then a swim in the 68°F spring-fed pool."},
	{Name: "Lady Bird Lake kayak rental", Category: "outdoor", Neighborhood: "Downtown", PricePerPerson: 22, DurationHours: 2, BestDay: "Saturday", Description: "Two-hour tandem kayak rental with skyline views; no experience needed."},
	{Name: "Franklin Barbecue lunch", Category: "food", Neighborhood: "East Austin", PricePerPerson: 28, DurationHours: 1.5, BestDay: "Saturday", Description: "Counter-service brisket and ribs; expect a line, it moves fast."},
	{Name: "Veracruz All Natural tacos", Category: "food", Neighborhood: "East Austin", PricePerPerson: 14, DurationHours: 1, BestDay: "Friday", Description: "Migas poblanas and al pastor tacos from the famous food truck."},
	{Name: "Live music at Continental Club", Category: "live music", Neighborhood: "South Congress", PricePerPerson: 25, DurationHours: 2.5, BestDay: "Friday", Description: "Austin's oldest live-music room; two sets a night, cash cover at the door."},
	{Name: "Blanton Museum of Art", Category: "culture", Neighborhood: "UT Campus", PricePerPerson: 15, DurationHours: 2, BestDay: "Sunday", Description: "Ellsworth Kelly's Austin chapel plus the Latin American collection."},
	{Name: "Mount Bonnell sunset hike", Category: "outdoor", Neighborhood: "Westlake", PricePerPerson: 0, DurationHours: 1.5, BestDay: "Sunday", Description: "102 steps up to the best sunset view over Lake Austin; free."},
	{Name: "South Congress Avenue stroll", Category: "walk", Neighborhood: "South Congress", PricePerPerson: 0, DurationHours: 2, BestDay: "Sunday", Description: "Cowboy boots, vintage shops, and food trucks along the avenue."},
	{Name: "Rainey Street evening bar walk", Category: "nightlife", Neighborhood: "Rainey Street", PricePerPerson: 30, DurationHours: 2, BestDay: "Saturday", Description: "Bungalow bars with food trucks; budget covers two drinks each."},
}

func searchActivities(args map[string]any) (any, *rpcErr) {
	interests := strSliceArg(args, "interests")
	maxPrice := floatArg(args, "max_price_per_person", 0)
	day := strArg(args, "day", "")

	options := make([]activity, 0, len(activityInventory))
	for _, a := range activityInventory {
		if maxPrice > 0 && a.PricePerPerson > maxPrice {
			continue
		}
		if day != "" && !strings.Contains(strings.ToLower(a.BestDay), strings.ToLower(day)) && !strings.Contains(strings.ToLower(day), strings.ToLower(a.BestDay)) {
			continue
		}
		if len(interests) > 0 {
			matched := false
			haystack := a.Category + " " + a.Name + " " + a.Description
			for _, in := range interests {
				if containsFold(haystack, in) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		options = append(options, a)
	}
	sort.Slice(options, func(i, j int) bool { return options[i].PricePerPerson < options[j].PricePerPerson })

	return map[string]any{
		"city":     "Austin",
		"currency": "USD",
		"options":  options,
		"note":     "Mock data — deterministic sample inventory for the Agent Bazaar demo, not live availability. Prices are per person.",
	}, nil
}

func activitiesDataset() dataset {
	return dataset{
		Name: "activities",
		Tools: []tool{{
			Name:        "search_activities",
			Description: "Find things to do in Austin: live music, food, outdoor activities, sightseeing, and neighborhood walks. Returns options sorted by price per person, cheapest first. Mock data: fixed sample inventory, not live availability.",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"interests":             map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "What the visitors like, e.g. [\"live music\", \"food\", \"outdoor\"]"},
					"max_price_per_person":  map[string]any{"type": "number", "description": "Optional ceiling on price per person in USD"},
					"day":                   map[string]any{"type": "string", "description": "Optional day filter, e.g. Saturday"},
				},
			},
			Call: searchActivities,
		}},
	}
}
