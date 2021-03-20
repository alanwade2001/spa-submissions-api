package repositories

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/alanwade2001/spa-submissions-api/models/generated/submission"
	"github.com/alanwade2001/spa-submissions-api/types"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"k8s.io/klog/v2"
)

// MongoService s
type MongoService struct {
}

// NewMongoService s
func NewMongoService() types.RepositoryAPI {
	return &MongoService{}
}

// MongoConnection s
type MongoConnection struct {
	client *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// Disconnect f
func (mc *MongoConnection) Disconnect() {
	mc.cancel()
	mc.client.Disconnect(mc.ctx)
}

// Connect f
func (ms MongoService) connect() MongoConnection {
	username := viper.GetString("MONGODB_USER")
	klog.Infof("mongo user: [%s]", username)
	password := viper.GetString("MONGODB_PASSWORD")
	uriTemplate := viper.GetString("MONGODB_URI_TEMPLATE")
	klog.Infof("uriTemplate: [%s]", uriTemplate)

	connectionURI := fmt.Sprintf(uriTemplate, username, password)

	structcodec, _ := bsoncodec.NewStructCodec(bsoncodec.JSONFallbackStructTagParser)
	reg := bson.NewRegistryBuilder().
		RegisterTypeEncoder(reflect.TypeOf(submission.SubmissionModel{}), structcodec).
		RegisterTypeDecoder(reflect.TypeOf(submission.SubmissionModel{}), structcodec).
		Build()

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI).SetRegistry(reg))
	if err != nil {
		klog.Warningf("Failed to create client: %v", err)
	}

	connectTimeout := viper.GetDuration("MONGODB_TIMEOUT") * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)

	err = client.Connect(ctx)
	if err != nil {
		klog.Warningf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		klog.Warningf("Failed to ping cluster: %v", err)
	}

	klog.Infof("Connected to MongoDB!")

	return MongoConnection{client, ctx, cancel}
}

func (ms MongoService) getCollection(connection MongoConnection) *mongo.Collection {
	database := viper.GetString("MONGODB_DATABASE")
	return connection.client.Database(database).Collection("Submissions")
}

// CreateSubmission f
func (ms MongoService) CreateSubmission(submission *submission.SubmissionModel) (*submission.SubmissionModel, error) {
	connection := ms.connect()
	defer connection.Disconnect()

	submission.Id = primitive.NewObjectID().Hex()

	result, err := ms.getCollection(connection).InsertOne(connection.ctx, submission)

	if err != nil {
		klog.Warningf("Could not create Submission: %v", err)
		return nil, err
	}

	klog.Infof("result:[%+v]", result)

	return submission, nil
}

// GetSubmission f
func (ms MongoService) GetSubmission(ID string) (sub *submission.SubmissionModel, err error) {
	connection := ms.connect()
	defer connection.Disconnect()

	sub = new(submission.SubmissionModel)

	filter := bson.M{"_id": ID}

	if err := ms.getCollection(connection).FindOne(connection.ctx, filter).Decode(sub); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	klog.Infof("submission:[%+v]", sub)

	return sub, nil
}

// GetSubmissions f
func (ms MongoService) GetSubmissions() ([]*submission.SubmissionModel, error) {
	connection := ms.connect()
	defer connection.Disconnect()

	var cursor *mongo.Cursor
	var err error
	submissions := []*submission.SubmissionModel{}

	filter := bson.M{}
	if cursor, err = ms.getCollection(connection).Find(connection.ctx, filter); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return submissions, nil
		}

		return nil, err
	}

	if err = cursor.All(connection.ctx, &submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}
