# Go Assignment Partner Matching

## Setup (macOS)

    > brew bundle
    > go build -o matching && ./matching

## Show Seed Data

    > open http://localhost:8080/partners

## Execute Example Query

    > PARM=$(echo -n '{"address_lon": 13.45, "address_lat": 52.5, "services": ["tiles", "carpet"], "floor_size": 120.5, "phone_number": "123-456-7890"}' | jq -s -R -r @uri)
    > curl "http://localhost:8080/matches/flooring?q=$(echo -n $PARM)"

## Open Swagger Docs

    > open "http://localhost:8080/swagger/index.html#/matches/get_matches_flooring"

