package stores

import (
	"institute-person-api/src/config"
	"institute-person-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	CollectionName    string `json:"name"`
	Version           string `json:"value"`
	DefaultQuery      bson.M `json:"defaultQuery"`
	DefaultProjection bson.D `json:"defaultProjection"`
	config            *config.Config
	collection        mongo.Collection
}

/**
* Constants for common Bson Query and Projection values
 */
func MongoQueryNotVersion() bson.M {
	return bson.M{"name": bson.M{"$ne": "VERSION"}}
}

func MongoShortNameProjection() bson.D {
	return bson.D{{Key: "name", Value: 1}}
}

/**
* Construct a PersonStore to handle person database io
 */
func NewMongoStore(cfg *config.Config, collectionName string, query bson.M, projection bson.D) *MongoStore {
	this := &MongoStore{}

	// Initilize Store
	this.config = cfg
	this.CollectionName = collectionName
	this.DefaultQuery = query
	this.DefaultProjection = projection
	this.collection = cfg.GetCollection(collectionName)
	this.Version = this.GetVersion()

	// Put the database Version in the Config
	this.config.AddConfigStore(this.AsStoreItem())

	return this
}

/**
* Simple wrapper for mongo InsertOne
 */
func (this *MongoStore) InsertOne(insertValues bson.M) (*mongo.InsertOneResult, error) {
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	return this.collection.InsertOne(context, insertValues)
}

/**
* Simple wrapper for mongo Find One
 */
func (this *MongoStore) FindOne(query bson.M) *mongo.SingleResult {
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	return this.collection.FindOne(context, query)
}

/**
* Simple wrapper for mongo Find Many
 */
func (this *MongoStore) FindMany(query bson.M, options *options.FindOptions) (*mongo.Cursor, error) {
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	return this.collection.Find(context, query, options)
}

/**
* Simple wrapper for mongo Find One and Update
 */
func (this *MongoStore) FindOneAndUpdate(query bson.M, update bson.M, options *options.FindOneAndUpdateOptions) *mongo.SingleResult {
	ctx, cancel := this.config.GetTimeoutContext()
	defer cancel()
	return this.collection.FindOneAndUpdate(ctx, query, update, options)
}

/**
* Default Query by ID
 */
func (this *MongoStore) FindId(id string) (*interface{}, error) {
	// get the bson ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}

	var result interface{}
	var err error
	err = this.FindOne(query).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/**
* Default FindMany Query
 */
func (this *MongoStore) FindDocuments() ([]interface{}, error) {
	var results []interface{}
	var err error

	// Put the default projection into options
	options := options.Find().SetProjection(this.DefaultProjection)

	// Query the database
	cursor, err := this.FindMany(this.DefaultQuery, options)
	if err != nil {
		return nil, err
	}

	// Fetch all the results
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	err = cursor.All(context, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

/**
* Get the collection schema version
 */
func (this *MongoStore) GetVersion() string {
	var theVersion models.VersionInfo
	var err error

	query := bson.M{"name": "VERSION"}
	context, cancel := this.config.GetTimeoutContext()
	defer cancel()
	err = this.collection.FindOne(context, query).Decode(&theVersion)
	if err != nil {
		return err.Error()
	}
	return theVersion.Version
}

/**
* Describe this Store as a ConfigStore
 */
func (this *MongoStore) AsStoreItem() *config.StoreItem {
	var storeItem config.StoreItem
	storeItem = config.StoreItem{
		CollectionName: this.CollectionName,
		Version:        this.Version,
	}
	return &storeItem
}
