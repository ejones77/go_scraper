package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Link struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

type Paragraph struct {
	Text  string `json:"text"`
	Links []Link `json:"links"`
}

type Section struct {
	Title      string      `json:"title"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

type Reference struct {
	Id    string `json:"id"`
	Text  string `json:"text"`
	Links []Link `json:"links"`
}

type Article struct {
	Title      string      `json:"title"`
	Sections   []Section   `json:"sections"`
	References []Reference `json:"references"`
}

func getLinks(s *goquery.Selection) []Link {
	var links []Link
	s.Find("a").Each(func(_ int, a *goquery.Selection) {
		link := Link{
			Text: a.Text(),
			URL:  a.AttrOr("href", ""),
		}
		links = append(links, link)
	})
	return links
}

func getParagraphs(s *goquery.Selection) []Paragraph {
	var paragraphs []Paragraph
	s.Each(func(_ int, n *goquery.Selection) {
		if n.Get(0).Data == "p" {
			paragraph := Paragraph{
				Text:  n.Text(),
				Links: getLinks(n),
			}
			paragraphs = append(paragraphs, paragraph)
		} else if n.Get(0).Data == "ul" || n.Get(0).Data == "ol" {
			n.Find("li").Each(func(_ int, li *goquery.Selection) {
				paragraph := Paragraph{
					Text:  li.Text(),
					Links: getLinks(li),
				}
				paragraphs = append(paragraphs, paragraph)
			})
		}
	})
	return paragraphs
}

func getSections(e *colly.HTMLElement) []Section {
	var sections []Section
	e.DOM.Find("h2").Each(func(_ int, s *goquery.Selection) {
		section := Section{
			Title:      s.Text(),
			Paragraphs: getParagraphs(s.NextUntil("h2")),
		}
		sections = append(sections, section)
	})
	return sections
}

func getReferences(e *colly.HTMLElement) []Reference {
	var references []Reference
	e.ForEach(".reflist li", func(_ int, e *colly.HTMLElement) {

		selection := e.DOM
		selection.Find("style").Remove()

		reference := Reference{
			Id:    e.Attr("id"),
			Text:  selection.Text(),
			Links: getLinks(e.DOM),
		}
		references = append(references, reference)
	})
	return references
}

func getArticle(c *colly.Collector, url string) Article {
	var article Article
	c.OnHTML(".mw-content-container", func(e *colly.HTMLElement) {
		// Assumes there's introductory text without a section
		introduction := Section{
			Title:      "Introduction",
			Paragraphs: getParagraphs(e.DOM.Find("p")),
		}
		article.Title = e.ChildText("#firstHeading")
		article.Sections = append(article.Sections, introduction)
		article.Sections = append(article.Sections, getSections(e)...)
		article.References = append(article.References, getReferences(e)...)
	})
	c.Visit(url)
	c.Wait()
	return article
}

func main() {
	urls := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
	}

	c := colly.NewCollector()

	file, err := os.Create("wikipedia-article-text.jl")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	for _, url := range urls {
		article := getArticle(c, url)
		data, _ := json.Marshal(article)
		_, err := file.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
