<script>
export default {
    data(){
            return{
                showMenu: false
            }
    },
    props: {
        message: {
            type: Object,
            required: true
        }
    },
    methods: {
        openProfile(){
            this.$router.push(`/users/${this.message.userid}`)
        },
        actions(){
            this.showMenu = !this.showMenu
        },
        react(emoji) {
            this.$emit("react", {
                message: this.message,
                emoji: emoji
            })
            this.showMenu = false
        },
        reply(){
            this.$emit("reply", this.message)
            this.showMenu = false
        }
    }
}
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
        <div v-if="showMenu" class="menuuu">
            <button @click="react('❤️')">❤️</button>
            <button @click="react('👍')">👍</button>
            <button @click="react('😂')">😂</button>
            <button @click="react('😢')">😢</button>
            <button @click="reply">↩ Reply</button>
        </div>
    </div>
</template>

<!--see what to do to avoid weird img stretching-->
<style scoped>
    .message-item img {
        max-width: 100%;
        max-height: 400px;
        width: auto;
        height: auto;
        display: block;
    }


</style>