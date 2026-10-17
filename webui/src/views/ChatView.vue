
<script>
export default{
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
        };
    },
    methods:{
        async getConversation(){
            this.loading = true;
            this.errormsg = null;
            try{
                const response = await this.$axios.get( `/chats/${this.$route.params.chatid}/messages`)
                console.log("gtconversation:\n", response.data)
                this.convo=response.data;
                this.msglist = response.data.messages || [];
                const chat = this.convo.chat;
                if (!chat.grouporchat && chat.members){
                    const otherID = chat.members.find(id=>id!==this.userID);
                    if (otherID){
                        const response=await this.$axios.get(`/users/${otherID}`);
                        console.log("getuser:\n",response.data)
                        this.otherUser=response.data;
                    }
                }
            } catch (e) {
				this.errormsg = e.response?.data?.error || e.message || "could not load chat..";
				this.msglist=[];
			} finally {
				this.loading = false;
			}
        },
        async loadingOlderMsg(){
            if (this.loadingOlder || !this.moremsg || this.msglist.length === 0){
                return;
            }
            this.loadingOlder = true;
            try {
                const oldest = this.msglist[0];
                const response = await this.$axios.get(`/chats/${this.$route.params.chatid}/messages`,
                    {params: {before: oldest.timestamp}});
                console.log("getmsgs older:\n", response.data)
                const oldermsgs= response.data.messages||[];
                if (oldermsgs.length===0){this.moremsg=false;}
                else {
                    this.msglist=[
                        ...oldermsgs,
                        ...this.msglist
                    ];
                }
            } catch (e) {
                this.errormsg = e.toString();
            } finally {this.loadingOlder=false;}
        },

        handleScroll(event){
            if (event.target.scrollTop===0){
                this.loadingOlderMsg();
            }
        },
        async openChatProfile(){
            if (!this.convo) return;
            const chat = this.convo.chat;
            if (chat.grouporchat){
                this.$router.push(`/chats/${chat.chatid}/settings`);
            } else {
                const otherID = chat.members.find(id=>id!==this.userID);
                if (otherID){
                    this.$router.push(`/users/${otherID}`);
                }
            }
        },

        addMembers(){
            this.$router.push(`/chats/${this.$route.params.chatid}/members`);
        },

        async addReaction({message, emoji}){
            try{
                const response = await this.$axios.post(
                    `/chats/${this.$route.params.chatid}/messages/${message.messageid}/reactions`, {emoji}
                );
                console.log("addcomment:\n",response.data)
                if (!message.reactions){
                    message.reactions=[];
                }
                const idx=message.reactions.findIndex(r=>r.userid===this.userID);
                if (idx!==-1){message.reactions.splice(idx,1,response.data);}
                // message.reactions.push(response.data);
            } catch(e){
                this.errormsg=e.response?.data?.error||e.message;
            }
        },
        replyToMsg(message){
            this.replyingto=message;
        },
        cancel(){this.replyingto=null;},

        chooseImg(){
            this.$refs.imageInput.click();
        },

        selectImg(event){
            const file = event.target.files?.[0];
            if(!file) return;
            if(!file.type.startsWith('image/')) {
                this.errormsg='must select an img file';
                event.target.value='';
                return;
            }
            if (file.size > 10*1024*1024){
                this.errormsg='img max size must be 10MB'
                event.target.value='';
                return;
            }
            this.errormsg=null;
            this.imgFile=file;
            if (this.imgPrev){
                URL.revokeObjectURL(this.imgPrev);
            }
            this.imgPrev=URL.createObjectURL(file);
        },

        removeImg(){
            if (this.imgPrev){
                URL.revokeObjectURL(this.imgPrev);
            }
            this.imgFile=null;
            this.imgPrev=null;
            if (this.$refs.imageInput){
                this.$refs.imageInput.value='';
            }
        },

        async sendMessage(){
            const content = this.msgtext.trim();
            if (!content &&!this.imgFile) return;
            if (content.length>20000){
                this.errormsg = 'max length (20000 chars) exceeded';
                return;
            }
            const formData=new FormData();
            if (content){
                formData.append('content', content);
            }
            if (this.imgFile){
                formData.append('photo', this.imgFile);
            }
            if (this.replyingto){
                formData.append("replyto", this.replyingto.messageid);
            }
            try{
                const response = await this.$axios.post(
                    `/chats/${this.$route.params.chatid}/messages`, formData);
                console.log("sendmessage",response.data)
                const newmsg={
                    ...response.data,
                    username:response.data.username || 'me',
                    userid: response.data.userid||this.userID,
                    reactions: response.data.reactions||[],
                };
                this.msglist.push(newmsg);
                this.msgtext="";
                this.replyingto=null;
                this.removeImg();
            } catch(e){
                this.errormsg=e.toString();
                this.msgtext="";
                this.replyingto=null;  
            }
        },
        async startFwd(message){
            this.fwdmsg=message;
            this.targetchatid="";
            try{
				const response = await this.$axios.get("/chats");
				console.log("getting convs for forwarding: \n",response.data)
				this.myChats = Array.isArray(response.data) ? response.data : [];
            } catch(e){
                this.errormsg=e.response?.data?.error || e.message;
                this.fwdmsg=null;
            }
        },
        cancelFwd(){
            this.fwdmsg=null;
            this.targetchatid="";
        },
        async forward(){
            if (!this.fwdmsg || !this.targetchatid) return;
            this.fwding=true;
            this.errormsg=null;
            try{
                await this.$axios.post(`/chats/${this.$route.params.chatid}/messages/${this.fwdmsg.messageid}/forward`,
                    {targetchat_id: this.targetchatid});
                this.cancelFwd();
                
            } catch(e){
                this.errormsg=e.response?.data?.error||e.message;
            } finally{this.fwding=false;}
        },
        async delReaction({message, reac}){
            try{
                const response = await this.$axios.delete(
                    `/chats/${this.$route.params.chatid}/messages/${message.messageid}/reactions/${reac.reactionid}`);
                console.log("uncommentmsg:\n", response.data)
                message.reactions=(message.reactions||[]).filter(r=>r.reactionid!==reac.reactionid);
            }catch(e){
                this.errormsg=e.response?.data?.error ||e.message;
            }
        },

        async delMsg(msg){
            try {
                const response = await this.$axios.delete(
                    `/chats/${this.$route.params.chatid}/messages/${msg.messageid}`);
                console.log("deletemsg:\n", response.data)
                this.msglist=this.msglist.filter(m=>m.messageid!==msg.messageid);
            } catch(e){
                this.errormsg=e.toString();
            }
        },

		togglebuttns(){
			this.showbuttns = !this.showbuttns;
		},
		async leave(){
            try{
                const response = await this.$axios.delete(`/chats/${this.$route.params.chatid}/members/me`);
                console.log("leavechatgroup:\n", response.data)
                this.$router.back();
            }catch(e){
                this.errormsg=e.toString();
            }
        },
    },
    mounted(){
        this.getConversation()
    },
    beforeUnmount(){
        if (this.imgPrev){
            URL.revokeObjectURL(this.imgPrev);
        }
    }
};

</script>
<template>
    <div class="chat-page">
        <div class="chat-header">
            <button @click="this.$router.replace('/chats')">↩</button>
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

            <div v-if="showbuttns" class="menuuu">
                <template v-if="convo?.chat.grouporchat">
                    <button @click="addMembers">Add users</button>
                    <button @click="leave" class="leavebttn">Leave</button>
                </template>
                <template v-else>
                    <button @click="leave"class="leavebttn">Leave</button>
                </template>
            </div>
        </div>

        <ErrorMsg v-if="errormsg" :msg="errormsg" />
        <LoadingSpinner v-if="loading" />

        <div v-else class="message-list" @scroll="handleScroll">
            <p v-if="msglist.length===0"></p>
            <MessageItem
                v-for="message in msglist"
                :key="message.messageid"
                :message="message"
                :userID="userID"
                @react="addReaction"
                @reply="replyToMsg"
                @unreact="delReaction"
                @delete="delMsg"
                @forward="startFwd"
            />
        </div>

        <div v-if="replyingto" class="reply-prev">
            <!--in case msg is deleted/not found-->
            <strong> Replying to {{ replyingto.username }} || 'message'</strong>
            <p>{{ replyingto.content }}</p> <!-- shorten it!-->
            <button type="button" @click="cancel">✕</button>
        </div>

        <div v-if="imgPrev" class="img-prev">
            <img :src="imgPrev" alt="msgphoto" />
            <button type="button" @click="removeImg">✕</button>
        </div>

        <div v-if="fwding" class="fwdable-list">
            <h4> foward to...</h4>
            <select v-model="targetchatid">
                <option value="" disabled>select a chat</option>
                <option v-for="chat in myChats" :key="chat.chatid" :value="chat.chatid">
                    {{ chat.grouporchat ? chat.chatname : "SHOULD DISPLAY NAME OF OTHER USER"}}
                </option>
            </select>
            <button @click="cancelFwd">cancel</button>
            <button @click="forward" :disabled="!targetchatid || forwarding">
                {{ forwarding ? "forwarding..." : "forward" }}
            </button>
        </div>
        <form class="msg-composer" @submit.prevent="sendMessage">
            <input ref="imageInput" type="file" accept="image/*" hidden @change="selectImg" />
            <button type="button" @click="chooseImg" aria-label="attach img">📷</button>
            <input v-model="msgtext" type="text" placeholder="write a message..." maxlength="20000"/>
            <button type="submit" Ldisabled="!msgtext.trim() && !imgFile" aria-label="send msg">➤</button>
            <!--disable if msg has img and length exceeds ... chars TBD-->
        </form>
        
    </div>
</template>
<style scoped>
.chat-title{
    display: flex;
    align-items: center;
    gap: 10px;
    border: none;
    background: transparent;
    cursor: pointer;
}
.chat-page{
    height: 100vh;
    display: flex;
    flex-direction: column;
}
.chat-header{
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: relative;
    padding: 8px;
}
.chat-profile-pic
.message-list{
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    border-bottom-style: solid gainsboro
}
.chat-profile-pic{
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
}
.msg-composer{
    display: flex;
    gap: 8px;
    padding: 10px;
}
.msg-composer input{
    flex: 1;
}
.fwdable-list{
    padding:12px;
    border: 1px solid lightgray;
    display: flex;
    align-items: center;
    gap: 8px;
}

.settings-button{
    background-color: blueviolet;
    text-decoration-color:white;
}

.menuuu{
    position:absolute;
    top:100%;
    right:8px;
    z-index:10;
    display: flex;
    border: 1px solid grey;
    border-radius: 8px;
    padding: 8px;
    background: white;
    flex-direction: column;
}

.reply-prev{
    border-left: 3px solid gray;
    align-items: center;
    display: flex;
    justify-content: space-between;
    gap: 12px;
    padding: 8px 12px;
}
.reply-prev img{
    padding: 8px 12px;
}
.reply-prev p {
    margin: 4px 0 0;
}

.chat-title{
    background: transparent;
}


.img-prev{
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
}
.img-prev img{
    max-width: 120px;
    max-height: 120px;
    object-fit: contain
    border-radius=8px;
}

.msg-composer button:disabled{
    opacity: 0.5;
    cursor: not-allowed;
}

.fallback-pic{
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: #000;
}

.leavebttn{
    background: red;
    text-decoration-color: white;
}

</style>