package Book

import "encoding/xml"

type FB2 struct {
	XMLName     xml.Name    `xml:"FictionBook"`
	Description Description `xml:"description"`
	Body        Body        `xml:"body"`
}

type Description struct {
	XMLName      xml.Name     `xml:"description"`
	TitleInfo    TitleInfo    `xml:"title-info"`
	DocumentInfo DocumentInfo `xml:"document-info"`
}

type TitleInfo struct {
	XMLName   xml.Name  `xml:"title-info"`
	Genre     Genre     `xml:"genre"`
	Author    Author    `xml:"author"`
	BookTitle BookTitle `xml:"book-title"`
}

type Genre struct {
	XMLName xml.Name `xml:"genre"`
	Text    string   `xml:",chardata"`
}

type Author struct {
	XMLName   xml.Name  `xml:"author"`
	FirstName FirstName `xml:"first-name"`
	LastName  LastName  `xml:"last-name"`
}

type FirstName struct {
	XMLName xml.Name `xml:"first-name"`
	Text    string   `xml:",chardata"`
}

type LastName struct {
	XMLName xml.Name `xml:"last-name"`
	Text    string   `xml:",chardata"`
}

type BookTitle struct {
	XMLName xml.Name `xml:"book-title"`
	Text    string   `xml:",chardata"`
}

type DocumentInfo struct {
	XMLName xml.Name `xml:"document-info"`
	Author  Author   `xml:"author"`
	Date    Date     `xml:"date"`
}

type Date struct {
	XMLName xml.Name `xml:"date"`
	Text    string   `xml:",chardata"`
}

type Body struct {
	XMLName xml.Name `xml:"body"`
	Section Section  `xml:"section"`
}

type Section struct {
	XMLName   xml.Name  `xml:"section"`
	Title     Title     `xml:"title"`
	Paragraph Paragraph `xml:"p"`
}

type Title struct {
	XMLName xml.Name `xml:"title"`
	Text    string   `xml:",chardata"`
}

type Paragraph struct {
	XMLName xml.Name `xml:"p"`
	Text    []string `xml:",p"`
}
