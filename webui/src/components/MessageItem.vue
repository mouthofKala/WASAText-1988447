<script>
export default {
    data(){
            return{
                showMenu: false,
                // senderProfile:null, // remove?
                profileerr:null,
            }
    },
    props: {
        message: {
            type: Object, required: true
        },
        userID: {type: String, required: true},
        senderProfile:{type:Object, default:null},
    },
    computed:{
        myReac()   {
            return (this.message.reactions||[]).find(reaction=>reaction.userid===this.userID);
        },
        isMyMsg(){
            return this.message.senderid===this.userID;
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
            <img v-if="senderProfile?.photo"
            :src="`${$axios.defaults.baseURL}/${senderProfile.photo.replace(/^\/+/,'')}`"
            alt="sender pfp"
            class="sender-pfp" />
            {{ message.username }}:</button>
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
                <button @click="react('❤️')">❤️</button>
                <button @click="react('👍')">👍</button>
                <button @click="react('😂')">😂</button>
                <button @click="react('😭')">😭</button>
            </template>
            <button v-else @click="unreact">
                remove comment
            </button>
            <button @click="reply">↩ Reply</button>
            <button @click="forward">forward</button>
            <button v-if="isMyMsg" @click="delmsg">delete</button>
        </div>
    </div>
</template>

<style scoped>
    .message-item img {
        max-width: 100%;
        max-height: 400px;
        width: auto;
        height: auto;
        display: block;
    }
    .sender-pic{
        width:30px;
        height:30px;
        border-radius:50%;
        object-fit:cover;
        display:inline-block;
        vertical-align: middle;
    }
    .msg-user .sender-pic{
        display: inline-block;
        width: 30px;
        height: 30px;
        object-fit: cover;
    }
    .reply{
        border-left: 3px solid gray;
        margin-bottom: 8px;
        align-items: center;
        display: flex;
        justify-content: space-between;
        gap: 12px;
        padding: 8px 12px;
    }
    .reply img{padding:8px 12px;}
    .reply p{
        margin: 4px 0 0;
        opacity: 0.8;
        overflow: hidden;
        text-overflow: ellipsis;
    }


</style>