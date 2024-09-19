package main

import (
	// "algorithm/mod/algoritmo"
	"algorithm/mod/algoritmo/database"
	// "algorithm/mod/algoritmo/domain"
	"context"
	"fmt"
	// "go/printer"
	// "net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	// "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func getTestRecommendedSongs(c echo.Context) error {
// 	mockData := domain.UserMusicPreferences{
// 		UserPreferences: []string{
// 			"Rock", "Pop", "Jazz", "Hip-Hop", "Classical", "Electronic", "Reggae", "Country", "Blues", "Metal",
// 		},
// 		UserFollowStyles: []string{
// 			"Pop", "Jazz",
// 		},
// 		UserLastLikes: []string{
// 			"Rock", "Blues", "Hip-Hop",
// 		},
// 		Random50NewSongs: []domain.Song{
// 			{Id: "NewSong1", Style: "Rock"}, {Id: "NewSong2", Style: "Pop"}, {Id: "NewSong3", Style: "Jazz"},
// 			{Id: "NewSong4", Style: "Hip-Hop"}, {Id: "NewSong5", Style: "Classical"}, {Id: "NewSong6", Style: "Electronic"},
// 			{Id: "NewSong7", Style: "Reggae"}, {Id: "NewSong8", Style: "Country"}, {Id: "NewSong9", Style: "Blues"},
// 			{Id: "NewSong10", Style: "Metal"}, {Id: "NewSong11", Style: "Rock"}, {Id: "NewSong12", Style: "Pop"},
// 			{Id: "NewSong13", Style: "Jazz"}, {Id: "NewSong14", Style: "Hip-Hop"}, {Id: "NewSong15", Style: "Classical"},
// 			{Id: "NewSong16", Style: "Electronic"}, {Id: "NewSong17", Style: "Reggae"}, {Id: "NewSong18", Style: "Country"},
// 			{Id: "NewSong19", Style: "Blues"}, {Id: "NewSong20", Style: "Metal"}, {Id: "NewSong21", Style: "Rock"},
// 			{Id: "NewSong22", Style: "Pop"}, {Id: "NewSong23", Style: "Jazz"}, {Id: "NewSong24", Style: "Hip-Hop"},
// 			{Id: "NewSong25", Style: "Classical"}, {Id: "NewSong26", Style: "Electronic"}, {Id: "NewSong27", Style: "Reggae"},
// 			{Id: "NewSong28", Style: "Country"}, {Id: "NewSong29", Style: "Blues"}, {Id: "NewSong30", Style: "Metal"},
// 			{Id: "NewSong31", Style: "Rock"}, {Id: "NewSong32", Style: "Pop"}, {Id: "NewSong33", Style: "Jazz"},
// 			{Id: "NewSong34", Style: "Hip-Hop"}, {Id: "NewSong35", Style: "Classical"}, {Id: "NewSong36", Style: "Electronic"},
// 			{Id: "NewSong37", Style: "Reggae"}, {Id: "NewSong38", Style: "Country"}, {Id: "NewSong39", Style: "Blues"},
// 			{Id: "NewSong40", Style: "Metal"}, {Id: "NewSong41", Style: "Rock"}, {Id: "NewSong42", Style: "Pop"},
// 			{Id: "NewSong43", Style: "Jazz"}, {Id: "NewSong44", Style: "Hip-Hop"}, {Id: "NewSong45", Style: "Classical"},
// 			{Id: "NewSong46", Style: "Electronic"}, {Id: "NewSong47", Style: "Reggae"}, {Id: "NewSong48", Style: "Country"},
// 			{Id: "NewSong49", Style: "Blues"}, {Id: "NewSong50", Style: "Metal"},
// 		},
// 		Random50Songs: []domain.Song{
// 			{Id: "Song1", Style: "Rock"}, {Id: "Song2", Style: "Pop"}, {Id: "Song3", Style: "Jazz"},
// 			{Id: "Song4", Style: "Hip-Hop"}, {Id: "Song5", Style: "Classical"}, {Id: "Song6", Style: "Electronic"},
// 			{Id: "Song7", Style: "Reggae"}, {Id: "Song8", Style: "Country"}, {Id: "Song9", Style: "Blues"},
// 			{Id: "Song10", Style: "Metal"}, {Id: "Song11", Style: "Rock"}, {Id: "Song12", Style: "Pop"},
// 			{Id: "Song13", Style: "Jazz"}, {Id: "Song14", Style: "Hip-Hop"}, {Id: "Song15", Style: "Classical"},
// 			{Id: "Song16", Style: "Electronic"}, {Id: "Song17", Style: "Reggae"}, {Id: "Song18", Style: "Country"},
// 			{Id: "Song19", Style: "Blues"}, {Id: "Song20", Style: "Metal"}, {Id: "Song21", Style: "Rock"},
// 			{Id: "Song22", Style: "Pop"}, {Id: "Song23", Style: "Jazz"}, {Id: "Song24", Style: "Hip-Hop"},
// 			{Id: "Song25", Style: "Classical"}, {Id: "Song26", Style: "Electronic"}, {Id: "Song27", Style: "Reggae"},
// 			{Id: "Song28", Style: "Country"}, {Id: "Song29", Style: "Blues"}, {Id: "Song30", Style: "Metal"},
// 			{Id: "Song31", Style: "Rock"}, {Id: "Song32", Style: "Pop"}, {Id: "Song33", Style: "Jazz"},
// 			{Id: "Song34", Style: "Hip-Hop"}, {Id: "Song35", Style: "Classical"}, {Id: "Song36", Style: "Electronic"},
// 			{Id: "Song37", Style: "Reggae"}, {Id: "Song38", Style: "Country"}, {Id: "Song39", Style: "Blues"},
// 			{Id: "Song40", Style: "Metal"}, {Id: "Song41", Style: "Rock"}, {Id: "Song42", Style: "Pop"},
// 			{Id: "Song43", Style: "Jazz"}, {Id: "Song44", Style: "Hip-Hop"}, {Id: "Song45", Style: "Classical"},
// 			{Id: "Song46", Style: "Electronic"}, {Id: "Song47", Style: "Reggae"}, {Id: "Song48", Style: "Country"},
// 			{Id: "Song49", Style: "Blues"}, {Id: "Song50", Style: "Metal"},
// 		},
// 		Random20IndieSongs: []domain.Song{
// 			{Id: "IndieSong1", Style: "Indie Rock"}, {Id: "IndieSong2", Style: "Indie Pop"}, {Id: "IndieSong3", Style: "Indie Folk"},
// 			{Id: "IndieSong4", Style: "Indie Rock"}, {Id: "IndieSong5", Style: "Indie Pop"}, {Id: "IndieSong6", Style: "Indie Folk"},
// 			{Id: "IndieSong7", Style: "Indie Rock"}, {Id: "IndieSong8", Style: "Indie Pop"}, {Id: "IndieSong9", Style: "Indie Folk"},
// 			{Id: "IndieSong10", Style: "Indie Rock"}, {Id: "IndieSong11", Style: "Indie Pop"}, {Id: "IndieSong12", Style: "Indie Folk"},
// 			{Id: "IndieSong13", Style: "Indie Rock"}, {Id: "IndieSong14", Style: "Indie Pop"}, {Id: "IndieSong15", Style: "Indie Folk"},
// 			{Id: "IndieSong16", Style: "Indie Rock"}, {Id: "IndieSong17", Style: "Indie Pop"}, {Id: "IndieSong18", Style: "Indie Folk"},
// 			{Id: "IndieSong19", Style: "Indie Rock"}, {Id: "IndieSong20", Style: "Indie Pop"},
// 		},
// 	}

// 	algo := algoritmo.Algo{}

// 	recommendedSongs := algo.Algoritmo(mockData)
// 	return c.JSON(http.StatusOK, recommendedSongs)
// }

// func getRecommendedSongs(c echo.Context) error {
// 	params := new(domain.UserMusicPreferences)
// 	if err := c.Bind(params); err != nil {
// 		return err
// 	}

// 	algo := algoritmo.Algo{}

// 	recommendedSongs := algo.Algoritmo(*params)
// 	return c.JSON(http.StatusOK, recommendedSongs)
// }

func main() {
	// e := echo.New()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	err := godotenv.Load()
	if err != nil {
		fmt.Printf(err.Error())
	}

	// e.GET("/test-recommended-songs", fmt.Println("fodase"))
	// e.GET("/recommended-songs/id", )
	uri := os.Getenv("URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	db := database.NewDatabase(client, cancel)
	db.Ping(ctx)
	db.GetAllUserStyles(ctx, "66c868e082cd6161c37c0f48")
	db.Random50Songs(ctx)

	// e.Logger.Fatal(e.Start(":8080"))
}
