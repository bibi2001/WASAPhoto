import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/search', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/user/:id/link', component: HomeView},
	]
})

export default router
