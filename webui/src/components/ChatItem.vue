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
            this.$router.push(`/chats/${this.chat.chatid}`)
        }
    }
}
</script>

<template>
    <div class="chat-item border rounded p-3 mb-2" @click="getConversation">
        <h5>{{ chat.chatName }}</h5>
        <div v-if="chat.mostRecentMsg">
            <p v-if="chat.mostRecentMsg.content && !chat.mostRecentMsg.photo">
                {{ chat.mostRecentMsg.content }}
            </p>
            <p v-else-if="chat.mostRecentMsg.photo && !chat.mostRecentMsg.content">
                📷 Photo
            </p>
            <p v-else-if="chat.mostRecentMsg.content && chat.mostRecentMsg.photo">
                📷 {{ chat.mostRecentMsg.content }}
                <!--to be shortened to fit device, max 2 rows of text-->
            </p>
            <p v-else>error with message preview...</p>
        </div>
        <p v-else></p>
    </div>
</template>

<style scoped>
    .chat-item{
        cursor:pointer;
    }
</style>