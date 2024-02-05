# Go Assignment Partner Matching

## Notes and Comments

1. In this implementation, the initial SQL request for Partners does not account
for filtering based on the distance from the request. This is applied later by
using the haversine function. In some instances, this approach may result in fewer
matches since the initial database request is limited for security reasons, and too
many Partners could be filtered out afterward. Therefore, integrating the distance
calculation directly into the database (e.g., as a stored procedure) is preferable.

2. The final sorting is implemented in a binary fashion according to your feature
description. However, since the Partner scores are represented as floats, the
secondary sorting by distance would never be applied. Therefore, we could either
round the score values or opt for a weighted sorting function that considers both
score and distance with varying weights.

3. The input parameters have not been thoroughly tested. Only the services are
verified to exclusively include the predefined terms (wood, tiles, carpet) as
an example (by the function validateServices).

4. Access control is not implemented (i.e. http-Auth or the like).

5. The Swagger input and output examples are not well-defined.

## Setup (macOS)

    > brew bundle
    > go build -o matching && ./matching

## Run Tests

    > go test
      PASS
      ok      matching        0.383s

## Show Seed Data

    > open http://localhost:8080/partners

![Partner Index](static/partner-index.png)

## Execute Example Query

![Postman Request](static/postman-request.png)

## Open Swagger Docs

    > open "http://localhost:8080/swagger/index.html#/matches/post_matches_flooring"

![Swagger Docs](static/swagger-docs.png)
