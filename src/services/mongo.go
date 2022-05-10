package services

import (
	"context"
	"log"
	"time"

	cfg "GO-API-template/src/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoServiceData struct {
	Client        *mongo.Client
	Context       context.Context
	CancelContext context.CancelFunc
	DBs           map[string]*mongo.Database
	createdClient bool
}

var Mongo MongoServiceData

// Creates a singleton mongo client instance and sets it as public vars for the package
func (m *MongoServiceData) getMongoClient() (*mongo.Client, *context.Context, *context.CancelFunc) {
	if m.createdClient {
		return Mongo.Client, &(Mongo.Context), &(Mongo.CancelContext)
	}
	ctx, cancelContext := context.WithTimeout(context.Background(), 50*time.Second)

	m.Context = ctx
	m.CancelContext = cancelContext

	// MongoDB connection uri and credentials
	uri := cfg.Config.Mongo.URI
	user := cfg.Config.Mongo.User
	passwd := cfg.Config.Mongo.Pass

	// login will still be done on admin database
	//TODO: env var to attempt login on specific db or db set to
	client, err := mongo.Connect(
		Mongo.Context,
		options.Client().ApplyURI(
			uri,
		).SetAuth(
			options.Credential{
				Username: user,
				Password: passwd,
			},
		))
	// check for connection errors
	if err != nil {
		log.Println(err.Error())
	}
	m.Client = client
	m.createdClient = true

	return Mongo.Client, &(Mongo.Context), &(Mongo.CancelContext)
}

//allow for multi database (usefull for mocking and testing)
func (m *MongoServiceData) initDBs() *map[string]*mongo.Database {
	m.DBs = map[string]*mongo.Database{
		//* here the databases in use will be stored, usually just one
		"mainDB": m.Client.Database(cfg.Config.Mongo.DBs[0]),
	}
	return &m.DBs
}

func (m *MongoServiceData) Init() *context.CancelFunc {
	_, _, cancelCtx := m.getMongoClient()

	// Connect
	m.Client.Connect(m.Context)
	log.Println("[DB Connection]: Created")
	// test connection
	log.Println("[DB Connection]: Pinging")
	err := m.Client.Ping(m.Context, nil)
	if err != nil {
		//? Should it panic the app so it fails to start if failed to connect to the db
		log.Println(err)
		log.Println("[DB Connection]: Could NOT be established")
	} else {
		log.Println("[DB Connection]: OK")
	}

	m.initDBs()

	return cancelCtx
}
