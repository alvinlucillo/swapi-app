import { defineStore } from 'pinia'
import { reactive, readonly, ref } from 'vue'
import { gql } from 'graphql-tag'
import { useApolloClient, provideApolloClient } from '@vue/apollo-composable'
import { apolloClient } from './apollo'

export interface CharactersResult {
  getCharacters: {
    Characters: {
      name: string
      films: string[]
      vehicleModels: string[]
    }[]
    SearchID: string
  }
}

export const useSwapiStore = defineStore('swapi', () => {
  const isLoading = ref(false)

  const getCharacters = async (name: string): Promise<CharactersResult> => {
    isLoading.value = true
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
    `
    const { data } = await apolloClient.query({
      query: query,
      variables: {
        name: name
      }
    })
    isLoading.value = false
    return data
  }

  return {
    getCharacters,
    isLoading
  }
})
