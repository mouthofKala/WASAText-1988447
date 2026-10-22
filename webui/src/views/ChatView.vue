
<script>
import "../assets/chat-page.css";
import "../assets/message-item.css";
import chatview from "../chatview.js";
export default{
    mixins:[chatview],
    data: function(){
        return{
            userID: sessionStorage.getItem("userID"),
            errormsg:null,
            loading:false,
            loadingOlder:false,
            convo: null,
            msglist:[],
            showbuttns: false,
            replyingto: null,
            msgtext:"",
            otherUser: null,
            imgFile: null,
            imgPrev:null,
            moremsg:true,
            fwdmsg: null,
            targetchatid: "",
            myChats:[],
            fwding: false,
            addOrFwdMode:null,
            userProfiles:{},
            refreshing: false,
            refresh: null,
            init: true,
        };
    },
    computed:{
        msgsWithReplies(){
            return this.msglist.map(message=>{
                const originalmsg=this.msglist.find(m=>m.messageid===message.replyto);
                return {...message,repliedMsg:originalmsg||null};
            });
        },
    },

    mounted(){
        this.getConversation();
        this.refresh=setInterval(()=>{this.getConversation();},3000);},
    beforeUnmount(){
        clearInterval(this.refresh);
        if (this.imgPrev){
            URL.revokeObjectURL(this.imgPrev);
        }
    }
};

</script>
<template>
    <div class="chat-page">
        <div class="chat-header">
            <button class="back" @click="this.$router.replace('/chats')">↩</button>
            <button v-if="convo" class="chat-title" @click="openChatProfile">
                <img 
                    v-if="convo.chat.grouporchat && convo.chat.photo"
                    :src="(convo.chat.photo ? `${$axios.defaults.baseURL}/${convo.chat.photo.replace(/^\/+/,'')}` : null)"
                    alt="group pic"
                    class="chat-profile-pic" 
                />

                <img v-else-if="otherUser && otherUser.photo"
                    :src="(otherUser.photo ? `${$axios.defaults.baseURL}/${otherUser.photo.replace(/^\/+/,'')}` : null)"
                    alt="user pic"
                    class="chat-profile-pic"/>
                <span v-else class="chat-profile-pic">😀</span>

                <span v-if="convo.chat.grouporchat">
                    {{ convo.chat.chatname }}
                </span>
                <span v-else-if="otherUser">
                    {{ otherUser.username }}
                </span>
                <span v-else>profile</span>
            </button>

            <button class="settings-button" @click="togglebuttns">⋮</button>

            <div v-if="showbuttns" class="chatmenuuu">
                <template v-if="convo?.chat.grouporchat">
                    <button class="addbuttn" @click="addMembers">Add users</button>
                    <button @click="leave" class="leavebttn">Leave</button>
                </template>
                <template v-else>
                    <button @click="leave" class="leavebttn">Leave</button>
                </template>
            </div>
        </div>

        <ErrorMsg v-if="errormsg" :msg="errormsg" />
        <LoadingSpinner v-if="loading" />
        <div v-else class="message-list" @scroll="handleScroll">
            <p v-if="msglist.length===0"></p>
            <MessageItem
                v-for="message in msgsWithReplies"
                :key="message.messageid"
                :message="message"
                :userID="userID"
                :senderprofile="userProfiles[message.userid]||null"
                @react="addReaction"
                @reply="replyToMsg"
                @unreact="delReaction"
                @delete="delMsg"
                @forward="startFwd"
            />
        </div>
        <AddOrFwd 
            v-if="addOrFwdMode"
            :mode="addOrFwdMode"
            :chatid="$route.params.chatid"
            :messageid="fwdmsg?.messageid ||''"
            :isgroup="convo?.chat?.grouporchat||false"
            :members="convo?.chat?.members||[]"
            @close="closeAddFwd"
            @completed="closeAddFwd"
        />

        <div v-if="replyingto" class="reply-prev">
            <!--in case msg is deleted/not found-->
            <strong> Replying to {{ replyingto.username }}</strong>
            <p>{{ replyingto.content }}</p> <!-- shorten it!-->
            <button type="button" @click="cancel" class="smol">✕</button>
        </div>

        <div v-if="imgPrev" class="img-prev">
            <img :src="imgPrev" alt="msgphoto" />
            <button type="button" @click="removeImg" class="smol">✕</button>
        </div>

        <div v-if="fwding" class="fwdable-list">
            <h4> foward to...</h4>
            <select v-model="targetchatid">
                <option value="" disabled>select a chat</option>
                <option v-for="chat in myChats" :key="chat.chatid" :value="chat.chatid">
                    {{ chat.grouporchat ? chat.chatname : "other user"}}
                </option>
            </select>
            <button @click="cancelFwd">cancel</button>
            <button @click="forward" :disabled="!targetchatid || fwding">
                {{ fwding ? "forwarding..." : "forward" }}
            </button>
        </div>
        <form class="msg-composer" @submit.prevent="sendMessage">
            <input ref="imageInput" type="file" accept="image/*" hidden @change="selectImg" />
            <button type="button" @click="chooseImg" aria-label="attach img" class="smol">📷</button>
            <input v-model="msgtext" type="text" placeholder="write a message..." maxlength="5000"/>
            <button type="submit" :disabled="(!msgtext.trim() && !imgFile) || msgtext.length>5000" aria-label="send msg" class="smol">➤</button>
        </form>
        
    </div>
</template>