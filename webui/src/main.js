import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import UserList from './components/UserList.vue'
import Photo from './components/Photo.vue'
import Comment from './components/Comment.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("UserList", UserList);
app.component("Photo", Photo);
app.component("Comment", Comment);
app.use(router)
app.mount('#app')
