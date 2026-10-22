export default{
    methods:{
        async getConversation(){
            if (this.refreshing) return;
            this.refreshing = true;
            if (this.init){this.loading=true;}
            this.errormsg = null;
            try{
                const response = await this.$axios.get( `/chats/${this.$route.params.chatid}/messages`)
                console.log("gtconversation:\n", response.data)
                this.convo=response.data;
                this.msglist = response.data.messages || [];

                const senderIDs=[...new Set(this.msglist.map(message=>message.userid).filter(Boolean))];
                await Promise.all(senderIDs.map(id=>this.getUserProfile(id)));
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
				this.errormsg = e.response?.data?.error || e.message || "could not load chat :(";
			} finally {
				this.loading = false;
                this.init=false;
                this.refreshing=false;
			}
        },
        async getUserProfile(userID){
            if (!userID) return;
            if (this.userProfiles[userID]) return;
            try{
                const response = await this.$axios.get(`/users/${userID}`);
                this.userProfiles[userID]=response.data;
            } catch(e){
                console.error(`couldn't get profile for ${userID}: `, e);
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
                    const senderIDs = [...new Set(oldermsgs.map(message=>message.userid).filter(Boolean))];

                    await Promise.all(senderIDs.map(id=>this.getUserProfile(id)));

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
            console.log("chat=", chat);
            console.log("members=", chat.members);
            console.log("my userID=", this.userID);
            if (chat.grouporchat){
                this.$router.push(`/chats/${chat.chatid}/settings`);
            } else {
                const otherID = chat.members.find(id=>id!==this.userID);
                if (otherID){
                    this.$router.push(`/users/${otherID}`);
                }
            }
        },

        async addMembers(){
            this.addOrFwdMode="add";
            const response= await this.$axios.post(`/chats/${this.$route.params.chatid}/members`);
            // !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
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
                else{message.reactions.push(response.data);}
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
            if (content.length>5000){
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
                this.getUserProfile(newmsg.userid);
                this.msgtext="";
                this.replyingto=null;
                this.removeImg();
            } catch(e){
                this.errormsg=e.toString();
                this.msgtext="";
                this.replyingto=null;  
            } finally{
                this.imgPrev=null;
            }
        },
        async startFwd(message){
            this.addOrFwdMode="fwd";
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
        closeAddFwd(){this.addOrFwdMode=null;},
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
        async delReaction({message, reaction}){
            try{
                const response = await this.$axios.delete(
                    `/chats/${this.$route.params.chatid}/messages/${message.messageid}/reactions/${reaction.reactionid}`);
                console.log("uncommentmsg:\n", response.data)
                message.reactions=(message.reactions||[]).filter(r=>r.reactionid!==reaction.reactionid);
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


}