package models

import (
	"context"
	"errors"
	"log"

	"GO-API-template/src/services"
	"GO-API-template/src/utils/radix"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Kaomoji struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	String     []rune             `bson:"string,omitempty" json:"string"`
	Desciption string             `bson:"description,omitempty" json:"description"`
	Categories []string           `bson:"categories,omitempty" json:"categories"`
	Revisions  []IssueMinimal     `bson:"revisions,omitempty" json:"revisions"`
}

//? The plurificated interfaces of the models are probably useless AAAND anoying
type Kaomojis interface {
	Create() *mongo.InsertOneResult
	ReadAll() []Kaomoji
	CreateSingletonDBAndCollection()
}

// basic reduced kaomojidata viewable by anyone to use on kaomojis listings
type KaomojiMinimal struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	String     []rune             `bson:"string,omitempty" json:"string"`
	Desciption string             `bson:"description,omitempty" json:"description"`
	Categories []string           `bson:"categories,omitempty" json:"categories"`
}

var KaomojisCollection *mongo.Collection

// this is the database where the collection is expected, could have multiple if necessary
var kaomojiModelDB *mongo.Database

const ErrPolicyCategoriesLength = "Error: ilegal Categories slice size"

// Returns the public/viewable kaomoji info
func (k Kaomoji) Minimal() KaomojiMinimal {
	return KaomojiMinimal{
		ID: k.ID,
	}
}

func (k Kaomoji) CreateSingletonDBAndCollection() {
	if kaomojiModelDB == nil {
		kaomojiModelDB = services.Mongo.DBs["mainDB"]
	}
	if KaomojisCollection == nil {
		KaomojisCollection = kaomojiModelDB.Collection("kaomojis")
	}
}

func (k *Kaomoji) Create() (*mongo.InsertOneResult, error) {
	k.CreateSingletonDBAndCollection()

	if len(k.Categories) > 10 {
		return nil, errors.New(ErrPolicyCategoriesLength)
	}
	radix.Sort(k.Categories)

	insertedRes, err := KaomojisCollection.InsertOne(context.Background(), k)
	if err != nil {
		log.Println(err)
	}
	k.ID = insertedRes.InsertedID.(primitive.ObjectID)

	return insertedRes, err
}

func (k Kaomoji) ReadAll() []Kaomoji {
	k.CreateSingletonDBAndCollection()

	filter := bson.D{
		//{"name", "uwu"}, // this works, Try it!
	}

	currsor, err := KaomojisCollection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer currsor.Close(services.Mongo.Context)

	var kaomojis []Kaomoji
	currsor.All(context.Background(), &kaomojis)

	return kaomojis
}

// fills the kaomoji checking  id, string or email, multiple of them can be used at the same time
func (k *Kaomoji) Fill(identificator string, id, _string bool) error {
	var fields = []bson.M{}
	if id {
		id, err := primitive.ObjectIDFromHex(identificator)
		if err == nil {
			fields = append(fields, bson.M{"_id": id})
		}
	}

	if _string {
		fields = append(fields, bson.M{"string": identificator})
	}

	filter := bson.M{"$or": fields}

	res := KaomojisCollection.FindOne(context.Background(), filter)
	if err := res.Err(); err != nil {
		return err
	}
	res.Decode(k)
	return nil
}

// Fills the kaomoji checking  id, string or email, multiple of them can be used at the same time
func (k *Kaomoji) FillMinimal(identificator string, id, _string bool) (KaomojiMinimal, error) {
	var fields = []bson.M{}
	if id {
		id, err := primitive.ObjectIDFromHex(identificator)
		if err == nil {
			fields = append(fields, bson.M{"_id": id})
		}
	}
	if _string {
		fields = append(fields, bson.M{"string": identificator})
	}

	var min KaomojiMinimal

	filter := bson.M{"$or": fields}
	projection := bson.M{
		"_id":         1,
		"string":      1,
		"description": 1,
		"categories":  1,
	}
	opts := &options.FindOneOptions{Projection: projection}
	res := KaomojisCollection.FindOne(context.Background(), filter, opts)
	if err := res.Err(); err != nil {
		return min, err
	}
	res.Decode(k)

	return k.Minimal(), nil
}

func (k *Kaomoji) Update() error {
	if len(k.Categories) > 10 {
		return errors.New(ErrPolicyCategoriesLength)
	}
	radix.Sort(k.Categories)
	_, err := IssuesCollection.UpdateByID(context.Background(), k.ID, k)
	return err
}

// Will only return false if it finds a diferent kaomoji with the same string
func (k Kaomoji) CheckUnique() (bool, error) {

	filter := bson.M{
		"_id":    bson.M{"$ne": k.ID},
		"string": k.String,
	}

	count, err := UsersCollection.CountDocuments(context.Background(), filter)

	return count == 0, err
}

//TODO: GetSimilar() ([]KaomojiMinimal, error) {...} get kaomojis with similar keys, or strings (using a threshold)
/*
Get kaomojis with similar keys, or strings, Only one of the requirements have to match tobe considered similar
 	- keys, the minimum amount (inclusive) of keys to match, 0 will completely skip this checking
 	- _string, the maximum fuzyness in characters in trimed search, 0 will skip checking
*/
func (k Kaomoji) GetSimilar(keysThreshold int, stringFuzyness int) ([]KaomojiMinimal, error) {
	filter := bson.M{}
	projection := bson.M{
		"_id":         1,
		"string":      1,
		"description": 1,
		"categories":  1,
	}
	opts := &options.FindOptions{Projection: projection}
	cursor, err := KaomojisCollection.Find(context.Background(), filter, opts)

	if err != nil {
		return nil, err
	}

	var res []KaomojiMinimal

	for cursor.Next(context.Background()) {
		var kao KaomojiMinimal
		err = cursor.Decode(&kao)
		if err != nil {
			return nil, err
		}

		similarStr := stringFuzyness != 0 && fuzzyStringIsIn(kao.String, k.String, stringFuzyness)
		if similarStr || matchedStringSlice(k.Categories, kao.Categories) >= keysThreshold {
			res = append(res, kao)
		}
	}

	return res, nil
}

//? should asume s1 is sorted and use binary search instead?
/*
Returns the amount of matched items in the two string slices in O(n1+n2)
*/
func matchedStringSlice(s1, s2 []string) int {
	count := 0
	set := make(map[string]bool)
	for _, v := range s1 {
		set[v] = true
	}
	for _, v := range s2 {
		if set[v] {
			count++
		}
	}
	return count
}

//TODO: fuzzyStringIsIn(stri, subject []rune, fuzyness ...int) bool {...}
/*
Trims str checking against subject in a fuzy way, and checks if subject is within str
str string to find subject in
subjet string to be found within str in a fuzy way
Fuzyness defaults to 2, so at most 2 operations[insertions, deletions or substitutions] on the trimed string
*/
func fuzzyStringIsIn(stri, subject []rune, fuzyness ...int) bool {

	return false
}

func levenshtein(str1, str2 []rune) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1

		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var inc int
			if str1[y-1] != str2[x-1] {
				inc = 1
			}
			column[y] = minimum(column[y]+1, column[y-1]+1, lastkey+inc)
			lastkey = oldkey
		}
	}
	return column[s1len]
}
func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else if b < c {
		return b
	}
	return c
}
