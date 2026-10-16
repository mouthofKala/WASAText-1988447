<script>
export default {
    props: {
        chat: {
            type: Object,
            required: true
        }
    },
    methods: {
        getConversation(){
            this.$router.push(`/chats/${this.chat.chatid}/messages`)
        }
    }
}
</script>

<template>
    <div class="chat-item border rounded p-3 mb-2" @click="getConversation">
        <img v-if="chat.photo" alt="chat pfp" :src="`${$axios.defaults.baseURL}/${chat.photo.replace(/^\/+/, '')}`" class="item-photo" @error="console.error('error loaging: ', chat.photo)" />
        <h5>{{ chat.chatname }}</h5>
        <div v-if="chat.mostRecentMsg">
            <p class="mb-1">
                {{ chat.mostRecentMsg.username }}:
            </p>
            <p v-if="chat.mostRecentMsg.content && !chat.mostRecentMsg.photo" class="msg-prev">
                {{ chat.mostRecentMsg.content }}
            </p>
            <p v-else-if="chat.mostRecentMsg.photo && !chat.mostRecentMsg.content">
                📷 Photo
            </p>
            <p v-else-if="chat.mostRecentMsg.content && chat.mostRecentMsg.photo" class="msg-prev">
                📷 {{ chat.mostRecentMsg.content }}
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