#!/bin/sh

go-jsonschema album.json collection.json common.json --package github.com/albums2/gopkg/types/schema --output=/workspaces/gopkg/filetypes/schema/schema.go
bunx -p json-schema-to-typescript json2ts -i '*.json' -o /workspaces/frontend/src/lib/server/schema/