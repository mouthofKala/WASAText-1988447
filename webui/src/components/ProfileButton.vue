<script>
    export default{
        data: function(){
            return{
                photo: null
            }
        },
        computed:{
            photoUrl(){
                if (!this.photo) return null;
                const path = this.photo.replace(/^\/+/,"");
                return `${this.$axios.defaults.baseURL}/${path}`;
            },
        },
        methods: {
            async getMyProfile(){
                try{
                    const response = await this.$axios.get("/profile");
                    this.photo=response.data.photo;
                } catch(e){

                    console.error(e);
                }
            },
            goToProfile(){
                this.$router.push(`/profile`)
            }
        },
        mounted(){
            this.getMyProfile();
        },
        updated(){
            this.getMyProfile();
        }

    }
</script>

<template>
    <button class="profile-buttn" @click="goToProfile">
        <img v-if="photo" :src="photoUrl" alt="pfp" @error="photo = null">
        <span v-else>🧑</span> <!--substitute with default black pic-->
    </button>
</template>

<style scoped>
.profile-buttn{
    width: 40px;
    height:40px;
    border-radius:50%;
    overflow: hidden;
    padding:0;
    border: none;
    cursor: pointer;
}
</style>