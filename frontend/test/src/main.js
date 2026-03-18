import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { userStoreKey, userStore } from './store'

console.log('Starting app initialization')
console.log('App component:', App)
console.log('Router:', router)

const app = createApp(App)
console.log('App created:', app)

app.use(router)
app.provide(userStoreKey, userStore)
console.log('Router added to app')

app.mount('#app')
console.log('App mounted to #app')
