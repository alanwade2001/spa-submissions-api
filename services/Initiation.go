package services

import (
	"github.com/alanwade2001/spa-submissions-api/models/generated/initiation"
	"github.com/alanwade2001/spa-submissions-api/types"
)

// InitiationService s
type InitiationService struct {
	xmlParserAPI        types.XMLParserAPI
	initiationMapperAPI types.InitiationMapperAPI
}

// NewInitiationService f
func NewInitiationService(xmlParserAPI types.XMLParserAPI, initiationMapperAPI types.InitiationMapperAPI) types.InitiationAPI {
	return InitiationService{
		xmlParserAPI:        xmlParserAPI,
		initiationMapperAPI: initiationMapperAPI,
	}
}

// Parse f
func (p InitiationService) Parse(data []byte) (initiation *initiation.InitiationModel, err error) {
	if doc, err := p.xmlParserAPI.Parse(data); err != nil {
		return nil, err
	} else if initiation, err = p.initiationMapperAPI.Map(doc); err != nil {
		return nil, err
	}

	return initiation, err
}
