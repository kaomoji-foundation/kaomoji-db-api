package models

import (
	"context"
	"log"

	"GO-API-template/src/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Role struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// Role string
	Role string `bson:"role" json:"role"`
	// Integer to indicate the permissons level, the lower, the higer perms.
	// 0 is full acess, this is only used to compare roles
	Level int `bson:"level" json:"level"`
	// Permissons object, this is used when performing a specific action
	Permissons Permissons `bson:"permissons" json:"permissons"`
}

// Theese permissons are checked every time a user wants to make an operation restricted to certain roles
type Permissons struct { //? shuld have a FullAcess flag?
	ReadUsers  bool `bson:"readUsers" json:"readUsers"`
	UsersAdmin bool `bson:"usersAdmin" json:"usersAdmin"`
	ReadRoles  bool `bson:"readRoles" json:"readRoles"`
	RolesAdmin bool `bson:"rolesAdmin" json:"rolesAdmin"`
}

//? The plurificated interfaces of the models are probably useless AAAND anoying
type Roles interface {
	Create() *mongo.InsertOneResult
	ReadAll() []User
	CreateSingletonDBAndCollection()
}

var RolesCollection *mongo.Collection

// this is the database where the collection is expected, could have multiple if necessary
var rolesModelDB *mongo.Database

func (r Role) CreateSingletonDBAndCollection() {
	if rolesModelDB == nil {
		rolesModelDB = services.Mongo.DBs["mainDB"]
	}
	if RolesCollection == nil {
		RolesCollection = rolesModelDB.Collection("roles")
	}
}

func (r Role) Create() (*mongo.InsertOneResult, error) {
	r.CreateSingletonDBAndCollection()

	insertedRes, err := RolesCollection.InsertOne(context.Background(), r)
	if err != nil {
		log.Println(err)
	}

	return insertedRes, err
}

func (r Role) ReadAll() []User {
	r.CreateSingletonDBAndCollection()

	filter := bson.D{
		//{"name", "uwu"}, // this works, Try it!
	}

	currsor, err := RolesCollection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer currsor.Close(services.Mongo.Context)

	var users []User
	currsor.All(context.Background(), &users)

	return users
}

func (r *Role) Fill(identity string, id, role bool) error {
	var fields = []bson.M{}
	if id {
		id, err := primitive.ObjectIDFromHex(identity)
		if err == nil {
			fields = append(fields, bson.M{"_id": id})
		}
	}

	if role {
		fields = append(fields, bson.M{"role": identity})
	}

	filter := bson.M{"$or": fields}

	res := RolesCollection.FindOne(context.Background(), filter)
	if err := res.Err(); err != nil {
		return err
	}
	res.Decode(r)
	return nil
}
