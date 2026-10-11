#!/bin/sh

go-jsonschema album.json collection.json common.json --package github.com/albums2/gopkg/types/schema --output=/workspaces/gopkg/filetypes/schema/schema.go