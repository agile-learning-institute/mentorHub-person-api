package stores

import (
	"mentorhub-person-api/src/config"
	"mentorhub-person-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	CollectionName string `json:"name"`
	Version        string `json:"value"`
	DefaultQuery   bson.M `json:"defaultQuery"`
	FilterQuery    bson.M `json:"filterQuery"`
	config         *config.Config
	collection     mongo.Collection
}

/**
* Constants for common Bson Query and Projection values
 */
func CollectionDefaultQuery() bson.M {
	return bson.M{"$and": []bson.M{
		{"name": bson.M{"$ne": "VERSION"}},
		{"status": bson.M{"$ne": "Archived"}},
	}}
}

func MongoShortNameProjection() bson.D {
	return bson.D{{Key: "ID", Value: "$_id"}, {Key: "name", Value: 1}}
}

/**
* Construct a PersonStore to handle person database io
 */
func NewMongoStore(cfg *config.Config, collectionName string, query bson.M) *MongoStore {
	store := &MongoStore{}

	// Initilize Store
	store.config = cfg
	store.CollectionName = collectionName
	store.collection = cfg.GetCollection(collectionName)
	store.Version = store.GetVersion()
	store.FilterQuery = query
	if query != nil {
		store.DefaultQuery = bson.M{"$and": []bson.M{
			CollectionDefaultQuery(),
			query,
		}}
	} else {
		store.DefaultQuery = CollectionDefaultQuery()
	}

	// Put the database Version in the Config
	store.config.AddConfigStore(store.AsStoreItem())
	return store
}

/**
* Simple wrapper for mongo InsertOne
 */
func (store *MongoStore) InsertOne(insertValues bson.M) (*mongo.InsertOneResult, error) {
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	return store.collection.InsertOne(context, insertValues)
}

/**
* Simple wrapper for mongo Find One
 */
func (store *MongoStore) FindOne(query bson.M) *mongo.SingleResult {
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	return store.collection.FindOne(context, query)
}

/**
* Simple wrapper for mongo Find Many
 */
func (store *MongoStore) FindMany(query bson.M, options *options.FindOptions) (*mongo.Cursor, error) {
	fullQuery := bson.M{"$and": []bson.M{
		store.DefaultQuery,
		query,
	}}

	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	return store.collection.Find(context, fullQuery, options)
}

/**
* Simple wrapper for mongo Find One and Update
 */
func (store *MongoStore) FindOneAndUpdate(query bson.M, update bson.M, options *options.FindOneAndUpdateOptions) *mongo.SingleResult {
	ctx, cancel := store.config.GetTimeoutContext()
	defer cancel()
	return store.collection.FindOneAndUpdate(ctx, query, update, options)
}

/**
* Get the collection schema version
 */
func (store *MongoStore) GetVersion() string {
	var theVersion models.VersionInfo
	var err error

	query := bson.M{"name": "VERSION"}
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	err = store.collection.FindOne(context, query).Decode(&theVersion)
	if err != nil {
		return err.Error()
	}
	return theVersion.Version
}

/**
* Describe this Store as a ConfigStore
 */
func (store *MongoStore) AsStoreItem() *config.StoreItem {
	var storeItem config.StoreItem = config.StoreItem{
		CollectionName: store.CollectionName,
		Version:        store.Version,
		Filter:         store.FilterQuery,
	}
	return &storeItem
}

/**
* Default Query by ID
 */
func (store *MongoStore) FindId(id string) (*map[string]interface{}, error) {
	// get the bson ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}

	var result map[string]interface{}
	var err error = store.FindOne(query).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/**
* Default FindMany Query
 */
func (store *MongoStore) FindDocuments() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	var err error

	// Query the database
	opts := options.Find()
	cursor, err := store.FindMany(store.DefaultQuery, opts)
	if err != nil {
		return nil, err
	}

	// Fetch all the results
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	err = cursor.All(context, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

/**
* Default Find Names Query
 */
func (store *MongoStore) FindNames() ([]models.ShortName, error) {
	var results []models.ShortName
	var err error

	// Query the database
	opts := options.Find().SetProjection(MongoShortNameProjection())
	cursor, err := store.FindMany(store.DefaultQuery, opts)
	if err != nil {
		return nil, err
	}

	// Fetch all the results
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	err = cursor.All(context, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
