// Copyright (c) Luke Atherton 2015

package authenticator

import (
	"fmt"
	"log"
	"time"

	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Repo interface {
	SaveCredentials(userId uuid.UUID, credentials *Credentials) (err error)
	GetCredentials(userId uuid.UUID) (credentials *Credentials, err error)

	FindEmail(email string) (id uuid.UUID, err error)
	Cleanup()
}

const (
	MongoDBHosts = "127.0.0.1:27000"
	AuthDatabase = "admin"
	AuthUserName = "admin"
	AuthPassword = "test"
	AppDatabase  = "authentication-prod"
	TestDatabase = "authentication-test"
)

type MongoDBRepo struct {
	db *mgo.Session
}

func NewMongoRepo(dbHosts []string, authDb string, dbUsername string, dbPassword string) Repo {

	fmt.Printf("Connecting to MongoDB cluster: %s\n", dbHosts)

	// if authDb != "" {
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    dbHosts,
		Timeout:  60 * time.Second,
		Database: authDb,
		Username: dbUsername,
		Password: dbPassword,
	}

	// } else {
	// 	mongoSession, err := mgo.Dial(dbHosts[0])
	// }

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	mongoSession.SetMode(mgo.Monotonic, true)

	credentialsCollection := mongoSession.DB(AppDatabase).C("credentials")

	// Index
	idx_email := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = credentialsCollection.EnsureIndex(idx_email)
	if err != nil {
		panic(err)
	}

	repo := &MongoDBRepo{
		db: mongoSession,
	}

	fmt.Println("Init app db done")

	return repo
}

func NewMongoTestRepo(dbHosts []string, authDb string, dbUsername string, dbPassword string) Repo {
	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward, the data read will be consistent
	// across sequential queries in the same session, and modifications made
	// within the session will be observed in following queries (read-your-writes).
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode
	mongoSession.SetMode(mgo.Monotonic, true)

	err = mongoSession.DB(TestDatabase).DropDatabase()
	if err != nil {
		panic(err)
	}

	credentialsCollection := mongoSession.DB(TestDatabase).C("credentials")

	// Index
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = credentialsCollection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	repo := &MongoDBRepo{
		db: mongoSession,
	}

	fmt.Println("Init test db done")

	return repo
}

func (repo *MongoDBRepo) Cleanup() {
	socketConnection := repo.db.Copy()
	defer socketConnection.Close()

	err := socketConnection.DB(TestDatabase).DropDatabase()

	if err != nil {
		panic(err)
	}

	fmt.Println("Database Cleared")
}

func (repo *MongoDBRepo) SaveCredentials(userId uuid.UUID, credentials *Credentials) (err error) {

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	socketConnection := repo.db.Copy()
	defer socketConnection.Close()

	// Get a collection to execute the query against.
	collection := socketConnection.DB(TestDatabase).C("credentials")

	_, err = collection.Upsert(bson.M{"id": userId}, credentials)

	return err
}

func (repo *MongoDBRepo) GetCredentials(userId uuid.UUID) (credentials *Credentials, err error) {

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	socketConnection := repo.db.Copy()
	defer socketConnection.Close()

	// Get a collection to execute the query against.
	collection := socketConnection.DB(TestDatabase).C("credentials")

	result := &Credentials{}
	err = collection.Find(bson.M{"id": userId}).One(&result)

	return result, err
}

func (repo *MongoDBRepo) FindEmail(email string) (id uuid.UUID, err error) {
	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	socketConnection := repo.db.Copy()
	defer socketConnection.Close()

	// Get a collection to execute the query against.
	collection := socketConnection.DB(TestDatabase).C("credentials")

	result := Credentials{}
	err = collection.Find(bson.M{"email": email}).One(&result)

	return result.Id, err
}
