package routers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"
	"github.com/alanwade2001/spa-submissions-api/repositories"
	"github.com/alanwade2001/spa-submissions-api/routers"
	"github.com/alanwade2001/spa-submissions-api/services"
	"github.com/gin-gonic/gin"
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
