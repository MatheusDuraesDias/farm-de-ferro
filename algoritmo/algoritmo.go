package algoritmo

import (
	"algorithm/mod/algoritmo/domain"
	"math/rand"
)

type Algo struct {
	Idk string
}

func (a *Algo) Algoritmo(params domain.UserMusicPreferencesDTO) []domain.Song {

	filteredNewSongs := filterOfSongs(params.Random50NewSongs, params.UserPreferences, params.UserLastLikes, params.UserFollowStyles, 10)
	filteredRandomSongs := filterOfSongs(params.Random50Songs, params.UserPreferences, params.UserLastLikes, params.UserFollowStyles, 10)
	filteredIndieSongs := filterOfSongs(params.Random20IndieSongs, params.UserPreferences, params.UserLastLikes, params.UserFollowStyles, 5)

	allSongs := filteredNewSongs
	allSongs = append(allSongs, filteredRandomSongs...)
	allSongs = append(allSongs, filteredIndieSongs...)

	rand.Shuffle(len(allSongs), func(i, j int) {
		allSongs[i], allSongs[j] = allSongs[j], allSongs[i]
	})

	return allSongs[:25]
}

func filterOfSongs(songs []domain.Song, userPreferences []string, likes []string, follows []string, minimumFiltered int) []domain.Song {

	var filteredSongs []domain.Song

	added := false
	var notAddedSongs []domain.Song

	for _, song := range songs {
		for _, follow := range follows {
			if follow == song.Style {
				filteredSongs = append(filteredSongs, song)
				added = true
			}
		}
		if !added {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	songs = notAddedSongs
	added = false
	for _, song := range songs {
		for _, like := range likes {
			if like == song.Style {
				filteredSongs = append(filteredSongs, song)
				added = true
			}
		}
		if !added {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	songs = notAddedSongs
	added = false
	for _, song := range songs {
		for _, preference := range userPreferences {
			if preference == song.Style {
				filteredSongs = append(filteredSongs, song)
				added = true
			}
		}
		if !added {
			notAddedSongs = append(notAddedSongs, song)
		}
	}

	// rand.Shuffle(len(filteredSongs), func(i, j int) {
	// 	filteredSongs[i], filteredSongs[j] = filteredSongs[j], filteredSongs[i]
	// })

	var addedIds []int
	for len(filteredSongs) < minimumFiltered {
		alreadyAdded := false
		randId := rand.Intn(len(notAddedSongs))
		for _, yeah := range addedIds {
			if yeah == randId {
				alreadyAdded = true
			}
		}
		if !alreadyAdded {
			filteredSongs = append(filteredSongs, notAddedSongs[randId])
			addedIds = append(addedIds, randId)
		}
	}

	return filteredSongs
}