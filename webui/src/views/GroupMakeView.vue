<script>
import Searchbar from '../components/Searchbar.vue';
import ErrorMsg from '../components/ErrorMsg.vue';
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
    components: {Searchbar, ErrorMsg},
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
                this.users=response.data;
            } catch(e){this.errorMsg = "couldn't load users...";}
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
                const chatID=response.data.chatid ?? response.data.chat?.chatid;
                if (chatID){
                    await this.$router.back();
                    await this.$router.push(`/chats/${chatID}`)
                } else{
                    await this.$router.back();
                }
            }catch(e){
                this.errorMsg=error.response?.data || "couldn't create group :(";
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
    <template v-if="errorMsg">
        <ErrorMsg :msg="errorMsg" />
    </template>
    <template v-else>
        <div class="group-make-page">
            <div class="group-info">
                <button @click="this.$router.back()">↩</button>
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
                        <img v-if="getPhoto(user)" :src="getPhoto(user)" alt="" />
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
                        <img v-if="getPhoto(user)" :src="getPhoto(user)" alt="" />
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
</template>


<!--1/3 horizontally: pfp 2/3 horizontally: group name
occupies max 1/5 of the screen vertically
-->
<style scoped>
.group-make-page{
    display: flex;
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
    height:100%;
}
.group-info{
    min-height:90px;
    height:33vh;
    align-items:center;
    display: flex;
    gap: 1rem;
}
.photo-picker{
    border-radius:50%;
    width:64px;
    height:64px;
    overflow:hidden;
    display:grid;
    place-items:center;
    background: gainsboro;
    cursor:pointer;
    font-size: 1.8rem;
    flex-shrink:0;
}
.photo-picker img, member-item img{
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.group-name{
    flex: 1;
    min-width: 0;
    padding: 0.7rem;
    border: none;
    border-bottom: 1px solid lightgray;
    font-size: 1.1rem;
}
.create-button{
    border:none;
    border-radius:50%;
    font-size:1.4rem;
    cursor: pointer;
    width: 42px;
    height:42px;
}

.members-n-added{
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    flex:1;
    min-height:0;
}
.listofusers, .listofselected{
    overflow-y: auto;
    min-width: 0;
}
.member-item{
    display:flex;
    align-items: center;
    gap: 0.6rem;
    border-bottom: 1px solid gainsboro;
    border-radius:8px;
    background: transparent;
    cursor: pointer;
    text-align: left;
    width: 100%;
    padding: 0.6rem;
}
.member-item.selected, member-item:hover{
    background:aliceblue;
}
.pfp-placeholder, .member-item img{
    width: 38px;
    height: 38px;
    border-radius: 50%;
    flex-shrink: 0;
}
.pfp-placeholder{
    display: grid;
    place-items: center;
    background: lightgray;
}
.username{
    overflow: hidden;
    text-overflow: ellipsis;
}
.circle-select{
    margin-left: auto;
    flex-shrink: 0;
}
@media (max0width: 600px){
    .members-n-added{
        grid-template-columns: 1fr 1fr;
        gap: 0.4rem;
    }
    .member-item{
        gap: 0.3rem;
        padding: 0.4rem 0.2rem;
    }
}
</style>

