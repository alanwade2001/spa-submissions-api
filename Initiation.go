package main

import types "github.com/alanwade2001/spa-common"

// InitiationService s
type InitiationService struct {
	xmlParserAPI        XMLParserAPI
	initiationMapperAPI InitiationMapperAPI
}

// NewInitiationService f
func NewInitiationService(xmlParserAPI XMLParserAPI, initiationMapperAPI InitiationMapperAPI) InitiationAPI {
	return InitiationService{
		xmlParserAPI:        xmlParserAPI,
		initiationMapperAPI: initiationMapperAPI,
	}
}

// Parse f
func (p InitiationService) Parse(data []byte) (initiation *types.Initiation, err error) {
	if doc, err := p.xmlParserAPI.Parse(data); err != nil {
		return nil, err
	} else if initiation, err = p.initiationMapperAPI.Map(doc); err != nil {
		return nil, err
	}

	return initiation, err
}
