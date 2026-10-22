<script>
import "../assets/group-make-page.css";
import groupmake from "../groupmake.js";
export default {
    mixins:[groupmake],
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
            </section>
            <section class="listofselected">
                <h3>Selected:</h3>
                <button v-for="user in selectedUsers" :key="getUserID(user)" class="member-item selected" @click="toggleUser(user)">
                    <img v-if="getPhoto(user)" :src="`${$axios.defaults.baseURL}/${getPhoto(user)}`" alt="" @error="$event.target.style.display='none'"/>
                    <span v-else class="pfp-placeholder"> </span>

                    <span class="username">{{ getUsername(user) }}</span>
                    <span class="circle-select"> </span>
                </button>
            </section>
        </div>
    </div>
</template>
