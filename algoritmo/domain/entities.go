package domain

type Song struct {
	Name  string `json:"name"`
	Style string `json:"style"`
}

type UserMusicPreferencesDTO struct {
	UserPreferences    []string `json:"user_preferences"`
	UserFollowStyles   []string `json:"user_follow_styles"`
	UserLastLikes      []string `json:"user_last_likes"`
	Random50NewSongs   []Song   `json:"random_50_new_songs"`
	Random50Songs      []Song   `json:"random_50_songs"`
	Random20IndieSongs []Song   `json:"random_20_indie_songs"`
}
