package main

import (
	"fmt"

	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"
)

// Pain001Mapper s
type Pain001Mapper struct {
	groupHeaderMapperAPI        GroupHeaderMapperAPI
	paymentInformationMapperAPI PaymentInformationMapperAPI
}

// NewPain001Mapper f
func NewPain001Mapper(groupHeaderMapperAPI GroupHeaderMapperAPI, paymentInformationMapperAPI PaymentInformationMapperAPI) Pain001MapperAPI {
	return Pain001Mapper{
		groupHeaderMapperAPI:        groupHeaderMapperAPI,
		paymentInformationMapperAPI: paymentInformationMapperAPI,
	}
}

// Map f
func (m Pain001Mapper) Map(doc types.Document) (*Pain001, error) {
	//var root types.Node
	//var ctx *xpath.Context
	var grpHdr *GroupHeader
	var pmtInfs []PmtInf

	if root, err := doc.DocumentElement(); err != nil {
		return nil, err
	} else if ctx, err := xpath.NewContext(root); err != nil {
		return nil, err
	} else if err = ctx.RegisterNS("ns", "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"); err != nil {
		return nil, err
	} else if grpHdr, err = m.groupHeaderMapperAPI.Map(ctx); err != nil {
		return nil, err
	} else if pmtInfs, err = m.paymentInformationMapperAPI.Map(ctx); err != nil {
		return nil, err
	}

	pain001 := Pain001{
		GroupHeader: *grpHdr,
		PmtInfs:     pmtInfs,
	}

	return &pain001, nil
}

// GroupHeaderMapper s
type GroupHeaderMapper struct {
}

// NewGroupHeaderMapper f
func NewGroupHeaderMapper() GroupHeaderMapperAPI {
	return GroupHeaderMapper{}
}

// Map f
func (ghm GroupHeaderMapper) Map(ctx *xpath.Context) (gh *GroupHeader, err error) {
	var grpHdrIter types.NodeIter

	if grpHdrIter = xpath.NodeIter(ctx.Find("/ns:Document/ns:CstmrCdtTrfInitn/ns:GrpHdr")); grpHdrIter.Next() == false {
		return nil, fmt.Errorf("Unable to find group header but it has passed xsd check")
	}

	grpHdrNode := grpHdrIter.Node()

	var grpHdrCtx *xpath.Context
	if grpHdrCtx, err = xpath.NewContext(grpHdrNode); err != nil {
		return nil, err
	} else if err = grpHdrCtx.RegisterNS("ns", "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"); err != nil {
		return nil, err
	}

	gh = new(GroupHeader)
	gh.MsgID = xpath.String(grpHdrCtx.Find("ns:MsgId"))
	gh.CreDtTm = xpath.String(grpHdrCtx.Find("ns:CreDtTm"))
	gh.NbOfTxs = xpath.String(grpHdrCtx.Find("ns:NbOfTxs"))
	gh.CtrlSum = xpath.String(grpHdrCtx.Find("ns:CtrlSum"))
	gh.InitgPty = xpath.String(grpHdrCtx.Find("ns:InitgPty/ns:Id/ns:OrgId/ns:Othr/ns:Id"))

	return gh, nil
}

// PaymentInformationMapper s
type PaymentInformationMapper struct {
}

// NewPaymentInformationMapper f
func NewPaymentInformationMapper() PaymentInformationMapperAPI {
	return PaymentInformationMapper{}
}

// Map f
func (pim PaymentInformationMapper) Map(ctx *xpath.Context) (pis PmtInfs, err error) {
	pmtInfNodes := xpath.NodeList(ctx.Find("/ns:Document/ns:CstmrCdtTrfInitn/ns:PmtInf"))

	pmtInfs := make([]PmtInf, len(pmtInfNodes))

	for i, n := range pmtInfNodes {
		var pi *PmtInf

		if pi, err = pim.MapPmtInf(n); err != nil {
			return nil, err
		}

		pmtInfs[i] = *pi
	}

	return pmtInfs, nil
}

// MapPmtInf f
func (pim PaymentInformationMapper) MapPmtInf(pmtInfNode types.Node) (pi *PmtInf, err error) {

	pmtInfCtx, err := xpath.NewContext(pmtInfNode)
	if err != nil {
		return nil, err
	}
	defer pmtInfCtx.Free()

	if err = pmtInfCtx.RegisterNS("ns", "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"); err != nil {
		return nil, err
	}

	pi = new(PmtInf)

	pi.PmtInfID = xpath.String(pmtInfCtx.Find("ns:PmtInfId"))
	pi.NbOfTxs = xpath.String(pmtInfCtx.Find("ns:NbOfTxs"))
	pi.CtrlSum = xpath.String(pmtInfCtx.Find("ns:CtrlSum"))
	pi.ReqdExctnDt = xpath.String(pmtInfCtx.Find("ns:ReqdExctnDt"))
	pi.Dbtr = Account{
		Name: xpath.String(pmtInfCtx.Find("ns:Dbtr/ns:Nm")),
		IBAN: xpath.String(pmtInfCtx.Find("ns:DbtrAcct/ns:Id/ns:IBAN")),
		BIC:  xpath.String(pmtInfCtx.Find("ns:DbtrAgt/ns:FinInstnId/ns:BIC")),
	}

	return pi, nil
}
