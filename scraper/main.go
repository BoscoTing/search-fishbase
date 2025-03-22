package scraper

import (
	"log"
)

func Scrape(name string) FishInfo {
	log.Printf("Scraping %s", name)

	html, err := Download(name)
	if err != nil {
		log.Printf("error downloading %s: %v", name, err)
		return FishInfo{Name: name}
	}
	defer html.Close()

	parsed, err := Parse(html)
	if err != nil {
		log.Printf("error parsing %s: %v", name, err)
		return FishInfo{Name: name}
	}
	fish := parsed

	log.Printf(
		"Fish data: %s, Shape: %s %s, Estimation: A=%s, B=%s",
		fish.Name,
		fish.Shape.MaxLength,
		fish.Shape.MaxLengthUnit,
		fish.Estimation.BayesianA,
		fish.Estimation.BayesianB,
	)
	return fish
}
