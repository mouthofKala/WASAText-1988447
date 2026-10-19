<script>
    export default{
        data: function(){
            return{
                photo: null,
                refraiche: null,
            }
        },
        computed:{
            photoUrl(){
                if (!this.photo) return null;
                if (this.photo.startsWith("http://") ||
                    this.photo.startsWith("https://")) {
                    return this.photo;
                }
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
            this.refraiche=setInterval(()=>{
                this.getMyProfile();
            },3000);
        },
        beforeUnmount(){
            clearInterval(this.refraiche);
        }

    }
</script>

<template>
    <button class="profile-buttn" @click="goToProfile">
        <img v-if="photo" :src="photoUrl" alt="pfp" @error="photo = null">
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
.profile-buttn img{
    width: 100%;
    height: 100%;
    object-fit: cover;
}
}
</style>