package main

import (
	"brimobile/app/account/repository"
	"brimobile/app/account/service"
	repository2 "brimobile/app/saving/repository"
	savingService "brimobile/app/saving/service"
	"brimobile/db/connection"
	"brimobile/graph"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("cant load env file :" + err.Error())
		return
	} else {
		fmt.Println("success load env file")
	}
}

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	connection.ConnectDB()
	defer connection.DB.Close()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		AccService:    service.NewAccountService(repository.NewAccountRepository(connection.DB)),
		SavingService: savingService.NewSavingService(repository2.NewSavingRepository(connection.DB)),
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err.Error())
	}
}
