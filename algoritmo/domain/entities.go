package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Song struct {
	Id    string `json:"id"`
	Style string `json:"style"`
}

type UserMusicPreferences struct {
	UserPreferences  []string `json:"user_preferences"`
	UserFollowStyles []string `json:"user_follow_styles"`
	UserLastLikes    []string `json:"user_last_likes"`
}

type PostDTO struct {
	ID    primitive.ObjectID `bson:"_id"`
	Style string             `bson:"style"`
}

func PostDTOToSong(post PostDTO) Song {
	objectIDString := post.ID.Hex()
	song := Song{
		Id:    objectIDString,
		Style: post.Style,
	}

	return song
}
