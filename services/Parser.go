package services

import (
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/lestrrat-go/libxml2"
	xml "github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xsd"
)

type lestrratXMLParser struct {
	schema *xsd.Schema
}

// NewXMLParserAPI f
func NewXMLParserAPI() types.XMLParserAPI {
	parser := new(lestrratXMLParser)
	var err error

	if parser.schema, err = xsd.ParseFromFile("schemas/pain.001.001.03.xsd"); err != nil {
		panic(err)
	}

	return parser
}

func (p *lestrratXMLParser) Parse(data []byte) (doc xml.Document, err error) {
	if doc, err = libxml2.Parse(data); err != nil {
		return nil, err

	} else if err = p.schema.Validate(doc); err != nil {
		return nil, err
	}

	return doc, nil
}
