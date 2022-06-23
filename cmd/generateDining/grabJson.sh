#!/bin/sh

curl --data '{}' --cookie-jar cookies.txt --user-agent 'chrome/102.0.0.0' 'https://disneyworld.disney.go.com/finder/api/v1/authz/public'
curl --cookie cookies.txt --user-agent 'chrome/102.0.0.0' 'https://disneyworld.disney.go.com/finder/api/v1/explorer-service/list-ancestor-entities/wdw/80007798;entityType=destination/2022-06-22/dining' --compressed | jq .results

