package main

import (
	"mtguru/packages/config"
	"mtguru/packages/custom_logger"
)

func init() {
	custom_logger.CreateLogger()
	config.CreateConfig()
}

func main() {
	populateIndex()
}
