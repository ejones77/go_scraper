package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestGetLinks(t *testing.T) {
	html := `<p>This is a <a href="https://example.com">link</a>.</p>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	selection := doc.Selection
	links := getLinks(selection)
	assert.Equal(t, 1, len(links), "Expected 1 link")
}

func TestGetParagraphs(t *testing.T) {
	html := `<body><p>This is a paragraph with a <a href="https://example.com">link</a>.</p>
	<ul><li>This is a list item with a <a href="https://example.com">link</a>.</li></ul></body>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	selection := doc.Find("p, ul")
	paragraphs := getParagraphs(selection)
	fmt.Println(paragraphs)
	assert.Equal(t, 2, len(paragraphs), "Expected 2 paragraphs")
	assert.Equal(t, "This is a paragraph with a link.", paragraphs[0].Text, "Expected first paragraph text to be 'This is a paragraph with a link.'")
	assert.Equal(t, "This is a list item with a link.", paragraphs[1].Text, "Expected second paragraph text to be 'This is a list item with a link.'")
}
