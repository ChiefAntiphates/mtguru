package main

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate"

	"mtguru/packages/config"
	"mtguru/packages/custom_logger"
)

var activeConfig config.EnvironmentConfig
var client *weaviate.Client

func init() {
	custom_logger.CreateLogger()
	activeConfig = config.CreateConfig()
	// client = createClient(activeConfig)
}

func main() {

	// createIndex(client) // currently done manually
	populateIndex()
	// searchDatabase(client)
	// updateCollection(client)
}
