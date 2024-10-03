package algoritmo

import (
	"algorithm/mod/algoritmo/database"
	"algorithm/mod/algoritmo/domain"
	"context"
	"math/rand"
)

type Algo struct {
	Db      database.Database
	Neo4JDb database.NeoDatabase
	Ctx     context.Context
}

func (a *Algo) Algoritmo(UserId string, limit int, viewedPosts []string) ([]string, error) {
	err := a.Neo4JDb.MarkSongsAsViewed(UserId, viewedPosts)
	if err != nil {
		return nil, err
	}

	userSets := a.Db.GetAllUserStyles(a.Ctx, UserId)
	RandomSongs, _ := a.Db.RandomSongs(a.Ctx, limit)
	RandomNewSongs, _ := a.Db.RandomNewSongs(a.Ctx, limit)
	RandomIndieSongs, _ := a.Db.RandomIndieSongs(a.Ctx, limit)

	params := domain.UserMusicPreferences{
		UserPreferences:  userSets["favoriteStyles"],
		UserFollowStyles: userSets["followStyles"],
		UserLastLikes:    userSets["lastLikedStyles"],
	}

	allSongs := append(RandomNewSongs, RandomSongs...)

	var filteredSongs []domain.Song

	if len(filteredSongs) < int(float64(limit)*1.8) {
		filteredSongs = allSongs
	} else {
		filteredSongs = filterOfSongs(allSongs, params.UserPreferences, params.UserLastLikes, params.UserFollowStyles, int(float64(limit)*1.8))
	}

	allSongs = append(filteredSongs, RandomIndieSongs...)
	allSongs = removeDuplicates(allSongs)

	allSongsIds, err := a.Neo4JDb.GetUnviewedPosts(UserId, allSongs)
	if err != nil {
		return nil, err
	}

	rand.Shuffle(len(allSongsIds), func(i, j int) {
		allSongsIds[i], allSongsIds[j] = allSongsIds[j], allSongsIds[i]
	})

	res := []string{}
	if len(allSongs) > limit {
		res = allSongsIds[:limit]
	} else {
		res = allSongsIds
	}

	return res, nil
}

func removeDuplicates(elements []domain.Song) []domain.Song {
	encountered := map[domain.Song]bool{}
	result := []domain.Song{}

	for _, v := range elements {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func filterOfSongs(songs []domain.Song, userPreferences []string, likes []string, follows []string, minimumFiltered int) []domain.Song {
	var filteredSongs []domain.Song
	notAddedSongs := make([]domain.Song, 0)

	followSet := make(map[string]struct{})
	for _, follow := range follows {
		followSet[follow] = struct{}{}
	}

	likeSet := make(map[string]struct{})
	for _, like := range likes {
		likeSet[like] = struct{}{}
	}

	preferenceSet := make(map[string]struct{})
	for _, preference := range userPreferences {
		preferenceSet[preference] = struct{}{}
	}

	for _, song := range songs {
		if _, exists := followSet[song.Style]; exists {
			filteredSongs = append(filteredSongs, song)
		} else {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	songs = notAddedSongs
	notAddedSongs = make([]domain.Song, 0)

	for _, song := range songs {
		if _, exists := likeSet[song.Style]; exists {
			filteredSongs = append(filteredSongs, song)
		} else {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	songs = notAddedSongs
	notAddedSongs = make([]domain.Song, 0)

	for _, song := range songs {
		if _, exists := preferenceSet[song.Style]; exists {
			filteredSongs = append(filteredSongs, song)
		} else {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	addedIds := make(map[int]struct{})
	for len(filteredSongs) < minimumFiltered {
		randId := rand.Intn(len(notAddedSongs))
		if _, alreadyAdded := addedIds[randId]; !alreadyAdded {
			filteredSongs = append(filteredSongs, notAddedSongs[randId])
			addedIds[randId] = struct{}{}
		}
	}

	return filteredSongs
}
