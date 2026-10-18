<script>
export default {
    props: {
        chat: {
            type: Object,
            required: true
        }
    },
    data(){
        return{
            userID: sessionStorage.getItem("userID"),
            otheruser:null,
        };
    },
    methods: {
        getConversation(){
            this.$router.push(`/chats/${this.chat.chatid}/messages`)
        },
        async getOtherUser(){
            if (this.chat.grouporchat || !this.chat.members) return;
            const otherID= this.chat.members.find(id=>id!==this.userID);
            if (!otherID) return;
            try{const response = await this.$axios.get(`/users/${otherID}`);
                this.otheruser=response.data;
            } catch(e){console.error("failed to dfind otheruser:",e);}
        } 
    },
    mounted(){this.getOtherUser();}
}
</script>

<template>
    <div class="chat-item border rounded p-3 mb-2" @click="getConversation">
        <pre>{{ JSON.stringify(chat, null, 2) }}</pre>
        <img v-if="chat.grouporchat ? chat.photo : otheruser?.photo"
            alt="chat pfp"
            :src="`${$axios.defaults.baseURL}/${(chat.grouporchat ? chat.photo : otheruser.photo).replace(/^\/+/, '')}`"
            class="item-photo"
            @error="$event.target.src = `${$axios.defaults.baseURL}/storage/data/pfp/black.png`" />
        <h5 v-if="chat.grouporchat">{{ chat.chatname }}</h5>
        <h5 v-else>{{ otheruser?.username || 'user loading...'}}</h5>
        <div v-if="chat.mostRecentMsg">
            <p class="mb-1">
                {{ chat.mostRecentMsg.username }}:
            </p>
            <p v-if="chat.mostRecentMsg.content && !chat.mostRecentMsg.photo" class="msg-prev">
                {{ chat.mostRecentMsg.content || 'msg loading...'}}
            </p>
            <p v-else-if="chat.mostRecentMsg.photo && !chat.mostRecentMsg.content">
                📷 Photo
            </p>
            <p v-else-if="chat.mostRecentMsg.content && chat.mostRecentMsg.photo" class="msg-prev">
                📷 {{ chat.mostRecentMsg.content || 'msg loading...'}}
            </p>
            <p v-else>error with message preview...</p>
        </div>
        <p v-else></p>
    </div>
</template>

<style scoped>
    .chat-item {
        cursor:pointer;

    }
    .msg-prev{
        display: block;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }
    .item-photo{
        height:55px;
        width:55px;
        border-radius:50%;
        overflow: hidden;
        border:none;
        padding:0;

    }
</style>