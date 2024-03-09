package stores

import (
	"encoding/json"

	"mentorhub-person-api/src/config"
	"mentorhub-person-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonStore struct {
	config *config.Config
}

/**
* Construct a PersonStore to handle person database io
 */
func NewPersonStore(cfg *config.Config) *PersonStore {
	store := &PersonStore{}
	store.config = cfg
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
func (store *PersonStore) Insert(information []byte, crumb *models.BreadCrumb) (*models.Person, error) {
	// Get the document values
	var insertValues bson.M
	err := json.Unmarshal(information, &insertValues)
	if err != nil {
		return nil, err
	}

	// Addres OID values
	ConvertToOid(insertValues, "mentorId")
	ConvertToOid(insertValues, "partnerId")

	// Add the breadcrumb
	insertValues["lastSaved"] = crumb

	// Insert the document
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	result, err := store.config.GetPersonCollection().InsertOne(context, insertValues)
	if err != nil {
		return nil, err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()

	return store.FindId(id)
}

func (store *PersonStore) FindId(id string) (*models.Person, error) {
	// get the bson ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}

	var result *models.Person
	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	err := store.config.GetPersonCollection().FindOne(context, query).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/**
* Find One person and Update with the data provided
 */
func (store *PersonStore) FindOneAndUpdate(id string, request []byte, crumb *models.BreadCrumb) (*models.Person, error) {
	var thePerson models.Person

	// Build the query on ID
	objectID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": objectID}

	// Create the update set values
	var updateValues bson.M
	err := json.Unmarshal(request, &updateValues)
	if err != nil {
		return nil, err
	}

	ConvertToOid(updateValues, "mentorId")
	ConvertToOid(updateValues, "partnerId")

	// add breadcrumb to update object
	updateValues["lastSaved"] = crumb.AsBson()

	// Create the update object
	update := bson.M{"$set": updateValues}

	// Set Options
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	// Update the document
	ctx, cancel := store.config.GetTimeoutContext()
	defer cancel()
	err = store.config.GetPersonCollection().FindOneAndUpdate(ctx, query, update, options).Decode(&thePerson)
	if err != nil {
		return nil, err
	}

	return &thePerson, nil
}

/**
* Find Names
 */
func (store *PersonStore) FindNames(query bson.M) ([]models.ShortName, error) {
	var results []models.ShortName
	var err error

	// Query the database
	mentorProjection := bson.D{
		{Key: "ID", Value: "_id"},
		{Key: "name", Value: bson.M{"$concat": bson.A{"$firstName", " ", "$lastName"}}},
	}
	opts := options.Find().
		SetProjection(mentorProjection).
		SetSort(bson.D{{Key: "name", Value: 1}})

	fullQuery := bson.M{"$and": []bson.M{
		{"name": bson.M{"$ne": "VERSION"}},
		{"status": bson.M{"$ne": "Archived"}},
		query,
	}}

	context, cancel := store.config.GetTimeoutContext()
	defer cancel()
	cursor, err := store.config.GetPersonCollection().Find(context, fullQuery, opts)
	if err != nil {
		return nil, err
	}

	// Fetch all the results
	context, cancel = store.config.GetTimeoutContext()
	defer cancel()
	err = cursor.All(context, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
