package scraper

import (
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an HTML link with its text and href
type FishInfo struct {
	Name       string
	Shape      shape      `json:"shape"`
	Estimation estimation `json:"estimation"`
}

type shape struct {
	MaxLength     string `json:"maxLength"`     // e.g. "54.0"
	MaxLengthUnit string `json:"maxLengthUnit"` // e.g. "cm"
}

type estimation struct {
	BayesianA    string `json:"bayesianA"`    // e.g. "0.03020"
	BayesianAMin string `json:"bayesianAMin"` // e.g. "0.01889"
	BayesianAMax string `json:"bayesianAMax"` // e.g. "0.04829"
	BayesianB    string `json:"bayesianB"`    // e.g. "2.94"
	BayesianBMin string `json:"bayesianBMin"` // e.g. "2.81"
	BayesianBMax string `json:"bayesianBMax"` // e.g. "3.07"
}

// Parse extracts all links from HTML content
func Parse(page io.Reader) (FishInfo, error) {
	doc, err := html.Parse(page)
	if err != nil {
		return FishInfo{}, err
	}

	var fish FishInfo
	var parseErr error

	var f func(*html.Node)
	f = func(n *html.Node) {
		if parseErr != nil {
			return
		}

		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "ss-sciname" {
					fish.Name = extractName(n)
					return
				}
			}
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
		return FishInfo{}, parseErr
	}

	return fish, nil
}

func collectAllContent(n *html.Node) (rawContent string) {
	var allText strings.Builder

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
	return content
}

func extractName(n *html.Node) (name string) {
	var sciNames []string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "sciname" {
					sciNames = append(sciNames, collectAllContent(n))
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)

	return strings.Join(sciNames, " ")
}

func extractShape(n *html.Node) (shape, error) {
	var s shape

	content := collectAllContent(n)

	// Parse Max Length with Unit (any unit)
	if maxLenMatch := regexp.MustCompile(`Max length\s*:\s*([\d\.]+)\s*([^;]+)`).FindStringSubmatch(content); len(maxLenMatch) > 2 {
		s.MaxLength = strings.TrimSpace(maxLenMatch[1])
		s.MaxLengthUnit = strings.TrimSpace(maxLenMatch[2])
	}

	return s, nil
}

func extractEstimation(n *html.Node) (estimation, error) {
	var e estimation

	content := collectAllContent(n)

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

	return e, nil
}
