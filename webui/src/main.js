import {createApp, reactive} from 'vue';
import App from './App.vue';
import router from './router';
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue';
import LoadingSpinner from './components/LoadingSpinner.vue';
import ChatItem from './components/ChatItem.vue';
import MessageItem from './components/MessageItem.vue';
import './assets/dashboard.css';
import './assets/main.css';
import ProfileButton from './components/ProfileButton.vue';
import Searchbar from './components/Searchbar.vue';
import AddOrFwd from './components/AddOrFwd.vue';

const app = createApp(App);
const userID = sessionStorage.getItem("userID");

if (userID) {
    axios.defaults.headers.common.Authorization = `Bearer ${userID}`;
}

app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("ChatItem",ChatItem);
app.component("MessageItem",MessageItem);
app.component("ProfileButton",ProfileButton);
app.component("Searchbar",Searchbar);
app.component("AddOrFwd", AddOrFwd);
app.use(router);
app.mount('#app');
