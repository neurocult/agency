package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

// TODO 1) add pipelines 2) add "with distance" examples
func main() {
	openAPIKey := os.Getenv("OPENAI_API_KEY")

	client, err := weaviate.NewClient(weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
		Headers: map[string]string{
			"X-OpenAI-Api-Key": openAPIKey,
		},
	})
	if err != nil {
		panic(err)
	}

	// add the schema
	classObj := &models.Class{
		Class:      "Records",
		Vectorizer: "text2vec-openai",
		ModuleConfig: map[string]interface{}{
			"text2vec-openai":   map[string]interface{}{},
			"generative-openai": map[string]interface{}{},
		},
	}
	if err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background()); err != nil {
		panic(err)
	}

	ctx := context.Background()

	// insert and vectorize via openai
	resp, err := client.Batch().ObjectsBatcher().WithObjects(data...).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("insertion", resp)

	// query
	fields := []graphql.Field{{Name: "content"}}

	nearText := client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{"programming"})

	result, err := client.GraphQL().Get().
		WithClassName("Records").
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(5).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("query about programming", result)
}

type ContentItem struct {
	Content string `json:"content"`
}

func readAndParseJSON(path string) []ContentItem {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var items []ContentItem

	err = json.Unmarshal(byteValue, &items)
	if err != nil {
		panic(err)
	}

	return items
}
