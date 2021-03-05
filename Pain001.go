package main

// Pain001Service s
type Pain001Service struct {
	xmlParserAPI     XMLParserAPI
	pain001MapperAPI Pain001MapperAPI
}

// NewPain001Service f
func NewPain001Service(xmlParserAPI XMLParserAPI, pain001MapperAPI Pain001MapperAPI) Pain001API {
	return Pain001Service{
		xmlParserAPI:     xmlParserAPI,
		pain001MapperAPI: pain001MapperAPI,
	}
}

// Parse f
func (p Pain001Service) Parse(data []byte) (pain001 *Pain001, err error) {
	if doc, err := p.xmlParserAPI.Parse(data); err != nil {
		return nil, err
	} else if pain001, err = p.pain001MapperAPI.Map(doc); err != nil {
		return nil, err
	}

	return pain001, err
}
