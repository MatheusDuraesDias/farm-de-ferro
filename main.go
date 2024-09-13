package main

import (
	"algorithm/mod/algoritmo"
	"algorithm/mod/algoritmo/database"
	"algorithm/mod/algoritmo/domain"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getTestRecommendedSongs(c echo.Context) error {
	mockData := domain.UserMusicPreferences{
		UserPreferences: []string{
			"Rock", "Pop", "Jazz", "Hip-Hop", "Classical", "Electronic", "Reggae", "Country", "Blues", "Metal",
		},
		UserFollowStyles: []string{
			"Pop", "Jazz",
		},
		UserLastLikes: []string{
			"Rock", "Blues", "Hip-Hop",
		},
		Random50NewSongs: []domain.Song{
			{Name: "NewSong1", Style: "Rock"}, {Name: "NewSong2", Style: "Pop"}, {Name: "NewSong3", Style: "Jazz"},
			{Name: "NewSong4", Style: "Hip-Hop"}, {Name: "NewSong5", Style: "Classical"}, {Name: "NewSong6", Style: "Electronic"},
			{Name: "NewSong7", Style: "Reggae"}, {Name: "NewSong8", Style: "Country"}, {Name: "NewSong9", Style: "Blues"},
			{Name: "NewSong10", Style: "Metal"}, {Name: "NewSong11", Style: "Rock"}, {Name: "NewSong12", Style: "Pop"},
			{Name: "NewSong13", Style: "Jazz"}, {Name: "NewSong14", Style: "Hip-Hop"}, {Name: "NewSong15", Style: "Classical"},
			{Name: "NewSong16", Style: "Electronic"}, {Name: "NewSong17", Style: "Reggae"}, {Name: "NewSong18", Style: "Country"},
			{Name: "NewSong19", Style: "Blues"}, {Name: "NewSong20", Style: "Metal"}, {Name: "NewSong21", Style: "Rock"},
			{Name: "NewSong22", Style: "Pop"}, {Name: "NewSong23", Style: "Jazz"}, {Name: "NewSong24", Style: "Hip-Hop"},
			{Name: "NewSong25", Style: "Classical"}, {Name: "NewSong26", Style: "Electronic"}, {Name: "NewSong27", Style: "Reggae"},
			{Name: "NewSong28", Style: "Country"}, {Name: "NewSong29", Style: "Blues"}, {Name: "NewSong30", Style: "Metal"},
			{Name: "NewSong31", Style: "Rock"}, {Name: "NewSong32", Style: "Pop"}, {Name: "NewSong33", Style: "Jazz"},
			{Name: "NewSong34", Style: "Hip-Hop"}, {Name: "NewSong35", Style: "Classical"}, {Name: "NewSong36", Style: "Electronic"},
			{Name: "NewSong37", Style: "Reggae"}, {Name: "NewSong38", Style: "Country"}, {Name: "NewSong39", Style: "Blues"},
			{Name: "NewSong40", Style: "Metal"}, {Name: "NewSong41", Style: "Rock"}, {Name: "NewSong42", Style: "Pop"},
			{Name: "NewSong43", Style: "Jazz"}, {Name: "NewSong44", Style: "Hip-Hop"}, {Name: "NewSong45", Style: "Classical"},
			{Name: "NewSong46", Style: "Electronic"}, {Name: "NewSong47", Style: "Reggae"}, {Name: "NewSong48", Style: "Country"},
			{Name: "NewSong49", Style: "Blues"}, {Name: "NewSong50", Style: "Metal"},
		},
		Random50Songs: []domain.Song{
			{Name: "Song1", Style: "Rock"}, {Name: "Song2", Style: "Pop"}, {Name: "Song3", Style: "Jazz"},
			{Name: "Song4", Style: "Hip-Hop"}, {Name: "Song5", Style: "Classical"}, {Name: "Song6", Style: "Electronic"},
			{Name: "Song7", Style: "Reggae"}, {Name: "Song8", Style: "Country"}, {Name: "Song9", Style: "Blues"},
			{Name: "Song10", Style: "Metal"}, {Name: "Song11", Style: "Rock"}, {Name: "Song12", Style: "Pop"},
			{Name: "Song13", Style: "Jazz"}, {Name: "Song14", Style: "Hip-Hop"}, {Name: "Song15", Style: "Classical"},
			{Name: "Song16", Style: "Electronic"}, {Name: "Song17", Style: "Reggae"}, {Name: "Song18", Style: "Country"},
			{Name: "Song19", Style: "Blues"}, {Name: "Song20", Style: "Metal"}, {Name: "Song21", Style: "Rock"},
			{Name: "Song22", Style: "Pop"}, {Name: "Song23", Style: "Jazz"}, {Name: "Song24", Style: "Hip-Hop"},
			{Name: "Song25", Style: "Classical"}, {Name: "Song26", Style: "Electronic"}, {Name: "Song27", Style: "Reggae"},
			{Name: "Song28", Style: "Country"}, {Name: "Song29", Style: "Blues"}, {Name: "Song30", Style: "Metal"},
			{Name: "Song31", Style: "Rock"}, {Name: "Song32", Style: "Pop"}, {Name: "Song33", Style: "Jazz"},
			{Name: "Song34", Style: "Hip-Hop"}, {Name: "Song35", Style: "Classical"}, {Name: "Song36", Style: "Electronic"},
			{Name: "Song37", Style: "Reggae"}, {Name: "Song38", Style: "Country"}, {Name: "Song39", Style: "Blues"},
			{Name: "Song40", Style: "Metal"}, {Name: "Song41", Style: "Rock"}, {Name: "Song42", Style: "Pop"},
			{Name: "Song43", Style: "Jazz"}, {Name: "Song44", Style: "Hip-Hop"}, {Name: "Song45", Style: "Classical"},
			{Name: "Song46", Style: "Electronic"}, {Name: "Song47", Style: "Reggae"}, {Name: "Song48", Style: "Country"},
			{Name: "Song49", Style: "Blues"}, {Name: "Song50", Style: "Metal"},
		},
		Random20IndieSongs: []domain.Song{
			{Name: "IndieSong1", Style: "Indie Rock"}, {Name: "IndieSong2", Style: "Indie Pop"}, {Name: "IndieSong3", Style: "Indie Folk"},
			{Name: "IndieSong4", Style: "Indie Rock"}, {Name: "IndieSong5", Style: "Indie Pop"}, {Name: "IndieSong6", Style: "Indie Folk"},
			{Name: "IndieSong7", Style: "Indie Rock"}, {Name: "IndieSong8", Style: "Indie Pop"}, {Name: "IndieSong9", Style: "Indie Folk"},
			{Name: "IndieSong10", Style: "Indie Rock"}, {Name: "IndieSong11", Style: "Indie Pop"}, {Name: "IndieSong12", Style: "Indie Folk"},
			{Name: "IndieSong13", Style: "Indie Rock"}, {Name: "IndieSong14", Style: "Indie Pop"}, {Name: "IndieSong15", Style: "Indie Folk"},
			{Name: "IndieSong16", Style: "Indie Rock"}, {Name: "IndieSong17", Style: "Indie Pop"}, {Name: "IndieSong18", Style: "Indie Folk"},
			{Name: "IndieSong19", Style: "Indie Rock"}, {Name: "IndieSong20", Style: "Indie Pop"},
		},
	}

	algo := algoritmo.Algo{}

	recommendedSongs := algo.Algoritmo(mockData)
	return c.JSON(http.StatusOK, recommendedSongs)
}

func getRecommendedSongs(c echo.Context) error {
	params := new(domain.UserMusicPreferences)
	if err := c.Bind(params); err != nil {
		return err
	}

	algo := algoritmo.Algo{}

	recommendedSongs := algo.Algoritmo(*params)
	return c.JSON(http.StatusOK, recommendedSongs)
}

func main() {
	e := echo.New()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(err.Error())
	}

	e.GET("/test-recommended-songs", getTestRecommendedSongs)
	e.GET("/recommended-songs", getRecommendedSongs)
	uri := os.Getenv("URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	db := database.NewDatabase(client, cancel)
	db.Ping(ctx)
	db.GetFavoriteStyles(ctx, "66c868e082cd6161c37c0f48")
	db.UserFollowStyles(ctx, "66c868e082cd6161c37c0f48")

	// e.Logger.Fatal(e.Start(":8080"))
}
