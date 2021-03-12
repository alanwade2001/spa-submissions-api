package main

import (
	"fmt"

	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"

	spatypes "github.com/alanwade2001/spa-common"
)

// InitiationMapper s
type InitiationMapper struct {
	groupHeaderMapperAPI        GroupHeaderMapperAPI
	paymentInformationMapperAPI PaymentInformationMapperAPI
}

// NewInitiationMapper f
func NewInitiationMapper(groupHeaderMapperAPI GroupHeaderMapperAPI, paymentInformationMapperAPI PaymentInformationMapperAPI) InitiationMapperAPI {
	return InitiationMapper{
		groupHeaderMapperAPI:        groupHeaderMapperAPI,
		paymentInformationMapperAPI: paymentInformationMapperAPI,
	}
}

// Map f
func (m InitiationMapper) Map(doc types.Document) (*spatypes.Initiation, error) {
	//var root types.Node
	//var ctx *xpath.Context
	var grpHdr *spatypes.GroupHeader
	var pmtInfs *[]spatypes.PaymentInstruction

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

	initiation := spatypes.Initiation{
		GroupHeader:         *grpHdr,
		PaymentInstructions: *pmtInfs,
	}

	return &initiation, nil
}

// GroupHeaderMapper s
type GroupHeaderMapper struct {
}

// NewGroupHeaderMapper f
func NewGroupHeaderMapper() GroupHeaderMapperAPI {
	return GroupHeaderMapper{}
}

// Map f
func (ghm GroupHeaderMapper) Map(ctx *xpath.Context) (gh *spatypes.GroupHeader, err error) {
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

	gh = new(spatypes.GroupHeader)
	gh.MessageID = xpath.String(grpHdrCtx.Find("ns:MsgId"))
	gh.CreationDateTime = xpath.String(grpHdrCtx.Find("ns:CreDtTm"))
	gh.NumberOfTransactions = xpath.String(grpHdrCtx.Find("ns:NbOfTxs"))
	gh.ControlSum = xpath.String(grpHdrCtx.Find("ns:CtrlSum"))
	gh.InitiatingParty.InitiatingPartyID = xpath.String(grpHdrCtx.Find("ns:InitgPty/ns:Id/ns:OrgId/ns:Othr/ns:Id"))

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
func (pim PaymentInformationMapper) Map(ctx *xpath.Context) (pis *[]spatypes.PaymentInstruction, err error) {
	pmtInfNodes := xpath.NodeList(ctx.Find("/ns:Document/ns:CstmrCdtTrfInitn/ns:PmtInf"))

	pmtInfs := make([]spatypes.PaymentInstruction, len(pmtInfNodes))

	for i, n := range pmtInfNodes {
		var pi *spatypes.PaymentInstruction

		if pi, err = pim.MapPmtInf(n); err != nil {
			return nil, err
		}

		pmtInfs[i] = *pi
	}

	return &pmtInfs, nil
}

// MapPmtInf f
func (pim PaymentInformationMapper) MapPmtInf(pmtInfNode types.Node) (pi *spatypes.PaymentInstruction, err error) {

	pmtInfCtx, err := xpath.NewContext(pmtInfNode)
	if err != nil {
		return nil, err
	}
	defer pmtInfCtx.Free()

	if err = pmtInfCtx.RegisterNS("ns", "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"); err != nil {
		return nil, err
	}

	pi = new(spatypes.PaymentInstruction)

	pi.PaymentID = xpath.String(pmtInfCtx.Find("ns:PmtInfId"))
	pi.NumberOfTransactions = xpath.String(pmtInfCtx.Find("ns:NbOfTxs"))
	pi.ControlSum = xpath.String(pmtInfCtx.Find("ns:CtrlSum"))
	pi.RequestedExecutionDate = xpath.String(pmtInfCtx.Find("ns:ReqdExctnDt"))
	pi.Debtor = spatypes.AccountReference{
		Name: xpath.String(pmtInfCtx.Find("ns:Dbtr/ns:Nm")),
		IBAN: xpath.String(pmtInfCtx.Find("ns:DbtrAcct/ns:Id/ns:IBAN")),
		BIC:  xpath.String(pmtInfCtx.Find("ns:DbtrAgt/ns:FinInstnId/ns:BIC")),
	}

	return pi, nil
}
