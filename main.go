package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

func main() {

	c := colly.NewCollector()
	url := "https://en.wikipedia.org/wiki/Robotics"

	c.OnHTML(".mw-parser-output p", func(e *colly.HTMLElement) {
		text := e.Text
		fmt.Println(text)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
