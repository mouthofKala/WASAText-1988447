import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ChatView from '../views/ChatView.vue'
import ProfileView from '../views/ProfileView.vue'
import LoginView from '../views/LoginView.vue'
import GroupMakeView from '../views/GroupMakeView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/chats', component: HomeView},
		{path: '/chats/:chatid', component: ChatView},
		{path:'/chats/newgroup', component: GroupMakeView},	
		{path:'/chats/:chatid/settings', component: ProfileView},	
		{path: '/users', component: HomeView},
		{path: '/profile', component: ProfileView},
		{path: '/users/:userid', component: ProfileView},
		{path: '/session', component: LoginView},

	]
})

export default router
