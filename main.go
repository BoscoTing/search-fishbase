package main

import (
	"fishbase/scraper"
	"fmt"
	"log"
)

func main() {
	html, err := scraper.Download("Bodianus-loxozonus")
	if err != nil {
		log.Fatal(err)
	}
	defer html.Close() // 呼叫方應該負責關閉 body

	parsed, err := scraper.Parse(html)
	if err != nil {
		log.Fatal(err)
	}
	fish := parsed
	fmt.Printf("Fish Info:\n")
	fmt.Printf("  Shape:\n")
	fmt.Printf("    Raw Content: %s\n", fish.Shape.RawContent)
	fmt.Printf("    Maturity: %s\n", fish.Shape.Maturity)
	fmt.Printf("    Length Range: %s\n", fish.Shape.LengthRange)
	fmt.Printf("    Max Length: %s %s\n", fish.Shape.MaxLength, fish.Shape.MaxLengthUnit)
	fmt.Printf("    Common Length: %s cm TL\n", fish.Shape.CommonLength)
	fmt.Printf("    Max Age: %s %s\n", fish.Shape.MaxAge, fish.Shape.MaxAgeUnit)
	fmt.Printf("  Estimation:\n")
	fmt.Printf("    Bayesian a: %s (%s - %s)\n", fish.Estimation.BayesianA, fish.Estimation.BayesianAMin, fish.Estimation.BayesianAMax)
	fmt.Printf("    Bayesian b: %s (%s - %s)\n", fish.Estimation.BayesianB, fish.Estimation.BayesianBMin, fish.Estimation.BayesianBMax)
	fmt.Printf("    Reference ID: %s\n", fish.Estimation.ReferenceID)
}
