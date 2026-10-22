export default{
    methods:{
        async getUsers(){
            try{
                const response = await this.$axios.get('/users');
                console.log("getusers:\n",response.data)
                this.users=response.data;
            } catch(e){
                this.errorMsg = "couldn't load users...";
            this.users=[]}
        },
        async selectPhoto(event){
            const file = event.target.files?.[0];
            if (!file) return;
            if (!file.type.startsWith('image/')){
                this.errorMsg='file must be an image!';
                return;
            }
            this.groupPhoto=file;
            this.photoPreview=URL.createObjectURL(file);
            this.errorMsg=null;
        },
        getUserID(user){
            return user.userid; 
        },
        getPhoto(user){
            return user.photo;
        },
        getUsername(user){
            return user.username;
        },
        toggleUser(user){
            const id = this.getUserID(user);
            const idx = this.selectedUsers.findIndex(selected => this.getUserID(selected) === id);
            if (idx!== -1){
                this.selectedUsers.splice(idx, 1);
            } else {
                this.selectedUsers.push(user)
            }
        },
        async makeChatGroup(){
            if (!this.groupName.trim()){
                this.errorMsg='please enter a name for the group';
                return;
            }
            if (this.selectedUsers.length===0){
                this.errorMsg='please select at least one user';
                return;
            }

            this.loading=true;
            this.errorMsg=null;
            try{
                let photo = '';
                if (this.groupPhoto){
                    photo = await new Promise((resolve, reject) =>{
                        const reader = new FileReader();
                        reader.onload = () => {resolve(reader.result.split(',')[1]);};
                        reader.onerror=reject;
                        reader.readAsDataURL(this.groupPhoto);
                    });
                }
                const response=await this.$axios.post('/chats',{
                    chatname: this.groupName.trim(),
                    photo,
                    members: this.selectedUsers.map(user=>this.getUserID(user)),
                    grouporchat: true,
                });
                console.log("makegroup:\n", response.data)
                const chatID=response.data.chatid ?? response.data.chat?.chatid;
                if (chatID){
                    await this.$router.replace(`/chats/${chatID}`);
                } else{
                    await this.$router.replace(`/chats/`);
                }
            }catch(e){
                console.error("failed to make group: ", e);
                this.errorMsg=e.response?.data?.error || e.response?.data || "couldn't create group :(";
            } finally{
                this.loading=false;
            }
        },


    },
}