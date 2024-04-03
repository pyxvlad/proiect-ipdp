#!/usr/bin/env bash

templ generate
npx tailwindcss -i css/tw.css -o css/style.css

echo "Starting server"
export IPDP_DB=appdata/ipdp.db
go run .

