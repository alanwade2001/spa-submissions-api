package routers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"
	"github.com/alanwade2001/spa-submissions-api/repositories"
	"github.com/alanwade2001/spa-submissions-api/routers"
	"github.com/alanwade2001/spa-submissions-api/services"
	"github.com/alanwade2001/spa-submissions-api/types/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

func TestSubmissionRouter_GetSubmissions(t *testing.T) {
	// change the database to be unittest
	os.Setenv("MONGODB_DATABASE", "unittest")
	os.Setenv("PAIN_001_XSD_FILE", "../schemas/pain.001.001.03.xsd")

	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	userAPI := services.NewUserService()
	idGeneratorAPI := services.NewMongoIdGenerator()
	repositoryAPI := repositories.NewMongoService()
	xmlParserAPI := services.NewXMLParserAPI()
	groupHeaderMapperAPI := services.NewGroupHeaderMapper()
	paymentInformationMapperAPI := services.NewPaymentInformationMapper()
	initiationMapperAPI := services.NewInitiationMapper(groupHeaderMapperAPI, paymentInformationMapperAPI)
	initiationAPI := services.NewInitiationService(xmlParserAPI, initiationMapperAPI)
	validatorAPI := services.NewValidator()
	customerAPI := services.NewCustomerService()
	messageAPI := services.NewMessageService()
	submissionServiceAPI := services.NewSubmissionService(idGeneratorAPI, repositoryAPI, initiationAPI, validatorAPI, customerAPI, messageAPI)
	submissionAPI := routers.NewSubmissionRouter(userAPI, submissionServiceAPI)

	services.NewConfigService().Load("..")

	engine.GET("/submissions", submissionAPI.GetSubmissions)

	type data struct {
		submissions []interface{}
	}

	tests := []struct {
		name string
		data data
	}{
		{
			name: "Test01",
			data: data{
				submissions: []interface{}{},
			},
		},
		{
			name: "Test02",
			data: data{
				submissions: []interface{}{
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234567",
					},
				},
			},
		},
		{
			name: "Test03",
			data: data{
				submissions: []interface{}{
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234567",
					},
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234568",
					},
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234569",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// clear down the database
			mongoRepo := repositoryAPI.(*repositories.MongoRepository)
			conn := mongoRepo.GetService().Connect()
			defer conn.Disconnect()

			// clear the database
			filter := bson.M{}
			if _, err := mongoRepo.GetService().GetCollection(conn).DeleteMany(conn.Ctx, filter); err != nil {
				t.Fatalf("error deleting submissions [%s]", err.Error())
			}

			if len(tt.data.submissions) > 0 {
				// insert the seed data
				if _, err := mongoRepo.GetService().GetCollection(conn).InsertMany(conn.Ctx, tt.data.submissions); err != nil {
					t.Fatalf("error inserting submitting [%s]", err.Error())
				}
			}

			req, err := http.NewRequest("GET", "/submissions", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder so you can inspect the response
			w := httptest.NewRecorder()

			// Perform the request
			engine.ServeHTTP(w, req)
			//fmt.Println(w.Body)

			// Check to see if the response was what you expected
			if w.Code == http.StatusOK {
				t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
			} else {
				t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
			}

			var result []submission.SubmissionModel
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatalf("Expected Json output")
			}

			// check if the size of the array is equal to the data
			if len(result) == len(tt.data.submissions) {
				t.Logf("Expected to get number of submissions %d is same as %d\n", len(tt.data.submissions), len(result))
			} else {
				t.Fatalf("Expected to get number of submissions %d but instead got %d\n", len(tt.data.submissions), len(result))
			}

			// random check on the first customer of it exists
			if len(result) > 0 {
				first := tt.data.submissions[0].(submission.SubmissionModel)
				if result[0].Id == first.Id {
					t.Logf("Expected id to match %s is same as %s\n", first.Id, result[0].Id)
				} else {
					t.Fatalf("Expected id to match %s but instead got %s\n", first.Id, result[0].Id)
				}
			} else {
				t.Logf("No customers in the result")
			}
		})
	}
}

func TestSubmissionRouter_GetSubmission(t *testing.T) {
	// change the database to be unittest
	os.Setenv("MONGODB_DATABASE", "unittest")
	os.Setenv("PAIN_001_XSD_FILE", "../schemas/pain.001.001.03.xsd")

	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	userAPI := services.NewUserService()
	idGeneratorAPI := services.NewMongoIdGenerator()
	repositoryAPI := repositories.NewMongoService()
	xmlParserAPI := services.NewXMLParserAPI()
	groupHeaderMapperAPI := services.NewGroupHeaderMapper()
	paymentInformationMapperAPI := services.NewPaymentInformationMapper()
	initiationMapperAPI := services.NewInitiationMapper(groupHeaderMapperAPI, paymentInformationMapperAPI)
	initiationAPI := services.NewInitiationService(xmlParserAPI, initiationMapperAPI)
	validatorAPI := services.NewValidator()
	customerAPI := services.NewCustomerService()
	messageAPI := services.NewMessageService()
	submissionServiceAPI := services.NewSubmissionService(idGeneratorAPI, repositoryAPI, initiationAPI, validatorAPI, customerAPI, messageAPI)
	submissionAPI := routers.NewSubmissionRouter(userAPI, submissionServiceAPI)

	services.NewConfigService().Load("..")

	engine.GET("/submissions/:id", submissionAPI.GetSubmissions)

	type data struct {
		submissions []interface{}
	}

	tests := []struct {
		name  string
		id    string
		data  data
		code  int
		index int
	}{
		{
			name: "Test01",
			id:   "123456",
			data: data{
				submissions: []interface{}{},
			},
			code:  http.StatusNotFound,
			index: -1,
		},
		{
			name: "Test02",
			id:   "sub_1234567",
			data: data{
				submissions: []interface{}{
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234567",
					},
				},
			},
			code:  http.StatusOK,
			index: 0,
		},
		{
			name: "Test03",
			id:   "sub_1234568",
			data: data{
				submissions: []interface{}{
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234567",
					},
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234568",
					},
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234569",
					},
				},
			},
			code:  http.StatusOK,
			index: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// clear down the database
			mongoRepo := repositoryAPI.(*repositories.MongoRepository)
			conn := mongoRepo.GetService().Connect()
			defer conn.Disconnect()

			// clear the database
			filter := bson.M{}
			if _, err := mongoRepo.GetService().GetCollection(conn).DeleteMany(conn.Ctx, filter); err != nil {
				t.Fatalf("error deleting submissions [%s]", err.Error())
			}

			if len(tt.data.submissions) > 0 {
				// insert the seed data
				if _, err := mongoRepo.GetService().GetCollection(conn).InsertMany(conn.Ctx, tt.data.submissions); err != nil {
					t.Fatalf("error inserting submitting [%s]", err.Error())
				}
			}

			req, err := http.NewRequest("GET", "/submissions/"+tt.id, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder so you can inspect the response
			w := httptest.NewRecorder()

			// Perform the request
			engine.ServeHTTP(w, req)
			//fmt.Println(w.Body)

			// Check to see if the response was what you expected
			if w.Code == http.StatusOK {
				t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
			} else {
				t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
			}

			var result []submission.SubmissionModel
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatalf("Expected Json output")
			}

			// check if the size of the array is equal to the data
			if len(result) == len(tt.data.submissions) {
				t.Logf("Expected to get number of submissions %d is same as %d\n", len(tt.data.submissions), len(result))
			} else {
				t.Fatalf("Expected to get number of submissions %d but instead got %d\n", len(tt.data.submissions), len(result))
			}

			// random check on the first customer of it exists
			if len(result) > 0 {
				first := tt.data.submissions[0].(submission.SubmissionModel)
				if result[0].Id == first.Id {
					t.Logf("Expected id to match %s is same as %s\n", first.Id, result[0].Id)
				} else {
					t.Fatalf("Expected id to match %s but instead got %s\n", first.Id, result[0].Id)
				}
			} else {
				t.Logf("No customers in the result")
			}
		})
	}
}

func TestSubmissionRouter_CreateSubmission(t *testing.T) {

	// change the database to be unittest
	os.Setenv("MONGODB_DATABASE", "unittest")
	os.Setenv("PAIN_001_XSD_FILE", "../schemas/pain.001.001.03.xsd")

	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	userAPI := &mocks.UserAPI{}
	idGeneratorAPI := services.NewMongoIdGenerator()
	repositoryAPI := repositories.NewMongoService()
	xmlParserAPI := services.NewXMLParserAPI()
	groupHeaderMapperAPI := services.NewGroupHeaderMapper()
	paymentInformationMapperAPI := services.NewPaymentInformationMapper()
	initiationMapperAPI := services.NewInitiationMapper(groupHeaderMapperAPI, paymentInformationMapperAPI)
	initiationAPI := services.NewInitiationService(xmlParserAPI, initiationMapperAPI)
	validatorAPI := services.NewValidator()
	customerAPI := &mocks.CustomerAPI{}
	messageAPI := &mocks.MessageAPI{}
	submissionServiceAPI := services.NewSubmissionService(idGeneratorAPI, repositoryAPI, initiationAPI, validatorAPI, customerAPI, messageAPI)
	submissionAPI := routers.NewSubmissionRouter(userAPI, submissionServiceAPI)

	services.NewConfigService().Load("..")

	engine.POST("/submissions", submissionAPI.CreateSubmission)

	type data struct {
		submissions   []interface{}
		newSubmission string
	}

	tests := []struct {
		name     string
		data     data
		user     *submission.UserReference
		customer *submission.CustomerReference
		code     int
		location string
	}{
		{
			name: "Test01",
			user: &submission.UserReference{Email: "alan@test.ie"},
			customer: &submission.CustomerReference{
				Id:   "cust_1234",
				Name: "Corporation ABC",
			},
			data: data{
				submissions: []interface{}{},
				newSubmission: `<?xml version="1.0" encoding="UTF-8"?>
<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pain.001.001.03">
    <CstmrCdtTrfInitn>
        <GrpHdr>
            <MsgId>message-id-002</MsgId>
            <CreDtTm>2010-09-28T14:07:00</CreDtTm>
            <NbOfTxs>2</NbOfTxs>
            <CtrlSum>30.3</CtrlSum>
            <InitgPty>
                <Nm>Bedrijfsnaam</Nm>
                <Id>
                    <OrgId>
                        <Othr>
                            <Id>123456789012345</Id>
                        </Othr>
                    </OrgId>
                </Id>
            </InitgPty>
        </GrpHdr>
        <PmtInf>
            <PmtInfId>minimaal gevuld</PmtInfId>
            <PmtMtd>TRF</PmtMtd>
            <NbOfTxs>1</NbOfTxs>
            <CtrlSum>10.1</CtrlSum>
            <ReqdExctnDt>2009-11-01</ReqdExctnDt>
            <Dbtr>
                <Nm>Naam</Nm>
            </Dbtr>
            <DbtrAcct>
                <Id>
                    <IBAN>IE12AIBK12345678901234</IBAN>
                </Id>
            </DbtrAcct>
            <DbtrAgt>
                <FinInstnId>
                    <BIC>RABONL2U</BIC>
                </FinInstnId>
            </DbtrAgt>
            <CdtTrfTxInf>
                <PmtId>
                    <EndToEndId>non ref</EndToEndId>
                </PmtId>
                <Amt>
                    <InstdAmt Ccy="EUR">10.1</InstdAmt>
                </Amt>
                <ChrgBr>SLEV</ChrgBr>
                <CdtrAgt>
                    <FinInstnId>
                        <BIC>ABNANL2A</BIC>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Nm>Naam creditor</Nm>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <IBAN>NL90ABNA0111111111</IBAN>
                    </Id>
                </CdtrAcct>
                <RmtInf>
                    <Ustrd>vrije tekst</Ustrd>
                </RmtInf>
            </CdtTrfTxInf>
        </PmtInf>
        <PmtInf>
            <PmtInfId>maximaal gevuld</PmtInfId>
            <PmtMtd>TRF</PmtMtd>
            <BtchBookg>true</BtchBookg>
            <NbOfTxs>1</NbOfTxs>
            <CtrlSum>20.2</CtrlSum>
            <PmtTpInf>
                <InstrPrty>NORM</InstrPrty>
                <SvcLvl>
                    <Cd>SEPA</Cd>
                </SvcLvl>
                <LclInstrm>
                    <Cd>IDEAL</Cd>
                </LclInstrm>
                <CtgyPurp>
                    <Cd>SECU</Cd>
                </CtgyPurp>
            </PmtTpInf>
            <ReqdExctnDt>2009-11-01</ReqdExctnDt>
            <Dbtr>
                <Nm>Naam</Nm>
                <PstlAdr>
                    <Ctry>NL</Ctry>
                    <AdrLine>Debtor straat 1</AdrLine>
                    <AdrLine>9999 XX Plaats debtor</AdrLine>
                </PstlAdr>
            </Dbtr>
            <DbtrAcct>
                <Id>
                    <IBAN>IE13AIBK23456789012345</IBAN>
                </Id>
            </DbtrAcct>
            <DbtrAgt>
                <FinInstnId>
                    <BIC>RABONL2U</BIC>
                </FinInstnId>
            </DbtrAgt>
            <UltmtDbtr>
                <Id>
                    <OrgId>
                        <Othr>
                            <Id>12345678</Id>
                            <SchmeNm>
                                <Prtry>klantnummer</Prtry>
                            </SchmeNm>
                            <Issr>klantnummer uitgifte instantie</Issr>
                        </Othr>
                    </OrgId>
                </Id>
            </UltmtDbtr>
            <ChrgBr>SLEV</ChrgBr>
            <CdtTrfTxInf>
                <PmtId>
                    <InstrId>debtor-to-debtor-bank-01</InstrId>
                    <EndToEndId>End-to-end-id-debtor-to-creditor-01</EndToEndId>
                </PmtId>
                <Amt>
                    <InstdAmt Ccy="EUR">20.2</InstdAmt>
                </Amt>
                <CdtrAgt>
                    <FinInstnId>
                        <BIC>ABNANL2A</BIC>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Nm>Naam creditor</Nm>
                    <PstlAdr>
                        <Ctry>NL</Ctry>
                        <AdrLine>Straat creditor 1</AdrLine>
                        <AdrLine>9999 XX Plaats creditor</AdrLine>
                    </PstlAdr>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <IBAN>NL90ABNA0111111111</IBAN>
                    </Id>
                </CdtrAcct>
                <UltmtCdtr>
                    <Id>
                        <PrvtId>
                            <DtAndPlcOfBirth>
                                <BirthDt>1969-07-03</BirthDt>
                                <CityOfBirth>PLAATS</CityOfBirth>
                                <CtryOfBirth>NL</CtryOfBirth>
                            </DtAndPlcOfBirth>
                        </PrvtId>
                    </Id>
                </UltmtCdtr>
                <Purp>
                    <Cd>CHAR</Cd>
                </Purp>
                <RmtInf>
                    <Strd>
                        <CdtrRefInf>
                            <Tp>
                                <CdOrPrtry>
                                    <Cd>SCOR</Cd>
                                </CdOrPrtry>
                                <Issr>CUR</Issr>
                            </Tp>
                            <Ref>1234567</Ref>
                        </CdtrRefInf>
                    </Strd>
                </RmtInf>
            </CdtTrfTxInf>
        </PmtInf>
    </CstmrCdtTrfInitn>
</Document>`,
			},
			code:     http.StatusCreated,
			location: "/customers/.+",
		},
		{
			name: "Test02",
			user: &submission.UserReference{Email: "alan@test.ie"},
			customer: &submission.CustomerReference{
				Id:   "cust_1234",
				Name: "Corporation ABC",
			},
			data: data{
				submissions: []interface{}{
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234567",
					},
				},
				newSubmission: `<?xml version="1.0" encoding="UTF-8"?>
<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pain.001.001.03">
    <CstmrCdtTrfInitn>
        <GrpHdr>
            <MsgId>message-id-002</MsgId>
            <CreDtTm>2010-09-28T14:07:00</CreDtTm>
            <NbOfTxs>2</NbOfTxs>
            <CtrlSum>30.3</CtrlSum>
            <InitgPty>
                <Nm>Bedrijfsnaam</Nm>
                <Id>
                    <OrgId>
                        <Othr>
                            <Id>123456789012345</Id>
                        </Othr>
                    </OrgId>
                </Id>
            </InitgPty>
        </GrpHdr>
        <PmtInf>
            <PmtInfId>minimaal gevuld</PmtInfId>
            <PmtMtd>TRF</PmtMtd>
            <NbOfTxs>1</NbOfTxs>
            <CtrlSum>10.1</CtrlSum>
            <ReqdExctnDt>2009-11-01</ReqdExctnDt>
            <Dbtr>
                <Nm>Naam</Nm>
            </Dbtr>
            <DbtrAcct>
                <Id>
                    <IBAN>IE12AIBK12345678901234</IBAN>
                </Id>
            </DbtrAcct>
            <DbtrAgt>
                <FinInstnId>
                    <BIC>RABONL2U</BIC>
                </FinInstnId>
            </DbtrAgt>
            <CdtTrfTxInf>
                <PmtId>
                    <EndToEndId>non ref</EndToEndId>
                </PmtId>
                <Amt>
                    <InstdAmt Ccy="EUR">10.1</InstdAmt>
                </Amt>
                <ChrgBr>SLEV</ChrgBr>
                <CdtrAgt>
                    <FinInstnId>
                        <BIC>ABNANL2A</BIC>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Nm>Naam creditor</Nm>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <IBAN>NL90ABNA0111111111</IBAN>
                    </Id>
                </CdtrAcct>
                <RmtInf>
                    <Ustrd>vrije tekst</Ustrd>
                </RmtInf>
            </CdtTrfTxInf>
        </PmtInf>
        <PmtInf>
            <PmtInfId>maximaal gevuld</PmtInfId>
            <PmtMtd>TRF</PmtMtd>
            <BtchBookg>true</BtchBookg>
            <NbOfTxs>1</NbOfTxs>
            <CtrlSum>20.2</CtrlSum>
            <PmtTpInf>
                <InstrPrty>NORM</InstrPrty>
                <SvcLvl>
                    <Cd>SEPA</Cd>
                </SvcLvl>
                <LclInstrm>
                    <Cd>IDEAL</Cd>
                </LclInstrm>
                <CtgyPurp>
                    <Cd>SECU</Cd>
                </CtgyPurp>
            </PmtTpInf>
            <ReqdExctnDt>2009-11-01</ReqdExctnDt>
            <Dbtr>
                <Nm>Naam</Nm>
                <PstlAdr>
                    <Ctry>NL</Ctry>
                    <AdrLine>Debtor straat 1</AdrLine>
                    <AdrLine>9999 XX Plaats debtor</AdrLine>
                </PstlAdr>
            </Dbtr>
            <DbtrAcct>
                <Id>
                    <IBAN>IE13AIBK23456789012345</IBAN>
                </Id>
            </DbtrAcct>
            <DbtrAgt>
                <FinInstnId>
                    <BIC>RABONL2U</BIC>
                </FinInstnId>
            </DbtrAgt>
            <UltmtDbtr>
                <Id>
                    <OrgId>
                        <Othr>
                            <Id>12345678</Id>
                            <SchmeNm>
                                <Prtry>klantnummer</Prtry>
                            </SchmeNm>
                            <Issr>klantnummer uitgifte instantie</Issr>
                        </Othr>
                    </OrgId>
                </Id>
            </UltmtDbtr>
            <ChrgBr>SLEV</ChrgBr>
            <CdtTrfTxInf>
                <PmtId>
                    <InstrId>debtor-to-debtor-bank-01</InstrId>
                    <EndToEndId>End-to-end-id-debtor-to-creditor-01</EndToEndId>
                </PmtId>
                <Amt>
                    <InstdAmt Ccy="EUR">20.2</InstdAmt>
                </Amt>
                <CdtrAgt>
                    <FinInstnId>
                        <BIC>ABNANL2A</BIC>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Nm>Naam creditor</Nm>
                    <PstlAdr>
                        <Ctry>NL</Ctry>
                        <AdrLine>Straat creditor 1</AdrLine>
                        <AdrLine>9999 XX Plaats creditor</AdrLine>
                    </PstlAdr>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <IBAN>NL90ABNA0111111111</IBAN>
                    </Id>
                </CdtrAcct>
                <UltmtCdtr>
                    <Id>
                        <PrvtId>
                            <DtAndPlcOfBirth>
                                <BirthDt>1969-07-03</BirthDt>
                                <CityOfBirth>PLAATS</CityOfBirth>
                                <CtryOfBirth>NL</CtryOfBirth>
                            </DtAndPlcOfBirth>
                        </PrvtId>
                    </Id>
                </UltmtCdtr>
                <Purp>
                    <Cd>CHAR</Cd>
                </Purp>
                <RmtInf>
                    <Strd>
                        <CdtrRefInf>
                            <Tp>
                                <CdOrPrtry>
                                    <Cd>SCOR</Cd>
                                </CdOrPrtry>
                                <Issr>CUR</Issr>
                            </Tp>
                            <Ref>1234567</Ref>
                        </CdtrRefInf>
                    </Strd>
                </RmtInf>
            </CdtTrfTxInf>
        </PmtInf>
    </CstmrCdtTrfInitn>
</Document>`,
			},
			code:     http.StatusCreated,
			location: "/customers/.+",
		},
		{
			name: "Test03",
			user: &submission.UserReference{Email: "alan@test.ie"},
			customer: &submission.CustomerReference{
				Id:   "cust_1234",
				Name: "Corporation ABC",
			},
			data: data{
				submissions: []interface{}{
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234567",
					},
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234568",
					},
					submission.SubmissionModel{
						Customer: &submission.CustomerReference{
							Id:   "cust_11223344",
							Name: "Corporation ABC",
						},
						GroupHeader: &submission.GroupHeaderReference{
							ControlSum:           64.0,
							CreationDateTime:     "2020-01-01T12:12:35",
							InitiatingPartyId:    "initpty_112233",
							MessageId:            "msg-id",
							NumberOfTransactions: 1,
						},
						SubmittedAt: "2020-01-01T12:12:35",
						SubmittedBy: &submission.UserReference{
							Email: "alan@test.ie",
						},
						Id: "sub_1234569",
					},
				},
				newSubmission: `<?xml version="1.0" encoding="UTF-8"?>
<Document xmlns="urn:iso:std:iso:20022:tech:xsd:pain.001.001.03">
    <CstmrCdtTrfInitn>
        <GrpHdr>
            <MsgId>message-id-002</MsgId>
            <CreDtTm>2010-09-28T14:07:00</CreDtTm>
            <NbOfTxs>2</NbOfTxs>
            <CtrlSum>30.3</CtrlSum>
            <InitgPty>
                <Nm>Bedrijfsnaam</Nm>
                <Id>
                    <OrgId>
                        <Othr>
                            <Id>123456789012345</Id>
                        </Othr>
                    </OrgId>
                </Id>
            </InitgPty>
        </GrpHdr>
        <PmtInf>
            <PmtInfId>minimaal gevuld</PmtInfId>
            <PmtMtd>TRF</PmtMtd>
            <NbOfTxs>1</NbOfTxs>
            <CtrlSum>10.1</CtrlSum>
            <ReqdExctnDt>2009-11-01</ReqdExctnDt>
            <Dbtr>
                <Nm>Naam</Nm>
            </Dbtr>
            <DbtrAcct>
                <Id>
                    <IBAN>IE12AIBK12345678901234</IBAN>
                </Id>
            </DbtrAcct>
            <DbtrAgt>
                <FinInstnId>
                    <BIC>RABONL2U</BIC>
                </FinInstnId>
            </DbtrAgt>
            <CdtTrfTxInf>
                <PmtId>
                    <EndToEndId>non ref</EndToEndId>
                </PmtId>
                <Amt>
                    <InstdAmt Ccy="EUR">10.1</InstdAmt>
                </Amt>
                <ChrgBr>SLEV</ChrgBr>
                <CdtrAgt>
                    <FinInstnId>
                        <BIC>ABNANL2A</BIC>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Nm>Naam creditor</Nm>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <IBAN>NL90ABNA0111111111</IBAN>
                    </Id>
                </CdtrAcct>
                <RmtInf>
                    <Ustrd>vrije tekst</Ustrd>
                </RmtInf>
            </CdtTrfTxInf>
        </PmtInf>
        <PmtInf>
            <PmtInfId>maximaal gevuld</PmtInfId>
            <PmtMtd>TRF</PmtMtd>
            <BtchBookg>true</BtchBookg>
            <NbOfTxs>1</NbOfTxs>
            <CtrlSum>20.2</CtrlSum>
            <PmtTpInf>
                <InstrPrty>NORM</InstrPrty>
                <SvcLvl>
                    <Cd>SEPA</Cd>
                </SvcLvl>
                <LclInstrm>
                    <Cd>IDEAL</Cd>
                </LclInstrm>
                <CtgyPurp>
                    <Cd>SECU</Cd>
                </CtgyPurp>
            </PmtTpInf>
            <ReqdExctnDt>2009-11-01</ReqdExctnDt>
            <Dbtr>
                <Nm>Naam</Nm>
                <PstlAdr>
                    <Ctry>NL</Ctry>
                    <AdrLine>Debtor straat 1</AdrLine>
                    <AdrLine>9999 XX Plaats debtor</AdrLine>
                </PstlAdr>
            </Dbtr>
            <DbtrAcct>
                <Id>
                    <IBAN>IE13AIBK23456789012345</IBAN>
                </Id>
            </DbtrAcct>
            <DbtrAgt>
                <FinInstnId>
                    <BIC>RABONL2U</BIC>
                </FinInstnId>
            </DbtrAgt>
            <UltmtDbtr>
                <Id>
                    <OrgId>
                        <Othr>
                            <Id>12345678</Id>
                            <SchmeNm>
                                <Prtry>klantnummer</Prtry>
                            </SchmeNm>
                            <Issr>klantnummer uitgifte instantie</Issr>
                        </Othr>
                    </OrgId>
                </Id>
            </UltmtDbtr>
            <ChrgBr>SLEV</ChrgBr>
            <CdtTrfTxInf>
                <PmtId>
                    <InstrId>debtor-to-debtor-bank-01</InstrId>
                    <EndToEndId>End-to-end-id-debtor-to-creditor-01</EndToEndId>
                </PmtId>
                <Amt>
                    <InstdAmt Ccy="EUR">20.2</InstdAmt>
                </Amt>
                <CdtrAgt>
                    <FinInstnId>
                        <BIC>ABNANL2A</BIC>
                    </FinInstnId>
                </CdtrAgt>
                <Cdtr>
                    <Nm>Naam creditor</Nm>
                    <PstlAdr>
                        <Ctry>NL</Ctry>
                        <AdrLine>Straat creditor 1</AdrLine>
                        <AdrLine>9999 XX Plaats creditor</AdrLine>
                    </PstlAdr>
                </Cdtr>
                <CdtrAcct>
                    <Id>
                        <IBAN>NL90ABNA0111111111</IBAN>
                    </Id>
                </CdtrAcct>
                <UltmtCdtr>
                    <Id>
                        <PrvtId>
                            <DtAndPlcOfBirth>
                                <BirthDt>1969-07-03</BirthDt>
                                <CityOfBirth>PLAATS</CityOfBirth>
                                <CtryOfBirth>NL</CtryOfBirth>
                            </DtAndPlcOfBirth>
                        </PrvtId>
                    </Id>
                </UltmtCdtr>
                <Purp>
                    <Cd>CHAR</Cd>
                </Purp>
                <RmtInf>
                    <Strd>
                        <CdtrRefInf>
                            <Tp>
                                <CdOrPrtry>
                                    <Cd>SCOR</Cd>
                                </CdOrPrtry>
                                <Issr>CUR</Issr>
                            </Tp>
                            <Ref>1234567</Ref>
                        </CdtrRefInf>
                    </Strd>
                </RmtInf>
            </CdtTrfTxInf>
        </PmtInf>
    </CstmrCdtTrfInitn>
</Document>`,
			},
			code:     http.StatusCreated,
			location: "/customers/.+",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// clear down the database
			mongoRepo := repositoryAPI.(*repositories.MongoRepository)
			conn := mongoRepo.GetService().Connect()
			defer conn.Disconnect()

			// clear the database
			filter := bson.M{}
			if _, err := mongoRepo.GetService().GetCollection(conn).DeleteMany(conn.Ctx, filter); err != nil {
				t.Logf("error deleting customers [%s]", err.Error())
			}

			if len(tt.data.submissions) > 0 {
				// insert the seed data
				if _, err := mongoRepo.GetService().GetCollection(conn).InsertMany(conn.Ctx, tt.data.submissions); err != nil {
					t.Logf("error inserting customers [%s]", err.Error())
				}
			}

			// mocks
			userAPI.On("Find", mock.Anything).Return(tt.user, nil)

			customerAPI.On("Find", mock.Anything).Return(tt.customer, nil)
			messageAPI.On("SendInitiation", mock.Anything).Return(nil)

			req, err := http.NewRequest(http.MethodPost, "/submissions", strings.NewReader(tt.data.newSubmission))
			if err != nil {
				t.Fatal(err)
			}

			// Create a response recorder so you can inspect the response
			w := httptest.NewRecorder()

			// Perform the request
			engine.ServeHTTP(w, req)
			//fmt.Println(w.Body)

			// Check to see if the response was what you expected
			if w.Code == tt.code {
				t.Logf("Expected to get status %d is same ast %d\n", tt.code, w.Code)
			} else {
				t.Fatalf("Expected to get status %d but instead got %d\n", tt.code, w.Code)
			}

			userAPI.AssertExpectations(t)
			customerAPI.AssertExpectations(t)
			messageAPI.AssertExpectations(t)
		})
	}
}
