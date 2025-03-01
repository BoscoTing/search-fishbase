package scraper

import (
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an HTML link with its text and href
type fishInfo struct {
	Shape      shape      `json:"shape"`
	Estimation estimation `json:"estimation"`
}

type shape struct {
	RawContent       string `json:"rawContent"`       // Store the complete raw content
	Maturity         string `json:"maturity"`         // e.g. "L"
	LengthRange      string `json:"lengthRange"`      // e.g. "? - ? cm"
	MaxLength        string `json:"maxLength"`        // e.g. "54.0"
	MaxLengthUnit    string `json:"maxLengthUnit"`    // e.g. "cm TL"
	CommonLength     string `json:"commonLength"`     // e.g. "35.0"
	CommonLengthUnit string `json:"commonLengthUnit"` // e.g. "cm TL"
	MaxAge           string `json:"maxAge"`           // e.g. "30"
	MaxAgeUnit       string `json:"maxAgeUnit"`       // e.g. "years"
}

type estimation struct {
	BayesianA    string `json:"bayesianA"`    // e.g. "0.03020"
	BayesianAMin string `json:"bayesianAMin"` // e.g. "0.01889"
	BayesianAMax string `json:"bayesianAMax"` // e.g. "0.04829"
	BayesianB    string `json:"bayesianB"`    // e.g. "2.94"
	BayesianBMin string `json:"bayesianBMin"` // e.g. "2.81"
	BayesianBMax string `json:"bayesianBMax"` // e.g. "3.07"
	ReferenceID  string `json:"referenceId"`  // e.g. "93245"
}

// Parse extracts all links from HTML content
func Parse(page io.Reader) (fishInfo, error) {
	doc, err := html.Parse(page)
	if err != nil {
		return fishInfo{}, err
	}

	var fish fishInfo
	var parseErr error

	var f func(*html.Node)
	f = func(n *html.Node) {
		if parseErr != nil {
			return
		}

		if n.Type == html.ElementNode && n.Data == "h1" {
			text := ""
			if n.FirstChild != nil {
				text = strings.TrimSpace(n.FirstChild.Data)
			}

			if text == "Size / Weight / Age" {
				nextNode := n.NextSibling
				for nextNode != nil && nextNode.Type != html.ElementNode {
					nextNode = nextNode.NextSibling
				}
				if nextNode != nil && nextNode.Data == "div" {
					shape, err := extractShape(nextNode)
					if err != nil {
						parseErr = err
						return
					}
					fish.Shape = shape
				}
			} else if text == "Estimates based on models" {
				nextNode := n.NextSibling
				for nextNode != nil && nextNode.Type != html.ElementNode {
					nextNode = nextNode.NextSibling
				}
				if nextNode != nil && nextNode.Data == "div" {
					estimation, err := extractEstimation(nextNode)
					if err != nil {
						parseErr = err
						return
					}
					fish.Estimation = estimation
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if parseErr != nil {
		return fishInfo{}, parseErr
	}

	return fish, nil
}

func extractShape(n *html.Node) (shape, error) {
	var s shape
	var allText strings.Builder

	// First collect all text content
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				allText.WriteString(text)
				allText.WriteString(" ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	content := strings.TrimSpace(allText.String())
	s.RawContent = content

	// Parse Maturity (L/M/...)
	if maturityMatch := regexp.MustCompile(`Maturity:\s*([LM])`).FindStringSubmatch(content); len(maturityMatch) > 1 {
		s.Maturity = maturityMatch[1]
	}

	// Parse Length Range (any unit)
	if rangeMatch := regexp.MustCompile(`range\s*(\?\s*-\s*\?[^;]*)`).FindStringSubmatch(content); len(rangeMatch) > 1 {
		s.LengthRange = strings.TrimSpace(rangeMatch[1])
	}

	// Parse Max Length with Unit (any unit)
	if maxLenMatch := regexp.MustCompile(`Max length\s*:\s*([\d\.]+)\s*([^;]+)`).FindStringSubmatch(content); len(maxLenMatch) > 2 {
		s.MaxLength = strings.TrimSpace(maxLenMatch[1])
		s.MaxLengthUnit = strings.TrimSpace(maxLenMatch[2])
	}

	// Parse Common Length with Unit (any unit)
	if commonLenMatch := regexp.MustCompile(`common length\s*:\s*([\d\.]+)\s*([^;]+)`).FindStringSubmatch(content); len(commonLenMatch) > 2 {
		s.CommonLength = strings.TrimSpace(commonLenMatch[1])
		s.CommonLengthUnit = strings.TrimSpace(commonLenMatch[2])
	}

	// Parse Max Age
	if maxAgeMatch := regexp.MustCompile(`max\.\s*reported age:\s*(\d+)\s*([^\s(]+)`).FindStringSubmatch(content); len(maxAgeMatch) > 2 {
		s.MaxAge = strings.TrimSpace(maxAgeMatch[1])
		s.MaxAgeUnit = strings.TrimSpace(maxAgeMatch[2])
	}

	return s, nil
}

func extractEstimation(n *html.Node) (estimation, error) {
	var e estimation
	var allText strings.Builder

	// First collect all text content
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				allText.WriteString(text)
				allText.WriteString(" ")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	content := strings.TrimSpace(allText.String())

	// Parse Bayesian a parameter and its range
	if aMatch := regexp.MustCompile(`a=(\d+\.\d+)\s*\((\d+\.\d+)\s*-\s*(\d+\.\d+)\)`).FindStringSubmatch(content); len(aMatch) > 3 {
		e.BayesianA = strings.TrimSpace(aMatch[1])
		e.BayesianAMin = strings.TrimSpace(aMatch[2])
		e.BayesianAMax = strings.TrimSpace(aMatch[3])
	}

	// Parse Bayesian b parameter and its range
	if bMatch := regexp.MustCompile(`b=(\d+\.\d+)\s*\((\d+\.\d+)\s*-\s*(\d+\.\d+)\)`).FindStringSubmatch(content); len(bMatch) > 3 {
		e.BayesianB = strings.TrimSpace(bMatch[1])
		e.BayesianBMin = strings.TrimSpace(bMatch[2])
		e.BayesianBMax = strings.TrimSpace(bMatch[3])
	}

	// Parse Reference ID
	if refMatch := regexp.MustCompile(`Ref\.\s*(\d+)`).FindStringSubmatch(content); len(refMatch) > 1 {
		e.ReferenceID = strings.TrimSpace(refMatch[1])
	}

	return e, nil
}
