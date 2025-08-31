// server.go

package main

import (

"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {

// ...
err := godotenv.Load(); if err != nil {
log.Fatal("Error loading .env file")
}

// ...

Database := graph.Connect()
srv := handler.NewDefaultServer(
generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: Database}}))
// ...
}