import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { DefaultApolloClient } from '@vue/apollo-composable'
import { createApolloProvider } from '@vue/apollo-option'

import { apolloClient } from './stores/apollo'

import App from './App.vue'
import router from './router'
import { Quasar } from 'quasar'
import quasarUserOptions from './quasar-user-options'

// // Configure the Apollo client
// const cache = new InMemoryCache()

// const apolloClient = new ApolloClient({
//   cache,
//   uri: 'http://localhost:8080/graphql'
// })

const app = createApp(App).use(Quasar, quasarUserOptions)

const apolloProvider = createApolloProvider({
  defaultClient: apolloClient
})

app.use(apolloProvider)
app.use(createPinia())
app.use(router)

app.mount('#app')
