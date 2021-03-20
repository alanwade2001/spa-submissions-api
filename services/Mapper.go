package services

import (
	"errors"
	"strconv"

	xml "github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"

	"github.com/alanwade2001/spa-submissions-api/models/generated/initiation"
	"github.com/alanwade2001/spa-submissions-api/types"
)

// InitiationMapper s
type InitiationMapper struct {
	groupHeaderMapperAPI        types.GroupHeaderMapperAPI
	paymentInformationMapperAPI types.PaymentInformationMapperAPI
}

// NewInitiationMapper f
func NewInitiationMapper(groupHeaderMapperAPI types.GroupHeaderMapperAPI, paymentInformationMapperAPI types.PaymentInformationMapperAPI) types.InitiationMapperAPI {
	return InitiationMapper{
		groupHeaderMapperAPI:        groupHeaderMapperAPI,
		paymentInformationMapperAPI: paymentInformationMapperAPI,
	}
}

// Map f
func (m InitiationMapper) Map(doc xml.Document) (*initiation.InitiationModel, error) {
	//var root types.Node
	//var ctx *xpath.Context
	var grpHdr *initiation.GroupHeaderReference
	var pmtInfs []*initiation.PaymentInstructionReference

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

	initiation := initiation.InitiationModel{
		GroupHeader:         grpHdr,
		PaymentInstructions: pmtInfs,
	}

	return &initiation, nil
}

// GroupHeaderMapper s
type GroupHeaderMapper struct {
}

// NewGroupHeaderMapper f
func NewGroupHeaderMapper() types.GroupHeaderMapperAPI {
	return GroupHeaderMapper{}
}

// Map f
func (ghm GroupHeaderMapper) Map(ctx *xpath.Context) (gh *initiation.GroupHeaderReference, err error) {
	var grpHdrIter xml.NodeIter

	if grpHdrIter = xpath.NodeIter(ctx.Find("/ns:Document/ns:CstmrCdtTrfInitn/ns:GrpHdr")); !grpHdrIter.Next() {
		return nil, errors.New("unable to find group header but it has passed xsd check")
	}

	grpHdrNode := grpHdrIter.Node()

	var grpHdrCtx *xpath.Context
	if grpHdrCtx, err = xpath.NewContext(grpHdrNode); err != nil {
		return nil, err
	} else if err = grpHdrCtx.RegisterNS("ns", "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"); err != nil {
		return nil, err
	}

	gh = new(initiation.GroupHeaderReference)
	gh.MessageId = xpath.String(grpHdrCtx.Find("ns:MsgId"))
	gh.CreationDateTime = xpath.String(grpHdrCtx.Find("ns:CreDtTm"))
	nbOfTxs := xpath.String(grpHdrCtx.Find("ns:NbOfTxs"))
	gh.NumberOfTransactions, _ = strconv.ParseFloat(nbOfTxs, 64)
	ctrlSum := xpath.String(grpHdrCtx.Find("ns:CtrlSum"))
	gh.ControlSum, _ = strconv.ParseFloat(ctrlSum, 64)
	gh.InitiatingPartyId = xpath.String(grpHdrCtx.Find("ns:InitgPty/ns:Id/ns:OrgId/ns:Othr/ns:Id"))

	return gh, nil
}

// PaymentInformationMapper s
type PaymentInformationMapper struct {
}

// NewPaymentInformationMapper f
func NewPaymentInformationMapper() types.PaymentInformationMapperAPI {
	return PaymentInformationMapper{}
}

// Map f
func (pim PaymentInformationMapper) Map(ctx *xpath.Context) (pis []*initiation.PaymentInstructionReference, err error) {
	pmtInfNodes := xpath.NodeList(ctx.Find("/ns:Document/ns:CstmrCdtTrfInitn/ns:PmtInf"))

	pmtInfs := make([]*initiation.PaymentInstructionReference, len(pmtInfNodes))

	for i, n := range pmtInfNodes {
		var pi *initiation.PaymentInstructionReference

		if pi, err = pim.MapPmtInf(n); err != nil {
			return nil, err
		}

		pmtInfs[i] = pi
	}

	return pmtInfs, nil
}

// MapPmtInf f
func (pim PaymentInformationMapper) MapPmtInf(pmtInfNode xml.Node) (pi *initiation.PaymentInstructionReference, err error) {

	pmtInfCtx, err := xpath.NewContext(pmtInfNode)
	if err != nil {
		return nil, err
	}
	defer pmtInfCtx.Free()

	if err = pmtInfCtx.RegisterNS("ns", "urn:iso:std:iso:20022:tech:xsd:pain.001.001.03"); err != nil {
		return nil, err
	}

	pi = new(initiation.PaymentInstructionReference)

	pi.PaymentId = xpath.String(pmtInfCtx.Find("ns:PmtInfId"))
	pi.NumberOfTransactions = xpath.Number(pmtInfCtx.Find("ns:NbOfTxs"))
	pi.ControlSum = xpath.Number(pmtInfCtx.Find("ns:CtrlSum"))
	pi.RequestedExecutionDate = xpath.String(pmtInfCtx.Find("ns:ReqdExctnDt"))
	pi.DebtorAccount = &initiation.AccountReference{
		Name: xpath.String(pmtInfCtx.Find("ns:Dbtr/ns:Nm")),
		IBAN: xpath.String(pmtInfCtx.Find("ns:DbtrAcct/ns:Id/ns:IBAN")),
		BIC:  xpath.String(pmtInfCtx.Find("ns:DbtrAgt/ns:FinInstnId/ns:BIC")),
	}

	return pi, nil
}
