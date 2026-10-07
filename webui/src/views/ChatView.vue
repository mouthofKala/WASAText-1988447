<script>
export default{
    data: function(){
        return{
            errormsg:null,
            loading:false,
            loadingOlder:false,
            convo: null,
            msglist:[],
        }
    },
    methods:{
        async getConversation(){
            this.loading = true;
            this.errormsg = null;
            try{
                let response = await this.$axios.get( `/chats/${this.$route.params.chatid}`)
                this.convo=response.data;
                this.msglist = response.data.Messages;
            } catch(e){
                this.errormsg=e.toString();
            }
            this.loading=false;
        },
        async loadingOlderMsg(){
            if (this.loadingOlder || this.msglist.length ===0){
                return;
            }
            this.loadingOlder = true;
            try {
                const oldest = this.msglist[0];
                let response = await this.$axios.get(`/chats/${this.$route.params.chatid}`,
                    {
                        params: {
                            before: oldest.timestamp
                        }
                    }
                );
                const oldermsgs= response.data.Messages;
                this.msglist=[
                    ...oldermsgs,
                    ...this.msglist
                ];
            } catch (e){
                this.errormsg = e.toString();
            }
            this.loadingOlder=false;
        },

        handleScroll(event){
            if (event.target.scrollTop===0){
                this.loadingOlderMsg();
            }
        }
    },
    mounted(){
        this.getConversation()
    }
}
</script>

<template>
    <div class="chat-page">
        <div class="chat-header">
            <button class="chat-profile-pic">😀</button>
            <!--pfp which is response data if group,
            otherwise it is the other user's pfp
            if you click it and it is a group, you can
            modify it (opens a popup on screen. if you
            click outside of it, it disappears)-->

            <h1 v-if="convo">{{ convo.Chat.chatname }}</h1>
            <!--if you click it and it is a groupchat name,
            you can change it. same as above-->

            <button class="settings-button">⋮</button>
            <!--three dot menu, if click ou get a small menu
            with addtogroup button and leave button.-->
        </div>

        <ErrorMsg v-if="errormsg" :msg="errormsg" />
        <LoadingSpinner v-if="loading" />

        <div v-else class="message-list" @scroll="handleScroll">
            <p v-if="msglist.length===0"></p>
            <MessageItem v-for="message in msglist" :key="message.messageid" :message="message" />
            <!--SORT THEM!
            clicking open options: forward, comment and,
            if already commented, uncomment. also delete
            message if it is your own msg-->
        </div>
        
    </div>
</template>

<style scoped>
    .chat-page{
        height: 100vh;
        display: flex;
        flex-direction: column;
    }
    .chat-header{
        flex-shrink: 0;
    }
    .message-list{
        flex: 1;
        overflow-y: auto;
    }


</style>