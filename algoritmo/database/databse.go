package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	Client  *mongo.Client
	Cancel  context.CancelFunc
	UserCol *mongo.Collection
	PostCol *mongo.Collection
}

func NewDatabase(client *mongo.Client, cancel context.CancelFunc) *Database {
	db := new(Database)
	db.Client = client
	db.Cancel = cancel
	db.UserCol = client.Database("dbFroggers").Collection("users")
	db.PostCol = client.Database("dbFroggers").Collection("posts")
	return db
}

func (db *Database) Close(ctx context.Context) {
	defer db.Cancel()
	defer func() {
		if err := db.Client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func (db *Database) Ping(ctx context.Context) error {

	if err := db.Client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func (db *Database) GetFavoriteStyles(ctx context.Context, userId string) []string {
	fmt.Println(userId)
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}
	projection := bson.M{"favorites": 1, "_id": 0}
	findOneOptions := options.FindOne().SetProjection(projection)

	var favoriteIds bson.M
	if err := db.UserCol.FindOne(ctx, filter, findOneOptions).Decode(&favoriteIds); err != nil {
		log.Fatal(err)
	}

	var idsParam []primitive.ObjectID
	if arr, ok := favoriteIds["favorites"].(primitive.A); ok {
		for _, item := range arr {
			if str, ok := item.(string); ok {
				id, err := primitive.ObjectIDFromHex(str)
				if err != nil {
					log.Fatal(err)
				}
				idsParam = append(idsParam, id)
			}
		}
	}

	matchStage := bson.D{
		{"$match", bson.M{"_id": bson.M{"$in": idsParam}}},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$style"},
		}},
	}
	projectStage := bson.D{
		{"$project", bson.M{"style": "$_id", "_id": 0}},
	}
	pipeline := mongo.Pipeline{matchStage, groupStage, projectStage}
	cursor, err := db.PostCol.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	var res []string
	for _, doc := range result {
		if style, ok := doc["style"].(string); ok {
			res = append(res, style)
		}
	}

	return res
}

func (db *Database) GetFollowStyles(ctx context.Context, userId string) []string {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}
	projection := bson.M{"followings": 1, "_id": 0}
	findOneOptions := options.FindOne().SetProjection(projection)

	var followings bson.M
	if err := db.UserCol.FindOne(ctx, filter, findOneOptions).Decode(&followings); err != nil {
		log.Fatal(err)
	}

	var idsParam []primitive.ObjectID
	if arr, ok := followings["followings"].(primitive.A); ok {
		for _, item := range arr {
			if str, ok := item.(string); ok {
				id, err := primitive.ObjectIDFromHex(str)
				if err != nil {
					log.Fatal(err)
				}
				idsParam = append(idsParam, id)
			}
		}
	}

	matchStage := bson.D{
		{"$match", bson.M{"_id": bson.M{"$in": idsParam}}},
	}
	projectStage := bson.D{
		{"$project", bson.M{"stylesPosted": 1, "_id": 0}},
	}
	pipeline := mongo.Pipeline{matchStage, projectStage}
	cursor, err := db.UserCol.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	var res []string
	for _, doc := range result {
		if stylesPosted, ok := doc["stylesPosted"].(bson.M); ok {
			for style := range stylesPosted {
				res = append(res, style)
			}
		}
	}

	res = removeDuplicates(res)
	return res
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func (db *Database) GetLastLikedStyles(ctx context.Context, userId string) []string {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}
	projection := bson.M{"likes": 1, "_id": 0}
	findOneOptions := options.FindOne().SetProjection(projection)

	var likes bson.M
	if err := db.UserCol.FindOne(ctx, filter, findOneOptions).Decode(&likes); err != nil {
		log.Fatal(err)
	}

	var idsParam []primitive.ObjectID
	if arr, ok := likes["likes"].(primitive.A); ok {
		for _, item := range arr {
			if str, ok := item.(string); ok {
				id, err := primitive.ObjectIDFromHex(str)
				if err != nil {
					log.Fatal(err)
				}
				idsParam = append(idsParam, id)
			}
		}
	}

	matchStage := bson.D{
		{"$match", bson.M{"_id": bson.M{"$in": idsParam}}},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$style"},
		}},
	}
	projectStage := bson.D{
		{"$project", bson.M{"style": "$_id", "_id": 0}},
	}
	pipeline := mongo.Pipeline{matchStage, groupStage, projectStage}
	cursor, err := db.PostCol.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	var res []string
	for _, doc := range result {
		if style, ok := doc["style"].(string); ok {
			res = append(res, style)
		}
	}

	fmt.Println(res)
	return res
}