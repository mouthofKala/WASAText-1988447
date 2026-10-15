<script>
export default{
    data() {
        return{
            errormsg:null,
            loading:false,
            user: null,
            profile: null,
            profileType: null,
            newName: "",
            newPhoto: null,
            imgPrev: null,
            savingPhoto: false,
        };
    },
    methods:{
        async getProfile(){
            this.loading = true;
            this.errormsg = null;

            try{
                if (this.$route.path ==="/profile"){
                    this.profileType="mine";
                    const response = await this.$axios.get("/profile");
                    console.log("getMYuserprofile:\n", response.data)
                    this.profile = response.data;
                    this.newName = this.profile.username;


                } else if(this.$route.params.userID){
                    this.profileType="user";
                    const response = await this.$axios.get(`
                    /users/${this.$route.params.userID}`);
                    console.log("getuserprofile:\n", response.data)
                    this.profile = response.data.chat;


                } else if (this.$route.params.chatid){
                    this.profileType = "group";
                    const response = await this.$axios.get(`
                    /chats/${this.$route.params.chatid}/messages`);
                    console.log("getgroupprofile:\n", response.data)
                    this.profile = response.data.chat;
                    this.newName = this.profile.chatname;


                }
            }catch(e){
                this.errormsg=e.toString();
                this.profile="";
                this.newName="";
            } finally {this.loading=false;}
        },
        async newNameSave(){
            try{
                if (this.profileType==="mine"){
                    const newUsername=this.newName.trim();
                    const response = await this.$axios.put("profile/username",{
                        username: newUsername,
                    });
                    console.log("setmyusername:\n", response.data)
                    this.profile.username = newUsername;
                } else if (this.profileType === "group"){
                    const newUsername=this.newName.trim();
                    const response = await this.$axios.put(`
                        /chats/${this.$route.params.chatid}/groupname`,
                            {groupname: newUsername,});
                    console.log("setgroupusername:\n", response.data);
                    this.profile.username = newUsername;

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
                this.errormsg=e.toString();
                this.newPhoto=this.profile.photo;
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
    mounted(){
        this.getProfile();
    },
    beforeUnmount(){
        if (this.imgPrev){
            URL.revokeObjectURL(this.imgPrev);
        }
    }

}
</script>


<template>
    <template v-if="errormsg">
        <ErrorMsg :msg="errormsg"></ErrorMsg>
    </template>
    <template v-else>
        <div class="profile-page">
            <button @click="this.$router.back()">↩</button>

            <ErrorMsg v-if="errormsg" :msg="errormsg" />
            <LoadingSpinner v-if="loading" />

            <div v-else-if="profile">
                <img
                    v-if="imgPrev || profile.photo"
                    :src="imgPrev || (profile.photo ? `${$axios.defaults.baseURL}/${profile.photo}` : null)"
                    alt="pfp"
                    class="profile-pic" />
                <input ref="photoInput" type="file" accept="image/*" hidden @change="selectPhoto"/>
                <template v-if="profileType==='mine'||profileType==='group'">
                    <button type="button" @click="$refs.photoInput.click()">choose photo</button>
                    <button type="button" @click="newPhotoSave" :disabled="!newPhoto ||savingPhoto">
                        {{ savingPhoto ? 'saving...' : 'save photo' }}
                    </button>
                </template>

                <h2 v-if="profileType==='mine' || profileType==='user'">
                    {{ profile.Username }}
                </h2>
                <h2 v-else-if="profileType==='group'">
                    {{ profile.ChatName }}
                </h2>
                <template v-if="profileType==='mine'||profileType==='group'">
                    <input v-model="newName" />
                    <button @click="newNameSave">update username</button>
                </template>

            </div>

        </div>
    </template>
</template>

<style scoped>
.profile-page{
    padding: 20px;
}
.profile-pic{
    width: 50vw;
    height: auto;
    border-radius:50%;
    object-fit: cover;
    display: block;
}
</style>
