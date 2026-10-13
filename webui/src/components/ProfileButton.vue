<script>
    export default{
        data: function(){
            return{
                photo: null
            }
        },
        methods: {
            async getMyProfile(){
                try{
                    const response = await this.$axios.get("/users/me");
                    this.photo=response.data.photo;
                } catch(e){

                    console.error(e);
                }
            },
            goToProfile(){
                this.$router.push(`/users/me`)
            }
        },
        mounted(){
            this.getMyProfile();
        }

    }
</script>

<template>
    <button class="profile-buttn" @click="goToProfile">
        <img v-if="photo" :src="photo" alt="🧑">
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