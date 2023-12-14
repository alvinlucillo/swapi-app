import { ApolloClient, InMemoryCache } from '@apollo/client/core'

export const apolloClient = new ApolloClient({
  uri: 'http://localhost:8080/graphql',
  cache: new InMemoryCache()
})

// Configure the Apollo client
// const cache = new InMemoryCache()

// const apolloClient = new ApolloClient({
//   cache,
//   uri: 'http://localhost:8080/graphql'
// })
