package repositories

import (
	"errors"
	"reflect"
	"time"

	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mgo "github.com/alanwade2001/spa-common/mongo"

	"k8s.io/klog/v2"
)

// MongoRepository s
type MongoRepository struct {
	service *mgo.MongoService
}

// NewMongoService s
func NewMongoService() types.RepositoryAPI {
	return &MongoRepository{}
}

func (ms *MongoRepository) GetService() *mgo.MongoService {

	if ms.service != nil {
		return ms.service
	}

	uriTemplate := viper.GetString("MONGODB_URI_TEMPLATE")
	username := viper.GetString("MONGODB_USER")
	password := viper.GetString("MONGODB_PASSWORD")
	connectTimeout := viper.GetDuration("MONGODB_TIMEOUT") * time.Second
	database := viper.GetString("MONGODB_DATABASE")
	collection := viper.GetString("MONGODB_COLLECTION")

	structcodec, _ := bsoncodec.NewStructCodec(bsoncodec.JSONFallbackStructTagParser)
	reg := bson.NewRegistryBuilder().
		RegisterTypeEncoder(reflect.TypeOf(submission.SubmissionModel{}), structcodec).
		RegisterTypeDecoder(reflect.TypeOf(submission.SubmissionModel{}), structcodec).
		Build()

	service := mgo.NewMongoService(uriTemplate, username, password, database, collection, connectTimeout, reg)

	ms.service = service

	return ms.service
}

// CreateSubmission f
func (ms MongoRepository) CreateSubmission(submission *submission.SubmissionModel) (*submission.SubmissionModel, error) {
	connection := ms.GetService().Connect()
	defer connection.Disconnect()

	submission.Id = primitive.NewObjectID().Hex()

	result, err := ms.GetService().GetCollection(connection).InsertOne(connection.Ctx, submission)

	if err != nil {
		klog.Warningf("Could not create Submission: %v", err)
		return nil, err
	}

	klog.Infof("result:[%+v]", result)

	return submission, nil
}

// GetSubmission f
func (ms MongoRepository) GetSubmission(ID string) (sub *submission.SubmissionModel, err error) {
	connection := ms.GetService().Connect()
	defer connection.Disconnect()

	sub = new(submission.SubmissionModel)

	filter := bson.M{"_id": ID}

	if err := ms.GetService().GetCollection(connection).FindOne(connection.Ctx, filter).Decode(sub); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	klog.Infof("submission:[%+v]", sub)

	return sub, nil
}

// GetSubmissions f
func (ms MongoRepository) GetSubmissions() ([]*submission.SubmissionModel, error) {
	connection := ms.GetService().Connect()
	defer connection.Disconnect()

	var cursor *mongo.Cursor
	var err error
	submissions := []*submission.SubmissionModel{}

	filter := bson.M{}
	if cursor, err = ms.GetService().GetCollection(connection).Find(connection.Ctx, filter); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return submissions, nil
		}

		return nil, err
	}

	if err = cursor.All(connection.Ctx, &submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}
