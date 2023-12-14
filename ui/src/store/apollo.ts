import { ApolloClient, InMemoryCache } from "@apollo/client/core";

// Create the apollo client for querying the graphql server
export const apolloClient = new ApolloClient({
  uri: process.env.GRAPHQL_URI || "http://localhost:8080/graphql",
  cache: new InMemoryCache(),
});
