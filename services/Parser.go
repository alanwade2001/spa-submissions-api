package services

import (
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/lestrrat-go/libxml2"
	xml "github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/spf13/viper"
)

type lestrratXMLParser struct {
	schema *xsd.Schema
}

// NewXMLParserAPI f
func NewXMLParserAPI() types.XMLParserAPI {
	parser := new(lestrratXMLParser)

	return parser
}

func (p *lestrratXMLParser) GetSchema() (*xsd.Schema, error) {
	if p.schema != nil {
		return p.schema, nil
	}

	var err error

	xsdPath := viper.GetString("PAIN_001_XSD_FILE")
	if p.schema, err = xsd.ParseFromFile(xsdPath); err != nil {
		return nil, err
	}

	return p.schema, nil
}

func (p *lestrratXMLParser) Parse(data []byte) (doc xml.Document, err error) {

	if schema, err := p.GetSchema(); err != nil {
		return nil, err
	} else if doc, err = libxml2.Parse(data); err != nil {
		return nil, err

	} else if err = schema.Validate(doc); err != nil {
		return nil, err
	}

	return doc, nil
}
