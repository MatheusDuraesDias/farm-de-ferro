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
	Random50Songs, _ := a.Db.Random50Songs(a.Ctx)
	Random50NewSongs, _ := a.Db.Random50NewSongs(a.Ctx)
	RandomIndieSongs, _ := a.Db.Random20IndieSongs(a.Ctx)

	params := domain.UserMusicPreferences{
		UserPreferences:  userSets["favoriteStyles"],
		UserFollowStyles: userSets["followStyles"],
		UserLastLikes:    userSets["lastLikedStyles"],
	}

	allSongs := append(Random50NewSongs, Random50Songs...)

	unviewed, err := a.Neo4JDb.GetUnviewedPosts(UserId, allSongs)
	if err != nil {
		return nil, err
	}

	unviewedIndies, err := a.Neo4JDb.GetUnviewedPosts(UserId, RandomIndieSongs)

	var filteredSongs []domain.Song

	if len(filteredSongs) < 40 {
		filteredSongs = unviewed
	} else {
		filteredSongs = filterOfSongs(unviewed, params.UserPreferences, params.UserLastLikes, params.UserFollowStyles, 40)
	}

	allSongs = append(filteredSongs, unviewedIndies...)
	allSongs = removeDuplicates(allSongs)

	rand.Shuffle(len(allSongs), func(i, j int) {
		allSongs[i], allSongs[j] = allSongs[j], allSongs[i]
	})

	res := []string{}
	if len(allSongs) > limit {
		for i := 0; i < limit; i++ {
			res = append(res, allSongs[i].Id)
		}
	} else {
		for song := range allSongs {
			res = append(res, allSongs[song].Id)
		}
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

// vou comentar por que achei mei bagunça o algo-ritmo
func filterOfSongs(songs []domain.Song, userPreferences []string, likes []string, follows []string, minimumFiltered int) []domain.Song {
	var filteredSongs []domain.Song
	notAddedSongs := make([]domain.Song, 0)

	// maps para otimizar a pesquisa(sem ter que fazer um monte de for²)
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

	// filtro por estilos seguidos pelo usuário
	for _, song := range songs {
		if _, exists := followSet[song.Style]; exists {
			filteredSongs = append(filteredSongs, song)
		} else {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	// Reset
	songs = notAddedSongs
	notAddedSongs = make([]domain.Song, 0)

	// filtro por likes do usuário
	for _, song := range songs {
		if _, exists := likeSet[song.Style]; exists {
			filteredSongs = append(filteredSongs, song)
		} else {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	// Reset
	songs = notAddedSongs
	notAddedSongs = make([]domain.Song, 0)

	// filtro por preferências do usuário
	for _, song := range songs {
		if _, exists := preferenceSet[song.Style]; exists {
			filteredSongs = append(filteredSongs, song)
		} else {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	// verificação de mínimo de músicas filtradas
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
