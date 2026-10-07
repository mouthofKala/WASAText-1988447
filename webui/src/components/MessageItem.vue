<script>
export default {
    data(){
        showMenu: false
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
        }
    }
}
</script>

<template>
    <div class="message-item">
        <button class="msg-user" @click="openProfile">{{ message.username }}:</button>
        <div class="all-msg" @click="showMenu">
            <img v-if="message.photouri" :src="message.photouri" alt="msgphoto" />
            <p v-if="message.content" class="msg-content">
                {{ message.content }} 
            </p>
        </div>
        <div v-if="showMenu" class="menuuu">
            <button>❤️</button>
            <button>👍</button>
            <button>😂</button>
            <button>😢</button>
            <button>↩ Reply</button>
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