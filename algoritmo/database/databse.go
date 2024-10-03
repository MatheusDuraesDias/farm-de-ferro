package database

import (
	"context"
	"log"
	"time"

	"algorithm/mod/algoritmo/domain"

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
	return nil
}

func (db *Database) GetAllUserStyles(ctx context.Context, userId string) map[string][]string {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": id}
	projection := bson.M{"favorites": 1, "followings": 1, "likes": 1, "_id": 0}
	findOneOptions := options.FindOne().SetProjection(projection)

	var userData bson.M
	if err := db.UserCol.FindOne(ctx, filter, findOneOptions).Decode(&userData); err != nil {
		log.Fatal(err)
	}

	var favorites, followings, likes []primitive.ObjectID

	if favArr, ok := userData["favorites"].(primitive.A); ok {
		for _, item := range favArr {
			if str, ok := item.(string); ok {
				id, err := primitive.ObjectIDFromHex(str)
				if err != nil {
					log.Fatal(err)
				}
				favorites = append(favorites, id)
			}
		}
	}

	if followArr, ok := userData["followings"].(primitive.A); ok {
		for _, item := range followArr {
			if str, ok := item.(string); ok {
				id, err := primitive.ObjectIDFromHex(str)
				if err != nil {
					log.Fatal(err)
				}
				followings = append(followings, id)
			}
		}
	}

	if likesArr, ok := userData["likes"].(primitive.A); ok {
		for _, item := range likesArr {
			if str, ok := item.(string); ok {
				id, err := primitive.ObjectIDFromHex(str)
				if err != nil {
					log.Fatal(err)
				}
				likes = append(likes, id)
			}
		}
	}

	favoriteStyles := db.GetStylesByPostIds(ctx, favorites)
	followStyles := db.GetStylesByUserIds(ctx, followings)
	lastLikedStyles := db.GetStylesByPostIds(ctx, likes)

	result := map[string][]string{
		"favoriteStyles":  favoriteStyles,
		"followStyles":    followStyles,
		"lastLikedStyles": lastLikedStyles,
	}

	return result
}
func (db *Database) GetStylesByPostIds(ctx context.Context, postIds []primitive.ObjectID) []string {
	if len(postIds) == 0 {
		return nil
	}

	matchStage := bson.D{
		{"$match", bson.M{"_id": bson.M{"$in": postIds}, "active": true}},
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

	var styles []string
	for _, doc := range result {
		if style, ok := doc["style"].(string); ok {
			styles = append(styles, style)
		}
	}

	return removeDuplicates(styles)
}

func (db *Database) GetStylesByUserIds(ctx context.Context, userIds []primitive.ObjectID) []string {
	if len(userIds) == 0 {
		return nil
	}

	matchStage := bson.D{
		{"$match", bson.M{"_id": bson.M{"$in": userIds}, "active": true}},
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

	var styles []string
	for _, doc := range result {
		if stylesPosted, ok := doc["stylesPosted"].(bson.M); ok {
			for style := range stylesPosted {
				styles = append(styles, style)
			}
		}
	}

	return removeDuplicates(styles)
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

func (db *Database) RandomSongs(ctx context.Context, limit int) ([]domain.Song, error) {

	limit = int(float64(limit) * 1.5)

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"isActive", true}}}},
		{{"$sample", bson.D{{"size", 50}}}},
	}

	cursor, err := db.PostCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var songs []domain.PostDTO
	if err := cursor.All(ctx, &songs); err != nil {
		return nil, err
	}

	var songsRes []domain.Song
	songsRes = []domain.Song{}

	for i := 0; i < len(songs); i++ {
		songsRes = append(songsRes, domain.PostDTOToSong(songs[i]))
	}

	return songsRes, nil
}

func (db *Database) RandomNewSongs(ctx context.Context, limit int) ([]domain.Song, error) {
	now := time.Now()
	dateLimit := now.AddDate(0, 0, -50)

	limit = int(float64(limit) * 1.5)

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"date", bson.D{{"$gte", dateLimit}}}, {"isActive", true}}}},
		{{"$sample", bson.D{{"size", limit}}}},
	}

	cursor, err := db.PostCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var songs []domain.PostDTO
	if err := cursor.All(ctx, &songs); err != nil {
		return nil, err
	}

	var songsRes []domain.Song
	songsRes = []domain.Song{}

	for i := 0; i < len(songs); i++ {
		songsRes = append(songsRes, domain.PostDTOToSong(songs[i]))
	}

	return songsRes, nil
}

func (db *Database) RandomIndieSongs(ctx context.Context, limit int) ([]domain.Song, error) {
	limit = int(float64(limit) * 0.8)

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"style", "indie"}, {"isActive", true}}}},
		{{"$sample", bson.D{{"size", limit}}}},
	}

	cursor, err := db.PostCol.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var songs []domain.PostDTO
	if err := cursor.All(ctx, &songs); err != nil {
		return nil, err
	}

	var songsRes []domain.Song
	songsRes = []domain.Song{}

	for i := 0; i < len(songs); i++ {
		songsRes = append(songsRes, domain.PostDTOToSong(songs[i]))
	}

	return songsRes, nil
}
