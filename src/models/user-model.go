package models

import (
	"context"
	"log"
	"sync"

	"kaomojidb/src/config"
	"kaomojidb/src/services"

	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	Name     string             `bson:"name" json:"name"`
	Role     string             `bson:"role" json:"role"`
	RoleID   primitive.ObjectID `bson:"roleID" json:"roleID"`
	// Tokens list is only used to be able to block the token later by placing said token onto the BlockedTokens list
	Tokens map[string]bool `bson:"tokens" json:"tokens"`
	// Any attempt to use the tokens stored here, will be blocked
	BlockedTokens map[string]bool `bson:"blockedTokens" json:"blockedTokens"`
}

// ? The plurificated interfaces of the models are probably useless AAAND anoying
type Users interface {
	Create() *mongo.InsertOneResult
	ReadAll() []User
	CreateSingletonDBAndCollection()
}

// basic reduced userdata viewable by anyone to use on users listings
type UserMinimal struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Role     string             `bson:"role" json:"role"`
}

// userdata viewable by anyone
type userPublic struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Role     string             `bson:"role" json:"role"`
}

// User data only shown to the user
type userPrivate struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Name     string             `bson:"name" json:"name"`
	Role     string             `bson:"role" json:"role"`
	//RoleID   primitive.ObjectID `bson:"roleID" json:"roleID"`
	// Tokens list is only used to be able to block the token later by placing said token onto the BlockedTokens list
	Tokens map[string]bool `bson:"tokens" json:"tokens"`
	// Any attempt to use the tokens stored here, will be blocked
	BlockedTokens map[string]bool `bson:"blockedTokens" json:"blockedTokens"`
}

/*
User tokens struct, only used for token pruning
this is to avoid a race condition with the rest of the data set for the user
*/
type userTokens struct {
	// Tokens list is only used to be able to block the token later by placing said token onto the BlockedTokens list
	Tokens map[string]bool `bson:"tokens" json:"tokens"`
	// Any attempt to use the tokens stored here, will be blocked
	BlockedTokens map[string]bool `bson:"blockedTokens" json:"blockedTokens"`
}

var UsersCollection *mongo.Collection

// this is the database where the collection is expected, could have multiple if necessary
var userModelDB *mongo.Database

// Returns the public/viewable user info
func (u User) Minimal() userPublic {
	return userPublic{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
}

// Returns the public/viewable user info
func (u User) Public() userPublic {
	return userPublic{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
}

// Returns the user viewable info
func (u User) Private() userPrivate {
	return userPrivate{
		ID:            u.ID,
		Username:      u.Username,
		Email:         u.Email,
		Name:          u.Name,
		Role:          u.Role,
		Tokens:        u.Tokens,
		BlockedTokens: u.BlockedTokens,
	}
}

func (u User) CreateSingletonDBAndCollection() {
	if userModelDB == nil {
		userModelDB = services.Mongo.DBs["mainDB"]
	}
	if UsersCollection == nil {
		UsersCollection = userModelDB.Collection("users")
	}
}

func (u *User) Create() (*mongo.InsertOneResult, error) {
	u.CreateSingletonDBAndCollection()

	insertedRes, err := UsersCollection.InsertOne(context.Background(), u)
	if err != nil {
		log.Println(err)
	}
	u.ID = insertedRes.InsertedID.(primitive.ObjectID)

	return insertedRes, err
}

func (u User) ReadAll() []User {
	u.CreateSingletonDBAndCollection()

	filter := bson.D{
		//{"name", "uwu"}, // this works, Try it!
	}

	currsor, err := UsersCollection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer currsor.Close(services.Mongo.Context)

	var users []User
	currsor.All(context.Background(), &users)

	return users
}

// fills the user checking  id, username or email, multiple of them can be used at the same time
func (u *User) Fill(identity string, id, username, email bool) error {
	var fields = []bson.M{}
	if id {
		id, err := primitive.ObjectIDFromHex(identity)
		if err == nil {
			fields = append(fields, bson.M{"_id": id})
		}
	}

	if username {
		fields = append(fields, bson.M{"username": identity})
	}
	if email {
		fields = append(fields, bson.M{"email": identity})
	}

	filter := bson.M{"$or": fields}

	res := UsersCollection.FindOne(context.Background(), filter)
	if err := res.Err(); err != nil {
		return err
	}
	res.Decode(u)
	return nil
}

// Will only return false if it finds a diferent user wit the same username or email
func (u User) CheckUnique() (bool, error) {
	fields := []bson.M{
		{"username": u.Username},
	}
	if u.Email != "" {
		fields = append(fields, bson.M{"email": u.Email})
	}
	filter := bson.M{
		"_id": bson.M{"$ne": u.ID},
		"$or": fields,
	}

	count, err := UsersCollection.CountDocuments(context.Background(), filter)

	return count == 0, err
}

// Returns the specific data about the role for the user from the DB
func (u User) RoleData() (Role, error) {
	filter := bson.M{"_id": u.RoleID}

	res := RolesCollection.FindOne(context.Background(), filter)

	var role Role
	err := res.Decode(&role)

	if err != nil {
		return role, err
	}

	return role, nil
}

// Sets the role id by searching it via the name stored on user.Role
func (u *User) SetRole() error {
	filter := bson.M{"role": u.Role}
	res := RolesCollection.FindOne(context.Background(), filter)
	if err := res.Err(); err != nil {
		return err
	}
	var role Role
	err := res.Decode(&role)

	if err == nil {
		u.RoleID = role.ID
	}
	return err
}

// * Requires to have a well formed ID
func (u *User) LoadTokens() error {

	res := UsersCollection.FindOne(
		context.Background(),
		bson.M{
			"_id": u.ID,
		},
		options.MergeFindOneOptions().SetProjection(bson.M{"tokens": 1, "blockedTokens": 1}),
	)
	if res.Err() != nil {
		return res.Err()
	}
	err := res.Decode(&u)
	return err // returns nil or the last error
}

// * Requires data to be filled first, will push to database.
func (u *User) PruneTokens() {
	var wg sync.WaitGroup
	tokenLists := []*map[string]bool{&u.BlockedTokens, &u.Tokens}
	for _, list := range tokenLists {
		wg.Add(1)
		go pruneTokenList(list, &wg)
	}
	wg.Wait()
	// update the tokens lists on the db
	var userTokens userTokens
	userTokens.BlockedTokens = u.BlockedTokens
	userTokens.Tokens = u.Tokens

	filter := bson.M{"_id": u.ID}
	update := bson.M{"$set": userTokens}
	_, err := UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("[TokenPruneRoutine]:" + err.Error())
	} else {
		log.Println("[TokenPruneRoutine]: Finished single user routine sucessfully")
	}
}

// checks if the tokens in a list are expired and deletes them in that case.
func pruneTokenList(list *map[string]bool, callerWGs ...*sync.WaitGroup) {

	var wg sync.WaitGroup
	for token := range *list {
		//set up function to call asyncronously
		pruneFromList := func(token string, wg *sync.WaitGroup) {
			_, expired := tokenIsInvalid(token)
			if expired {
				//log.Println("[TokenPruneRoutine]: Found expired token: ")
				delete(*list, token)
				log.Printf("[TokenPruneRoutine]: %v", len(*list))
			}
			wg.Done()
		}
		wg.Add(1)
		go pruneFromList(token, &wg)
	}

	wg.Wait()
	if len(callerWGs) == 1 {
		// release one from the wait group of the caller function
		callerWGs[0].Done()
	}
}

// Parses a token and checks if token is valid, retuns parsed Token and invalid tag
func tokenIsInvalid(token string) (*jwt.Token, bool) {
	tok, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JWT.Secret), nil
	})

	// check if the token has expired
	return tok, !tok.Valid
}

// * Requires to have a well formed ID
func (u *User) BlockToken() {

}
