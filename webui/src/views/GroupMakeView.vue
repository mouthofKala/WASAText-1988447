<script>
import "../assets/group-make-page.css";
export default {
    data: function(){
        return {
            groupName:'',
            groupPhoto: null,
            photoPreview: null,
            users:[],
            selectedUsers:[],
            searchQuery: '',
            loading: false,
            errorMsg:null,
        };
    },
    computed: {
        filteredUsers(){
            const query = this.searchQuery.trim().toLowerCase();
            return this.users.filter(user=>{
                const id= user.userid;
                const name=user.username;
                return !this.selectedUsers.some(
                    selected=>(selected.userid)===id) && (!query ||name?.toLowerCase().includes(query));
            });
        },
    },
    methods:{
        async getUsers(){
            try{
                const response = await this.$axios.get('/users');
                console.log("getusers:\n",response.data)
                this.users=response.data;
            } catch(e){
                this.errorMsg = "couldn't load users...";
            this.users=[]}
        },
        async selectPhoto(event){
            const file = event.target.files?.[0];
            if (!file) return;
            if (!file.type.startsWith('image/')){
                this.errorMsg='file must be an image!';
                return;
            }
            this.groupPhoto=file;
            this.photoPreview=URL.createObjectURL(file);
            this.errorMsg=null;
        },
        getUserID(user){
            return user.userid; 
        },
        getPhoto(user){
            return user.photo;
        },
        getUsername(user){
            return user.username;
        },
        toggleUser(user){
            const id = this.getUserID(user);
            const idx = this.selectedUsers.findIndex(selected => this.getUserID(selected) === id);
            if (idx!== -1){
                this.selectedUsers.splice(idx, 1);
            } else {
                this.selectedUsers.push(user)
            }
        },
        async makeChatGroup(){
            if (!this.groupName.trim()){
                this.errorMsg='please enter a name for the group';
                return;
            }
            if (this.selectedUsers.length===0){
                this.errorMsg='please select at least one user';
                return;
            }

            this.loading=true;
            this.errorMsg=null;
            try{
                let photo = '';
                if (this.groupPhoto){
                    photo = await new Promise((resolve, reject) =>{
                        const reader = new FileReader();
                        reader.onload = () => {resolve(reader.result.split(',')[1]);};
                        reader.onerror=reject;
                        reader.readAsDataURL(this.groupPhoto);
                    });
                }
                const response=await this.$axios.post('/chats',{
                    chatname: this.groupName.trim(),
                    photo,
                    members: this.selectedUsers.map(user=>this.getUserID(user)),
                    grouporchat: true,
                });
                console.log("makegroup:\n", response.data)
                const chatID=response.data.chatid ?? response.data.chat?.chatid;
                if (chatID){
                    await this.$router.replace(`/chats/${chatID}`);
                } else{
                    await this.$router.replace(`/chats/`);
                }
            }catch(e){
                console.error("failed to make group: ", e);
                this.errorMsg=e.response?.data?.error || e.response?.data || "couldn't create group :(";
            } finally{
                this.loading=false;
            }
        },


    },
    beforeUnmount(){
        if (this.photoPreview){
            URL.revokeObjectURL(this.photoPreview);
        }
    },
    async mounted(){
        await this.getUsers();
    },
};
</script>

<template>
    <div class="group-make-page">
        <div class="group-info">
            <button @click="this.$router.push('/chats')" class="back">↩</button>
            <label class="photo-picker">
                <img v-if="photoPreview" :src="photoPreview" class="photo-picker img" @click="selectPhoto" />
                <span v-else>＋</span>
                <input type="file" accept="image/*" hidden @change="selectPhoto"/>
            </label>
            
            <input v-model="groupName" class="group-name" placeholder="group name..." maxlength="40"/>
            <button class="create-button" :disabled="loading" @click="makeChatGroup" aria-label="create group">➜</button>
        </div>

        <div class="search-bar">
            <Searchbar :select-mode="true" @select-user="toggleUser"/> <!-- modify inner (user) buttons-->
        </div>
        <ErrorMsg v-if="errorMsg"/>

        <div class="members-n-added">
            <section class="listofusers">
                <h3>select users...</h3>
                <button v-for="user in filteredUsers" :key="getUserID(user)" class="member-item" @click="toggleUser(user)">
                    <img v-if="getPhoto(user)" :src="`${$axios.defaults.baseURL}/${getPhoto(user)}`" alt="" @error="$event.target.style.display='none'"/>
                    <span v-else class="pfp-placeholder"> </span>

                    <span class="username">{{ getUsername(user) }}</span>
                    <span class="circle-select"> </span>
                </button>
                <!--member item: pfp and username. possibly sorted by most recent msg
                the item has a circle button which, if you press on the item in general
                gets coloured in and moves the item to listofselected-->
            </section>
            <section class="listofselected">
                <h3>Selected:</h3>
                <button v-for="user in selectedUsers" :key="getUserID(user)" class="member-item selected" @click="toggleUser(user)">
                    <img v-if="getPhoto(user)" :src="`${$axios.defaults.baseURL}/${getPhoto(user)}`" alt="" @error="$event.target.style.display='none'"/>
                    <span v-else class="pfp-placeholder"> </span>

                    <span class="username">{{ getUsername(user) }}</span>
                    <span class="circle-select"> </span>
                </button>
                <!-- list of members selected from listofusers, if you click on it
                they get deselected so they go back to the list on the left -->
            </section>
        </div>
    </div>
</template>


<!--1/3 horizontally: pfp 2/3 horizontally: group name
occupies max 1/5 of the screen vertically
-->

