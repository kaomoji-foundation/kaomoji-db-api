package models

import (
	"context"
	"log"
	"reflect"
	"time"

	"Kaomoji-DB/src/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Issue struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Author      UserMinimal        `bson:"author,omitempty" json:"author"`
	Reviewer    UserMinimal        `bson:"reviewer,omitempty" json:"reviewer"`
	StartDate   time.Time          `bson:"startDate,omitempty" json:"startDate"`
	EndDate     time.Time          `bson:"endDate,omitempty" json:"endDate"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	// Status: "acepted", "closed", "on review", "assigned", "pending"
	Status string `bson:"status" json:"status"`
	// Operation: "create", "modify", "delete", "report"
	Operation string  `bson:"operation" json:"operation"`
	Kaomoji   Kaomoji `bson:"kaomoji" json:"kaomoji"`
	Solution  Kaomoji `bson:"solution" json:"solution"`
}

// basic reduced issuedata viewable by anyone to use on issues listings
type IssueMinimal struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Author   UserMinimal        `bson:"author,omitempty" json:"author"`
	Reviewer UserMinimal        `bson:"reviewer,omitempty" json:"reviewer"`
	Title    string             `bson:"title" json:"title"`
	// Status: "acepted", "closed", "on review", "assigned", "pending"
	Status string `bson:"status" json:"status"`
	// Operation: "create", "modify", "delete", "report"
	Operation string         `bson:"operation" json:"operation"`
	Kaomoji   KaomojiMinimal `bson:"kaomoji" json:"kaomoji"`
}

//? The plurificated interfaces of the models are probably useless AAAND anoying
type Issues interface {
	Create() *mongo.InsertOneResult
	ReadAll() []Issue
	CreateSingletonDBAndCollection()
}

var IssuesCollection *mongo.Collection

// this is the database where the collection is expected, could have multiple if necessary
var issueModelDB *mongo.Database

// Returns the public/viewable issue info
func (i Issue) Minimal() IssueMinimal {
	return IssueMinimal{
		ID:        i.ID,
		Author:    i.Author,
		Reviewer:  i.Reviewer,
		Status:    i.Status,
		Operation: i.Operation,
		Kaomoji:   i.Kaomoji.Minimal(),
	}
}

func (i Issue) CreateSingletonDBAndCollection() {
	if issueModelDB == nil {
		issueModelDB = services.Mongo.DBs["mainDB"]
	}
	if IssuesCollection == nil {
		IssuesCollection = issueModelDB.Collection("issues")
	}
}

func (i *Issue) Create() (*mongo.InsertOneResult, error) {
	i.CreateSingletonDBAndCollection()

	insertedRes, err := IssuesCollection.InsertOne(context.Background(), i)
	if err != nil {
		log.Println(err)
	}
	i.ID = insertedRes.InsertedID.(primitive.ObjectID)

	return insertedRes, err
}

func (i *Issue) Update() error {
	_, err := IssuesCollection.UpdateByID(context.Background(), i.ID, i)
	return err
}

func (i Issue) ReadAll() []Issue {
	i.CreateSingletonDBAndCollection()

	filter := bson.D{
		//{"name", "uwu"}, // this works, Try it!
	}

	currsor, err := IssuesCollection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer currsor.Close(services.Mongo.Context)

	var issues []Issue
	currsor.All(context.Background(), &issues)

	return issues
}

// fills the issue checking id
func (i *Issue) Fill(idHex string) error {

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	res := IssuesCollection.FindOne(context.Background(), filter)
	if err := res.Err(); err != nil {
		return err
	}
	err = res.Decode(i)
	return err
}

// Sets the status to "acepted" and adds and apply the revision to the propper kaomoji
func (i *Issue) Accept() error {
	// i.Kaomoji should only contain ID and the fields to update

	var current Kaomoji
	err := current.Fill(i.Kaomoji.ID.Hex(), true, false)
	if err != nil {
		return err
	}

	i.Status = "acepted"
	err = i.Update()
	if err != nil {
		return err
	}

	// Do nothing to the Kaomoji if if is already up to date OR if is empty
	if reflect.DeepEqual(current, i.Solution) || reflect.DeepEqual(i.Solution, Kaomoji{}) {
		return nil
	}

	err = i.Solution.Update()

	return err
}
