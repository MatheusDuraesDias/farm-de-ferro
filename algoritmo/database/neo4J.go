package database

import (
	"algorithm/mod/algoritmo/domain"
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type NeoDatabase struct {
	Driver neo4j.DriverWithContext
}

func (db *NeoDatabase) GetUnviewedPosts(userId string, posts []domain.Song) ([]domain.Song, error) {
	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (u:User {id: $userId})
			WITH u
			MATCH (u)-[:VIEWED]->(p:Post)
			RETURN p.id AS postId
		`
		records, err := tx.Run(context.Background(), query, map[string]interface{}{
			"userId": userId,
		})
		if err != nil {
			return nil, err
		}

		var postIds []string
		for records.Next(context.Background()) {
			postId, _ := records.Record().Get("postId")
			postIds = append(postIds, postId.(string))
		}
		return postIds, nil
	})
	if err != nil {
		return nil, err
	}

	res := filterSongs(posts, result.([]string))

	return res, nil
}

func filterSongs(songs []domain.Song, allowedIds []string) []domain.Song {
	allowedMap := make(map[string]struct{}, len(allowedIds))
	for _, id := range allowedIds {
		allowedMap[id] = struct{}{}
	}

	var filteredSongs []domain.Song
	for _, song := range songs {
		if _, exists := allowedMap[song.Id]; !exists {
			filteredSongs = append(filteredSongs, song)
		}
	}

	return filteredSongs
}

func (db *NeoDatabase) MarkSongsAsViewed(userId string, postIds []string) error {
	session := db.Driver.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MERGE (u:User {id: $userId})
			WITH u
			UNWIND $postIds AS postId
			MERGE (p:Post {id: postId})
			MERGE (u)-[:VIEWED]->(p)
		`
		_, err := tx.Run(context.Background(), query, map[string]interface{}{
			"userId":  userId,
			"postIds": postIds,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
