<script>
import ProfileButton from '../components/ProfileButton.vue';
import ErrorMsg from '../components/ErrorMsg.vue';
import LoadingSpinner from '../components/LoadingSpinner.vue';
import Searchbar from '../components/Searchbar.vue';

export default {
	components: {Searchbar},
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
		}
	},
	watch: {
		searchkey(){
			clearTimeout(this.searchTimeout);
			this.searchTimeout=setTimeout(()=>{
				this.searchUsers();
			},300); // wrapped searchUsers with a timeout
		}
	},
	methods: {
		async getMyConversations() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/chats");
				this.chatlist = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		togglebuttns(){
			this.showbuttns = !this.showbuttns;
		},
		gomakegroupview(){
			// might collide with chatid!! so, will go back and then push newgroup
			this.$router.back();
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
				this.searchlist=response.data || [];
			} catch(e){
				this.errormsg=e.toString();
			} finally {
				this.searchin=false;
			}
		},
		logout(){
			localStorage.removeItem("userID")
			delete this.$axios.defaults.headers.common.Authorization;
			this.$router.push("/session");
		}
	},

	mounted() {
		this.getMyConversations()
	},
	beforeUnmount(){clearTimeout(this.searchTimeout);},	
}
</script>

<template>
	<template v-if="errormsg">
		<ErrorMsg :msg="errormsg"></ErrorMsg>
	</template>
	<template v-else>
		<div class="home-page">
			<div class="home-header">
				<h1 class="h2">Conversations</h1>
				<Searchbar />
				<div class="search-bar">
					<input v-model="searchkey" type="text" class="form-control" placeholder="search users...">
					<div v-if="searchkey.trim()" class="search-list">
						<button v-for="user in searchlist" :key="user.userID" class="search-res" @click="$router.push(`/users/${user.userID}`)">
							{{ user.username }}
						</button>
					</div>
					<p v-if="!searchin && searchlist.length===0">
						No users found. Try another search key...
					</p>
					<!--...small defect: this setup does not search groups. will see later how to fix-->
				</div>
				<button class="menu-button" @click="togglebuttns">⋮</button>
					<div v-if="showbuttns" class="smolmenu">
						<button @click="gomakegroupview">New group</button>
						<button @click="logout" class="logout-buttn">Log out</button>
					</div>

				<ProfileButton />
			</div>
			<ErrorMsg v-if="errormsg" :msg="errormsg" />
			<LoadingSpinner v-if="loading" />
			<div v-else-if="chatlist.length === 0">
				<p>you have no chats :[ </p>
			</div>
			<div v-else class="chat-list">
				<ChatItem
					v-for="chat in chatlist" :key="chat.chatid" :chat="chat" />
			</div>
		</div>
	</template>
</template>

<style scoped>
	.home-page{
		height: 100vh;
		display: flex;
		flex-direction: column;
	}

	.home-header{
		flex-shrink: 0;

	}

	.chat-list{
		flex: 1;
		overflow-y: auto;
	}

	.search-bar{
		position: relative;
	}

	.search-list{
		position: absolute;
		top: 100%;
		left: 0;
		right: 0;
		max-height:33vh;
		overflow-y: auto;
		background: white;
		border: 1px solid whitesmoke;
		z-index: 10;
	}

	.search-res{
		display: block;
		width: 100%;
		padding: 10px;
		text-align: left;
		border: none;
		background: white;
	}

	.search-res:hover{
		background: whitesmoke;
	}
</style>
