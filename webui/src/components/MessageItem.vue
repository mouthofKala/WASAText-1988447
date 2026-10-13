<script>
export default {
    data(){
            return{
                showMenu: false,
            }
    },
    props: {
        message: {
            type: Object, required: true
        },
        userID: {type: String, required: true},
    },
    computed:{
        myReac()   {
            return (this.message.Reactions||[]).find(reaction=>reaction.UserID===this.userID);
        },
        isMyMsg(){
            return this.message.SenderID===this.UserID;
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
        <button class="msg-user" @click="openProfile">{{ message.username }}:</button>
        <button class="all-msg" @click="actions">
            <img v-if="message.photouri" :src="message.photouri" alt="msgphoto" />
            <p v-if="message.content" class="msg-content">
                {{ message.content }} 
            </p>
        </button>
        <!--already present reactions on msg-->
        <div v-if="message.Reactions?.length" class="reacs">
            <span v-for="reac in message.Reactions" :key="reac.ReactionID">{{ reaction.Emoji }}</span>
        </div>
        <div v-if="showMenu" class="menuuu">
            <template v-if="!myReac">
                <button @click="react('❤️')">❤️</button>
                <button @click="react('👍')">👍</button>
                <button @click="react('😂')">😂</button>
                <button @click="react('😢')">😢</button>
            </template>
            <button v-else @click="unreact">
                remove comment
            </button>
            <button @click="reply">↩ Reply</button>
            <button v-if="isMyMsg" @click="delmsg">
                delete
            </button>
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


</style>