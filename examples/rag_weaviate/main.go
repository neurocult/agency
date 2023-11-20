package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func main() {
	openAPIKey := os.Getenv("OPENAI_API_KEY")

	ctx := context.Background()

	client, err := prepareDB(openAPIKey, ctx)
	if err != nil {
		panic(err)
	}

	fields := []graphql.Field{{Name: "content"}}

	result, err := client.GraphQL().Get().
		WithClassName("Records").
		WithFields(fields...).
		WithNearText(
			client.GraphQL().
				NearTextArgBuilder().
				WithConcepts([]string{"programming"}),
		).
		WithLimit(10).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	// factory := openai.New(openAPIKey)
	// pipeline.New(
	// 	factory.TextToText(openai.TextToTextParams{
	// 		Model: "gpt3dot5turbo",
	// 	}).WithOptions(core.WithPrompt("")),
	// )

	// TODO add TTS example with pipes

	fmt.Println("query about programming", result)
}

func prepareDB(openAPIKey string, ctx context.Context) (*weaviate.Client, error) {
	client, err := weaviate.NewClient(weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
		Headers: map[string]string{
			"X-OpenAI-Api-Key": openAPIKey,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := client.Schema().AllDeleter().Do(ctx); err != nil {
		return nil, err
	}

	classObj := &models.Class{
		Class:      "Records",
		Vectorizer: "text2vec-openai",
		ModuleConfig: map[string]interface{}{
			"text2vec-openai":   map[string]interface{}{},
			"generative-openai": map[string]interface{}{},
		},
	}
	if err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background()); err != nil {
		return nil, err
	}

	if _, err := client.Batch().ObjectsBatcher().WithObjects(data...).Do(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
