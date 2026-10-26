export default{
    methods:{
        async getProfile(){
            if (this.refreshing) return;
            this.refreshing=true;
            if (this.init){this.loading=true;}
            this.errormsg = null;

            try{
                if (this.$route.path ==="/profile"){
                    this.profileType="mine";
                    const response = await this.$axios.get("/profile");
                    console.log("getMYuserprofile:\n", response.data)
                    this.profile = response.data;
                    this.newName = this.profile.username;


                } else if(this.$route.params.userid){
                    this.profileType="user";
                    const response = await this.$axios.get(`/users/${this.$route.params.userid}`);
                    console.log("getuserprofile:\n", response.data)
                    this.profile = response.data;


                } else if (this.$route.params.chatid){
                    this.profileType = "group";
                    const response = await this.$axios.get(`/chats/${this.$route.params.chatid}/messages`);
                    console.log("getgroupprofile:\n", response.data)
                    this.profile = response.data.chat;
                    this.newName = this.profile.chatname;


                }
            }catch(e){
                this.errormsg=e.toString();
                // this.profile="";
                // this.newName="";

            } finally {this.loading=false;this.init=false;this.refreshing=false;}
        },
        async newNameSave(){
            try{
                if (this.profileType==="mine"){
                    const newUsername=this.newName.trim();
                    const response = await this.$axios.put("profile/username",{
                        targetusername: newUsername,
                    });
                    console.log("setmyusername:\n", response.data)
                    this.profile.username = newUsername;
                } else if (this.profileType === "group"){
                    const newUsername=this.newName.trim();
                    const response = await this.$axios.put(`/chats/${this.$route.params.chatid}/groupname`,
                            {newname: newUsername,});
                    console.log("setgroupusername:\n", response.data);
                    this.profile.chatname = newUsername;

                };
                
            } catch(e){
                this.errormsg=e.toString();
                this.newName=this.profile.username;
            }
        },

        async newPhotoSave(){
            if (!this.newPhoto||this.savingPhoto) return;
            this.savingPhoto=true;
            this.errormsg=null;
            try{
                let url;
                if (this.profileType==="mine"){
                    url='/profile/photo';
                } else if (this.profileType==="group"){
                    url=`/chats/${this.$route.params.chatid}/photo`;
                } else {
                    return;
                }
                const response = await this.$axios.put(url,this.newPhoto,
                    {
                        headers: {'Content-Type': this.newPhoto.type,},
                    }
                );
                console.log("setnewphoto:\n", response.data)
                const photoURI = response.data;
                if (this.profileType==='mine'){
                    this.profile.photo=photoURI;
                } else{this.profile.photo=photoURI;}
                this.newPhoto=null;
                if (this.imgPrev){
                    URL.revokeObjectURL(this.imgPrev);
                    this.imgPrev=null;
                }
                if (this.$refs.photoInput){
                    this.$refs.photoInput.value='';
                }
            } catch(e){
                this.errormsg=e.response?.data?.error || e.message ||e.toString();
                this.newPhoto=null;
            } finally {
                this.savingPhoto=false;
            }
        },

        selectPhoto(event){
            const file = event.target.files?.[0];
            if (!file) return;
            if (!file.type.startsWith('image/')){
                this.errormsg='selected file must be an img';
                event.target.value='';
                return;
            }
            if (file.size > 10*1024*1024){
                this.errormsg='img size limit is 10 MB';
                event.target.value='';
                return;
            }
            this.errormsg=null;
            this.newPhoto=file;

            if (this.imgPrev){
                URL.revokeObjectURL(this.imgPrev);
            }
            this.imgPrev=URL.createObjectURL(file);
        },

    },
}