# mtguru

# Project folder structure 

- apps:     user facing apps/websites
- services: for backend services 
- packages: for packages/modules 


# Data sources

MTG card dump downloaded from [here](https://scryfall.com/docs/api/bulk-data) (default cards)

# How to run the search stuff

In `services/ingestion/main.go`, moddify the search_string on line 50

Needs to run using: `go run .\services\ingestion\` from the root of the mtguru repo