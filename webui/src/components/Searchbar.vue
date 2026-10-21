<script>
import "../assets/search-bar.css";
export default{
	props:{selectmode:{type: Boolean, default: false,},},
	emits: ['select-user'],

    data: function(){
        return {
            searchkey: "",
            searchlist: [],
            searchin: false,
            searchTimeout:null,
			errormsg: null,
			userID: sessionStorage.getItem("userID"),

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

	methods:{
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

		async makeChat(user){
			this.errormsg=null;
			try{
				const response= await this.$axios.post('/chats',{
					chatname: "",
					photo: null,
					members: [user.userid],
					grouporchat: false,
				});
				console.log("makechat:\n", response.data);
				const chatID=response.data.chatid;
				if (!chatID){
					this.errormsg="server did not return chat ID";
					return;
				} else{
					this.searchkey="";
					this.searchlist=[];
					await this.$router.replace(`/chats/${chatID}/messages`);
				}
				}catch(e){
					console.error("failed to make or open chat: ", e);
					this.errormsg=e.response?.data?.error || e.response?.data || "couldn't create/open chat :(";	
				}
		},
		async selectuser(user){
			if (this.selectmode){
				this.$emit('select-user', user);
			} else{
				await this.makeChat(user);
			}
			this.searchkey='';
			this.searchlist=[];
		},
	},
	beforeUnmount(){clearTimeout(this.searchTimeout);},	
}
</script>


<template>
    <div class="search-bar">
        <input v-model="searchkey" type="text" class="form-control" placeholder="search users...">
        <div v-if="searchkey.trim()" class="search-list">
            <button v-for="user in searchlist" :key="user.userid" class="itemmm" v-on:click="selectuser(user)">
				<img v-if="user.photo" :src="`${this.$axios.defaults.baseURL}/${user.photo}`" alt="pfp" v-on:error="$event.target.style.display='none'"  class="item-pic"/>
                <p class="item-text">{{ user.username }}</p>
			</button>
        </div>
        <p v-if="!searchin && searchkey.trim() && searchlist.length===0">
            No users found. Try another searchkey...
        </p>
		<p v-if="errormsg">{{ errormsg }}</p>
    </div>
</template>