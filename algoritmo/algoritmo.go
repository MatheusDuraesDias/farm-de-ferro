package algoritmo

import (
	"algorithm/mod/algoritmo/database"
	"algorithm/mod/algoritmo/domain"
	"context"
	"math/rand"
)

type Algo struct {
	Db  database.Database
	Ctx context.Context
}

func (a *Algo) Algoritmo(UserId string) []domain.Song {
	userSets := a.Db.GetAllUserStyles(a.Ctx, UserId)
	Random50Songs, _ := a.Db.Random50Songs(a.Ctx)
	Random50NewSongs, _ := a.Db.Random50NewSongs(a.Ctx)
	RandomIndieSongs, _ := a.Db.Random20IndieSongs(a.Ctx)

	params := domain.UserMusicPreferences{
		UserPreferences:    userSets["favoriteStyles"],
		UserFollowStyles:   userSets["followStyles"],
		UserLastLikes:      userSets["lastLikedStyles"],
		Random50Songs:      Random50Songs,
		Random50NewSongs:   Random50NewSongs,
		Random20IndieSongs: RandomIndieSongs,
	}

	var filteredNewSongs []domain.Song

	var filteredRandomSongs []domain.Song

	if len(params.Random50NewSongs) < 10 {
		filteredNewSongs = params.Random50NewSongs
	} else {
		filteredNewSongs = filterOfSongs(params.Random50NewSongs, params.UserPreferences, params.UserLastLikes, params.UserFollowStyles, 10)
	}

	if len(params.Random50Songs)<10{
		filteredRandomSongs = params.Random50Songs
	}else{
		filteredRandomSongs = filterOfSongs(params.Random50Songs, params.UserPreferences, params.UserLastLikes, params.UserFollowStyles, 10)
	}

	allSongs := filteredNewSongs
	allSongs = append(allSongs, filteredRandomSongs...)
	allSongs = append(allSongs, params.Random20IndieSongs...)

	rand.Shuffle(len(allSongs), func(i, j int) {
		allSongs[i], allSongs[j] = allSongs[j], allSongs[i]
	})

	if len(allSongs) > 25{
		return allSongs[:25]
	}else{
		return allSongs
	}
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

// func filterOfSongs2(songs []domain.Song, userPreferences []string, likes []string, follows []string, minimumFiltered int) []domain.Song {

// 	var filteredSongs []domain.Song

// 	added := false
// 	var notAddedSongs []domain.Song

// 	for _, song := range songs {
// 		for _, follow := range follows {
// 			if follow == song.Style {
// 				filteredSongs = append(filteredSongs, song)
// 				added = true
// 			}
// 		}
// 		if !added {
// 			notAddedSongs = append(notAddedSongs, song)
// 		}
// 	}

// 	songs = notAddedSongs
// 	added = false
// 	for _, song := range songs {
// 		for _, like := range likes {
// 			if like == song.Style {
// 				filteredSongs = append(filteredSongs, song)
// 				added = true
// 			}
// 		}
// 		if !added {
// 			notAddedSongs = append(notAddedSongs, song)
// 		}
// 	}

// 	songs = notAddedSongs
// 	added = false
// 	for _, song := range songs {
// 		for _, preference := range userPreferences {
// 			if preference == song.Style {
// 				filteredSongs = append(filteredSongs, song)
// 				added = true
// 			}
// 		}
// 		if !added {
// 			notAddedSongs = append(notAddedSongs, song)
// 		}
// 	}

// 	// esse aqui vai verificar se tem músicas filtradas suficientes
// 	var addedIds []int
// 	for len(filteredSongs) < minimumFiltered {
// 		alreadyAdded := false
// 		randId := rand.Intn(len(notAddedSongs))
// 		for _, yeah := range addedIds {
// 			if yeah == randId {
// 				alreadyAdded = true
// 			}
// 		}
// 		if !alreadyAdded {
// 			filteredSongs = append(filteredSongs, notAddedSongs[randId])
// 			addedIds = append(addedIds, randId)
// 		}
// 	}

// 	return filteredSongs
// }
