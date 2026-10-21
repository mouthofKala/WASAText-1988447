<script>
import "../assets/message-item.css";
export default {
    data(){
            return{
                showMenu: false,
                // senderProfile:null,
                profileerr:null,
            }
    },
    props: {
        message: {type: Object, required: true},
        userID: {type: String, required: true},
        senderprofile:{type:Object,default:null,},
    },
    computed:{
        myReac()   {
            return (this.message.reactions||[]).find(reaction=>reaction.userid===this.userID);
        },
        isMyMsg(){
            return this.message.userid===this.userID; // idk y it worked anyway b4 lol
        }
    },
    methods: {
        openProfile(){
            this.$router.push(`/users/${this.message.userid}`);
        },
        actions(){
            this.showMenu = !this.showMenu
        },
        react(emoji) {
            this.$emit("react", {
                message: this.message,
                emoji
            });
            this.showMenu = false
        },
        reply(){
            this.$emit("reply", this.message)
            this.showMenu = false
        },
        forward(){
            this.$emit("forward", this.message)
            this.showMenu=false
        },
        unreact(){
            if (!this.myReac) return;
            this.$emit("unreact", {
                message: this.message,
                reaction:this.myReac,
            });
            this.showMenu = false
        },
        delmsg(){
            this.$emit("delete", this.message);
            this.showMenu=false;
        },
    },
};
</script>

<template>
    <div class="message-item">
        <button class="msg-user" @click="openProfile">
            <img v-if="senderprofile?.photo"
            :src="`${$axios.defaults.baseURL}/${senderprofile.photo.replace(/^\/+/,'')}`"
            alt="sender pfp"
            class="sender-pfp" />
            <p>{{ message.username }}:</p>
            </button>
        <button class="all-msg" @click="actions">
            <div v-if="message.replyto" class="reply">
                <template v-if="message.repliedMsg">
                    replying to {{ message.repliedMsg.username }}:
                    <p class="reply-content">{{ message.repliedMsg.content || '📷 photo'}}</p>
                </template>
                <template v-else>replying to a msg </template> <!--when msg got deleted-->
            </div>
            <p v-if="message.fwdfrom" class="fwd">forwarded</p>
            <img v-if="message.photouri"
                :src="(message.photouri ? `${$axios.defaults.baseURL}/${message.photouri.replace(/^\/+/,'')}` : null)"
                alt="msgphoto" />
            <p v-if="message.content" class="msg-content">
                {{ message.content }} 
            </p>
            <p>{{ message.status === "sent" ? "🗸" : "🗸🗸" }}</p>
        </button>

        <div v-if="message.reactions?.length" class="reacs">
            <span v-for="reac in message.reactions" :key="reac.reactionid">{{ reac.emoji }}</span>
        </div>
        <div v-if="showMenu" class="menuuu">
            <template v-if="!myReac">
                <button @click="react('❤️')" class="buttn">❤️</button>
                <button @click="react('👍')" class="buttn">👍</button>
                <button @click="react('😂')" class="buttn">😂</button>
                <button @click="react('😭')" class="buttn">😭</button>
            </template>
            <button v-else @click="unreact">
                remove comment
            </button>
            <button @click="reply" class="buttn">↩ Reply</button>
            <button @click="forward" class="buttn">forward</button>
            <button v-if="isMyMsg" @click="delmsg" class="buttn">delete</button>
        </div>
    </div>
</template>