/********************************************************************************
** MongoIO
**    This class is a wrapper for MongoDB interactions. It manages
**    the database connection, and executes all mongoDb calls with
**    easy to mock wrappers.
********************************************************************************/
package config

import (
	"go.mongodb.org/mongo-driver/bson"
)

type MongoIO struct {
	// config             Config
	// client             *mongo.Client
	// database           *mongo.Database
	// peopleCollection   *mongo.Collection
	// enumsCollection    *mongo.Collection
	// partnersCollection *mongo.Collection
	// versionsCollection *mongo.Collection
	// cancel             context.CancelFunc
}

// Common MongoDB Query Objects
func GetAllQuery() bson.M {
	return bson.M{"$and": []bson.M{
		{"name": bson.M{"$ne": "VERSION"}},
		{"status": bson.M{"$ne": "Archived"}},
	}}
}

func GetMentorsQuery() bson.M {
	return bson.M{"$and": []bson.M{
		{"status": bson.M{"$ne": "Archived"}},
		{"roles": "Mentor"},
	}}
}

func GetMentorsProjection() bson.D {
	return bson.D{
		{Key: "ID", Value: "_id"},
		{Key: "name", Value: bson.M{"$concat": bson.A{"$firstName", " ", "$lastName"}}},
	}
}

/**********************************************************************
* Constructor - initilize configuration values
 */
func NewMongoIO() *MongoIO {
	return &MongoIO{}
}

// /********************************************************************************
// * Connect to the Database, load the Versions and Enumerators Config values
//  */
// func (mongo *MongoIO) Connect() {
// 	// Connect to the database
// 	ctx, cancel := mongo.GetTimeoutContext()
// 	mongo.cancel = cancel
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.getConnectionString()))
// 	if err != nil {
// 		cancel()
// 		log.Fatal("Database Connection Failed:", err)
// 	}

// 	// Get the database object, and collection objects
// 	cfg.client = client
// 	cfg.database = cfg.client.Database(cfg.databaseName)
// 	cfg.peopleCollection = cfg.database.Collection(PeoplesCollectionName)
// 	cfg.enumsCollection = cfg.database.Collection(EnumeratorsCollectionName)
// 	cfg.partnersCollection = cfg.database.Collection(PartnersCollectionName)
// 	cfg.versionsCollection = cfg.database.Collection(VersionsCollectionName)

// 	// Find Versions
// 	query := bson.M{}
// 	opts := options.Find()
// 	cfg.FindDocuments(cfg.versionsCollection, query, opts, &cfg.Versions)
// 	if len(cfg.Versions) == 0 {
// 		log.Fatal("Versions not found")
// 	}

// 	// Find the current version for the people collection
// 	versionString := ""
// 	for _, v := range cfg.Versions {
// 		if v.CollectionName == PeoplesCollectionName {
// 			versionString = v.CurrentVersion
// 		}
// 	}
// 	if versionString == "" {
// 		log.Fatalf("People Collection not found %d", len(cfg.Versions))
// 	}

// 	// Extract the enumerators version number from the version string
// 	var versionNumber int
// 	parts := strings.Split(versionString, ".")
// 	versionNumber, err = strconv.Atoi(parts[len(parts)-1])
// 	if err != nil {
// 		log.Fatalf("Version Number not an integer: %s", versionString)
// 	}

// 	// Query Enumerators
// 	results := []*Enumerators{}
// 	query = bson.M{"version": versionNumber}
// 	opts = options.Find()
// 	opts.SetSort(bson.D{{Key: "version", Value: 1}})
// 	cfg.FindDocuments(cfg.enumsCollection, query, opts, &results)

// 	// Set enumerators
// 	if len(results) == 1 {
// 		cfg.Enumerators = results[0].Enumerators
// 	} else {
// 		log.Fatalf("Enumerators List is wrong size %d", len(results))
// 	}
// }

// func (cfg *Config) FindDocuments(collection *mongo.Collection, query bson.M, opts *options.FindOptions, results interface{}) {
// 	context, cancel := cfg.GetTimeoutContext()
// 	defer cancel()

// 	cursor, err := collection.Find(context, query, opts)
// 	if err != nil {
// 		log.Fatalf("Query Failed: %s %v", collection.Name(), err)
// 	}
// 	defer cursor.Close(context)

// 	// Directly pass results to cursor.All
// 	if err = cursor.All(context, results); err != nil {
// 		log.Fatalf("Fetch Failed: %s %v", collection.Name(), err)
// 	}
// }

// /********************************************************************************
// * Disconnect fromthe Database
//  */
// func (cfg *Config) Disconnect() {
// 	ctx, cancel := cfg.GetTimeoutContext()
// 	defer cancel()
// 	cfg.client.Disconnect(ctx)
// 	cfg.cancel()
// }

// /********************************************************************************
// * Simple Getters
//  */
// func (cfg *Config) GetPort() string {
// 	return cfg.port
// }

// func (cfg *Config) GetPersonCollection() *mongo.Collection {
// 	return cfg.peopleCollection
// }

// func (cfg *Config) GetTimeoutContext() (context.Context, context.CancelFunc) {
// 	timeout := time.Duration(cfg.databaseTimeout) * time.Second
// 	return context.WithTimeout(context.Background(), timeout)
// }

// /********************************************************************************
// * Load the Mentors and Partners list
//  */
// func (cfg *Config) LoadLists() error {

// 	// Fetch Mentors
// 	mentors, err := cfg.findNames(
// 		cfg.peopleCollection,
// 		options.Find().SetProjection(GetMentorsProjection()),
// 		GetMentorsQuery(),
// 	)

// 	if err != nil {
// 		log.Printf("ERROR: Load Mentors failed %s", err)
// 		return err
// 	}
// 	cfg.Mentors = mentors

// 	// Fetch Partners
// 	partners, err := cfg.findNames(
// 		cfg.partnersCollection,
// 		options.Find(),
// 		GetAllQuery(),
// 	)

// 	if err != nil {
// 		log.Printf("ERROR: Load Partners failed %s", err)
// 		return err
// 	}
// 	cfg.Partners = partners

// 	return nil
// }

// /**
// * Get a list of people names based on the query and options provided
//  */
// func (cfg *Config) findNames(collection *mongo.Collection, opts *options.FindOptions, query bson.M) ([]*models.ShortName, error) {
// 	var results []*models.ShortName
// 	var err error

// 	// Query the database
// 	context, cancel := cfg.GetTimeoutContext()
// 	defer cancel()
// 	cursor, err := collection.Find(context, query, opts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Fetch all the results
// 	context, cancel = cfg.GetTimeoutContext()
// 	defer cancel()
// 	err = cursor.All(context, &results)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return results, nil
// }
