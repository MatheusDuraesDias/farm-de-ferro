package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Song struct {
	Name   string `json:"name"`
	Style  string `json:"style"`
	Artist string `json:"artist"`
}

type UserMusicPreferences struct {
	UserPreferences []string `json:"user_preferences"`
	//mongo logic:
	// query for userFavorites, then take the ids(max of 100) and search for the styles of the following ids.
	UserFollowStyles []string `json:"user_follow_styles"`
	//mongo logic:
	// query for userFollowings, then take the ids(max of 20) and search for the firt music of each artist, then take the style of that music.
	UserLastLikes []string `json:"user_last_likes"`
	//mongo logic:
	// like UserPreferences logic, but it takes the liked music Ids, not the userFavorites.
	Random50Songs []Song `json:"random_50_songs"`
	//mongo logic:
	// dar um get no count de documents na collection de posts, com esse count eu aleatorizo 50 ids de posts e faço uma query buscando eles(e adapto
	// pra struct de Song)
	Random50NewSongs []Song `json:"random_50_new_songs"`
	//mongo logic:
	// tipo o de cima, só que limitando a músicas que tenham lançado no máximo há 50 dias
	Random20IndieSongs []Song `json:"random_20_indie_songs"`
	//mongo logic:
	// Tipo o Random50Songs, só que se limita a músicas com o
}

type PostDTO struct {
	ID     primitive.ObjectID `bson:"_id"`
	Title  string             `bson:"name"`
	Artist string             `bson:"userId"`
	Style  string             `bson:"style"`
}

func PostDTOToSong(post PostDTO) Song {
	song := Song{
		Name:   post.Title,
		Style:  post.Style,
		Artist: post.Artist,
	}

	return song
}
