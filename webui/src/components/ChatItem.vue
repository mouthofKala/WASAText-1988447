<script>
import "../assets/chat-item.css";
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
    watch:{
        'chat.members'(){this.getOtherUser();},
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
    mounted(){this.getOtherUser();},
}
</script>

<template>
    <div class="chat-item" >
        <div class="img-box">
            <img v-if="chat.grouporchat ? chat.photo : otheruser?.photo"
                alt="chat pfp"
                :src="`${$axios.defaults.baseURL}/${(chat.grouporchat ? chat.photo : otheruser.photo).replace(/^\/+/, '')}`"
                class="item-photo"
                @error="$event.target.src = `${$axios.defaults.baseURL}/storage/data/pfp/black.png`" />            
        </div>
        <!-- if (!$event.target.src.endsWith('/black.png')) -->

        <div v-if="chat.mostrecentmsg" class="chat-box" @click="getConversation">
            <h5 class="chatname" v-if="chat.grouporchat">{{ chat.chatname }}</h5>
            <h5 class="chatname" v-else>{{ otheruser?.username || 'user loading...'}}</h5>  
            <p class="usrname">
                {{ chat.mostrecentmsg.username }}:
            </p>
            <p v-if="chat.mostrecentmsg.content && !chat.mostrecentmsg.photo" class="msg-prev">
                {{ chat.mostrecentmsg.content || 'msg loading...'}}
            </p>
            <p v-else-if="chat.mostrecentmsg.photo && !chat.mostrecentmsg.content">
                📷 Photo
            </p>
            <p v-else-if="chat.mostrecentmsg.content && chat.mostrecentmsg.photo" class="msg-prev">
                📷 {{ chat.mostrecentmsg.content || 'msg loading...'}}
            </p>
            <p v-else>error with message preview...</p>
        </div>
        <p v-else></p>
    </div>
</template>