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

func getSections(e *colly.HTMLElement) []Section {
	var sections []Section
	e.ForEach("h2", func(_ int, e *colly.HTMLElement) {
		section := Section{
			Title:      e.Text,
			Paragraphs: getParagraphs(e),
		}
		sections = append(sections, section)
	})
	return sections
}

func getReferences(e *colly.HTMLElement) []Reference {
	var references []Reference
	e.ForEach(".reflist li", func(_ int, e *colly.HTMLElement) {
		reference := Reference{
			Id:    e.Attr("id"),
			Text:  e.ChildText(".cite"),
			Links: getLinks(e),
		}
		references = append(references, reference)
	})
	return references
}

func getArticle(c *colly.Collector, url string) Article {
	var article Article
	c.OnHTML(".mw-content-container", func(e *colly.HTMLElement) {
		introduction := Section{
			Title:      "Introduction",
			Paragraphs: getParagraphs(e),
		}
		article.Title = e.ChildText("#firstHeading")
		article.Sections = append(article.Sections, introduction)
		article.Sections = append(article.Sections, getSections(e)...)
		article.References = append(article.References, getReferences(e)...)
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
