/********************************************************************************
** Person Store
**    This class implementes the mongo calls needed by Person Handlers
********************************************************************************/
package stores

import (
	"encoding/json"
	"log"

	"mentorhub-person-api/src/config"
	"mentorhub-person-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonStore struct {
	mongoIO config.MongoIOInterface
	person  *mongo.Collection
}

type PersonStoreInterface interface {
	Insert(information []byte, crumb *models.BreadCrumb) (string, error)
	FindId(id string) (*models.Person, error)
	UpdateId(id string, request []byte, crumb *models.BreadCrumb) (*models.Person, error)
	FindNames() ([]*config.ShortName, error)
}

var _ PersonStoreInterface = (*PersonStore)(nil)

/**
* Construct a PersonStore to handle person database io
 */
func NewPersonStore(io config.MongoIOInterface) *PersonStore {
	store := &PersonStore{}
	store.mongoIO = io
	store.person = io.GetPeopleCollection()
	if store.person == nil {
		log.Fatal("MongoIO GetPeopleCollection returns nil")
	}
	return store
}

/**
* Convert Id string to ObjectId if present
 */
func ConvertToOid(values bson.M, fieldName string) {
	if _, idExists := values[fieldName]; idExists {
		idValue := values[fieldName].(string)
		newObjectID, _ := primitive.ObjectIDFromHex(idValue)
		values[fieldName] = newObjectID
	}
}

/**
* Insert a new person with the information provided
 */
func (store *PersonStore) Insert(information []byte, crumb *models.BreadCrumb) (string, error) {
	// Get the document values
	var insertValues bson.M
	err := json.Unmarshal(information, &insertValues)
	if err != nil {
		return "", err
	}

	// Addres OID values
	ConvertToOid(insertValues, "mentorId")
	ConvertToOid(insertValues, "partnerId")

	// Add the breadcrumb
	insertValues["lastSaved"] = crumb.AsBson()

	// Insert the document
	result, err := store.mongoIO.InsertOne(store.person, insertValues)
	if err != nil {
		return "", err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (store *PersonStore) FindId(id string) (*models.Person, error) {
	// get the bson ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	result := models.Person{}
	err := store.mongoIO.FindOne(store.person, query, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/**
* Find One person and Update with the data provided
 */
func (store *PersonStore) UpdateId(id string, request []byte, crumb *models.BreadCrumb) (*models.Person, error) {

	// Create the update set values
	var updateValues bson.M
	err := json.Unmarshal(request, &updateValues)
	if err != nil {
		return nil, err
	}

	// Conviert strings to OID where needed
	ConvertToOid(updateValues, "mentorId")
	ConvertToOid(updateValues, "partnerId")

	// add breadcrumb to update object
	updateValues["lastSaved"] = crumb.AsBson()

	// Setup the Update
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}
	update := bson.M{"$set": updateValues}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	resposne := &models.Person{}

	// Update the document
	err = store.mongoIO.UpdateOne(store.person, query, options, update, resposne)
	if err != nil {
		return nil, err
	}
	return resposne, nil
}

/**
* Find Names
 */
func (store *PersonStore) FindNames() ([]*config.ShortName, error) {
	// Setup Projections and Sort
	sortOrder := bson.D{
		{Key: "firstName", Value: 1},
		{Key: "lastName", Value: 1},
	}

	// Setup the Query
	var results []*config.ShortName
	query := bson.M{"status": bson.M{"$ne": "Archived"}}
	opts := options.Find()
	opts.SetProjection(config.NameProjection())
	opts.SetSort(sortOrder)

	// Find the documents
	err := store.mongoIO.Find(store.person, query, opts, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
