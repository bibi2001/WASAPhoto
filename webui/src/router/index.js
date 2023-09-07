import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import SearchView from '../views/SearchView.vue'
import LoginView from '../views/LoginView.vue'
import UploadView from '../views/UploadView.vue'
import ProfileView from '../views/ProfileView.vue'
import Followers from '../views/Followers.vue'
import Following from '../views/Following.vue'
import Bans from '../views/Bans.vue'
import ChangeUsername from '../views/ChangeUsername.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/home', component: HomeView},
		{path: '/search', component: SearchView},
		{path: '/new', component: UploadView},
		{path: '/profile/:username', component: ProfileView},
		{path: '/profile/:username/username', component: ChangeUsername},
		{path: '/profile/:username/followers', component: Followers},
		{path: '/profile/:username/following', component: Following},
		{path: '/profile/:username/bans', component: Bans},
	]
})

export default router
