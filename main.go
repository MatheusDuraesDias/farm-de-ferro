package main

import (
	"algorithm/mod/algoritmo"
	"algorithm/mod/algoritmo/database"
	"algorithm/mod/handler"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Me and my froggas")
}

func main() {
	e := echo.New()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err.Error())
	}

	mongoUri := os.Getenv("URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	neo4JUri := os.Getenv("NEO4J_URI")
	neo4JUsername := os.Getenv("NEO4J_USERNAME")
	neo4JPass := os.Getenv("NEO4J_PASSWORD")
	driver, err := neo4j.NewDriverWithContext(neo4JUri, neo4j.BasicAuth(neo4JUsername, neo4JPass, ""))
	if err != nil {
		panic(err)
	}

	db := database.NewDatabase(client, cancel)
	neo4JDb := database.NeoDatabase{
		Driver: driver,
	}

	algo := algoritmo.Algo{
		Db:      *db,
		Neo4JDb: neo4JDb,
		Ctx:     ctx,
	}

	h := handler.Handler{
		Algo: algo,
	}

	e.GET("recommended-songs/:userId", h.GetRecommendedSongs())
	e.GET("health", healthHandler)
	e.Logger.Fatal(e.Start(":8080"))

}
