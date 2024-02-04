package main

import (
	"fmt"

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

type Article struct {
	Title      string      `json:"text"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

/*
type Section struct {
	Title string `json:"title"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

type Article struct {
	Title        string      `json:"title"`
	Introduction []Paragraph   `json:"introduction"`
	Sections     []Section  `json:"sections"`
}
*/

func getLinks(e *colly.HTMLElement) []Link {
	var links []Link
	e.ForEach("a", func(_ int, e *colly.HTMLElement) {
		link := Link{
			Text: e.Text,
			URL:  e.Attr("href"),
		}
		links = append(links, link)
	})
	return links
}

func getParagraphs(e *colly.HTMLElement) []Paragraph {
	var paragraphs []Paragraph
	e.ForEach("p", func(_ int, e *colly.HTMLElement) {
		paragraph := Paragraph{
			Text:  e.Text,
			Links: getLinks(e),
		}
		paragraphs = append(paragraphs, paragraph)
	})
	return paragraphs
}

func getArticle(c *colly.Collector, url string) Article {
	var article Article
	c.OnHTML(".mw-content-container", func(e *colly.HTMLElement) {
		article = Article{
			Title:      e.ChildText("#firstHeading"),
			Paragraphs: getParagraphs(e),
		}
	})
	c.Visit(url)
	return article
}

func main() {
	url := "https://en.wikipedia.org/wiki/Robotics"
	c := colly.NewCollector()
	article := getArticle(c, url)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	fmt.Printf("%+v\n", article)
}
