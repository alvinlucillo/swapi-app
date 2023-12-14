import { defineStore } from "pinia";
import { gql } from "graphql-tag";
import { apolloClient } from "./apollo";

export interface CharactersResult {
  getCharacters: {
    Characters: Character[];
    SearchID: string;
  };
}

export interface Character {
  films: string[];
  vehicleModels: string[];
  name: string;
}

export interface SavedSearchesByIDResult {
  getSavedSearchesByID: Character[];
}

export interface SavedSearch {
  ID: string;
  SearchKey: string;
}

export interface SavedSearchesResult {
  getSavedSearches: SavedSearch[];
}

export const useSwapiStore = defineStore("swapi", () => {
  const searchCharacters = async (name: string): Promise<CharactersResult> => {
    const query = gql`
      query GetCharacters($name: String!) {
        getCharacters(name: $name) {
          Characters {
            name
            films
            vehicleModels
          }
          SearchID
        }
      }
    `;
    const { data } = await apolloClient.query({
      query: query,
      variables: {
        name: name,
      },
    });
    return data;
  };

  const getSavedSearches = async (): Promise<SavedSearchesResult> => {
    const query = gql`
      query GetSavedSearches {
        getSavedSearches {
          ID
          SearchKey
        }
      }
    `;
    const { data } = await apolloClient.query({
      query: query,
      fetchPolicy: "network-only",
    });

    return data;
  };

  const getSavedSearchByID = async (
    searchID: string
  ): Promise<SavedSearchesByIDResult> => {
    const query = gql`
      query GetSavedSearchByID($searchID: String!) {
        getSavedSearchesByID(searchID: $searchID) {
          films
          vehicleModels
          name
        }
      }
    `;
    const { data } = await apolloClient.query({
      query: query,
      variables: {
        searchID: searchID,
      },
    });

    return data;
  };

  const saveSearch = async (searchID: string): Promise<CharactersResult> => {
    const query = gql`
      mutation SaveSearch($searchID: String!) {
        saveSearch(searchID: $searchID)
      }
    `;
    const { data } = await apolloClient.mutate({
      mutation: query,
      variables: {
        searchID: searchID,
      },
    });

    return data;
  };

  return {
    searchCharacters,
    getSavedSearches,
    getSavedSearchByID,
    saveSearch,
  };
});
