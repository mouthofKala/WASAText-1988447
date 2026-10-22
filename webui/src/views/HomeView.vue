<script>
import "../assets/home-page.css";
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			chatlist: [],
			showbuttns: false,
			searchkey: "",
			searchlist:[],
			searchin: false,
			searchTimeout:null,
			refresh: false,
			init: true,
			refreshing: null,
		}
	},
	watch: {
		searchkey(){
			clearTimeout(this.searchTimeout);
			this.searchTimeout=setTimeout(()=>{
				this.searchUsers();
			},300);},
		'$route.path'(newpath){if (newpath==='/chats'){this.getMyProfile();}},
	},
	methods: {
		async getMyConversations() {
			if (this.refresh) return;
			this.refresh = true;
			if (this.init){this.loading=true;}
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/chats");
				console.log("getmyconversations: \n",response.data)
				this.chatlist = Array.isArray(response.data) ? response.data : [];
			} catch (e) {
				this.errormsg = e.response?.data?.error || e.message || "could not load chats :(";
				this.chatlist=[];
				if (this.init){this.chatlist=[];}
			} finally {
				this.refresh=false;
				this.loading = false;
				this.init=false;
			}
		},
		togglebuttns(){
			this.showbuttns = !this.showbuttns;
		},
		gomakegroupview(){
			this.$router.push("/newgroup"); 
			
		},
		async searchUsers(){
			const query = this.searchkey.trim();
			if (!query){
				this.searchlist=[];
				this.searchin=false;
				return;
			}

			this.searchin = true;
			this.errormsg=null;
			try{
				const response = await this.$axios.get("/users", {params: {searchkey: query}});
				console.log("searchusers:\n", response.data)
				this.searchlist=response.data || [];
			} catch (e) {
				this.errormsg = e.response?.data?.error || e.message || "could not load users...";
				this.searchlist=[];
			} finally {
				this.searchin = false;
			}
		},
		logout(){
			sessionStorage.removeItem("userID")
			delete this.$axios.defaults.headers.common.Authorization;
			this.$router.push("/");
		}
	},

	mounted() {
		this.getMyConversations();
		this.refreshing=setInterval(()=>{
			this.getMyConversations();}, 3000);
	},
	beforeUnmount(){clearTimeout(this.searchTimeout); clearInterval(this.refreshing);},	
}
</script>

<template>
	<div class="home-page">
		<div class="home-header">
			<h1 class="header-left">Conversations</h1>
			<Searchbar /> <!--...small defect: this setup does not search groups.-->
			<button class="menu-button" @click="togglebuttns">⋮
				<div v-if="showbuttns" class="smolmenu">
					<button @click="gomakegroupview" class="groupbuttn">New group</button>
					<button @click="logout" class="logout-buttn">Log out</button>
				</div>
			</button>
			<ProfileButton />
		</div>
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<LoadingSpinner v-if="loading" />
		<div v-else-if="chatlist.length === 0">
			<p>you have no chats :[ </p>
		</div>
		<div v-else class="chat-list">
			<ChatItem
				v-for="chat in chatlist" :key="chat.chatid" :chat="chat"></ChatItem>
		</div>
	</div>
</template>
