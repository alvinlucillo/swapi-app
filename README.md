# Star Wars Trivia App

As a Star Wars super fan that loves trivia, the goal is to develop a UI that allows users to list the film(s) and vehicle model(s) for Star Wars characters, aiding in preparation for trivia nights.

#### Acceptance Criteria

- The UI should be a locally hosted web based application
  - Use any technology stack you are comfortable with – the choice is yours
  - A text box should be used for text entry of the person
  - The films and vehicle model should be returned in a table – you decide how to
structure the table
  - The UI should have the ability to store search results
  - The UI should have the ability to retrieve search results that are saved
- The Star Wars API should be used as your data source: https://swapi.dev/
  - HINT: Make sure you read the documentation to find the search function
- The interface used by the UI for querying and mutating data should be GraphQL based
  - Use any GraphQL server technology/tooling you are comfortable with
- The storage of search results should be in a database
  - Use any database you are comfortable with and more importantly, one you feel
will address the problem best
- (Bonus - Optional) If multiple people are retuned in a search result, list out the results in
separate tables for each person found

#### Test Case

Searching “Darth Maul” should list “The Phantom Menace” and the vehicle model “FC-20 speeder bike”.

## Background
### Backend
- The backend uses Go for the GraphQL server and MongoDB for the database. 
- The GraphQL server takes these queries and mutation:
  - `getCharacters`: returns a list of characters based on the search term
    - Every time a search is made, a search is created in the database with expiration (TTL)
        - If the user saves the search via `saveSearch`, the expiration is removed
    - Every time a search is made, characters, films, and vehicles are saved in the database with TTL based on environment variable DOCUMENT_TTL so future queries using the same objects will be faster. 
       - If objects don't exist in the database, they are fetched from the Star Wars API. This happens if they don't exist in the first place or if they have expired.
  - `getSavedSearches`: returns a list of saved searches
  - `getSavedSearchesByIDs`: returns the characters based on the IDs
  - `saveSearch`: saves a search to the database
    - Accepts the search ID created by `getCharacters`. This is used to find the search in the database.

## Getting started
### Run locally via docker-compose
- Run `docker-compose up` to start the ui, server, and database
- Access the UI at http://localhost:8081/
- Access the GraphQL playground at http://localhost:8080/graphql
- Connect to the database at mongodb://localhost:27017
- Run `docker-compose down` to stop the containers
- Run `docker-compose up --build` to rebuild the containers if you make changes to the code

### Run locally without docker-compose
- Run `docker run --name mongodb -p 27017:27017 -d mongo` to start the database
- Run `cd server && go run ./cmd/` to start the GraphQL server
- Run `cd ui && npm install && npm run dev` to start the UI