<script>
import "../assets/profile-page.css";
import profileview from "../profileview.js";
export default{
    mixins:[profileview],
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
            refresh: null,
            refreshing: false,
            init: true,
        };
    },

    mounted(){
        this.getProfile();
        this.refresh=setInterval(()=>{if (!this.savingPhoto){this.getProfile();}},5000);
    },
    beforeUnmount(){
        clearInterval(this.refresh);
        if (this.imgPrev){
            URL.revokeObjectURL(this.imgPrev);
        }
    }

}
</script>


<template>
    <div class="profile-page">
        <button class="back" @click="this.$router.back()">↩</button>

        <ErrorMsg v-if="errormsg" :msg="errormsg" />
        <LoadingSpinner v-if="loading" />

        <div v-else-if="profile" class="profilediv">
            <img
                v-if="imgPrev || profile.photo"
                :src="imgPrev || (profile.photo ? `${$axios.defaults.baseURL}/${profile.photo.replace(/^\/+/,'')}` : null)"
                alt="pfp"
                class="profile-pic" />
            <input ref="photoInput" type="file" accept="image/*" hidden @change="selectPhoto"/>
            <template v-if="profileType==='mine'||profileType==='group'">
                <button class="botton" type="button" @click="$refs.photoInput.click()">choose photo</button>
                <button class="botton" type="button" @click="newPhotoSave" :disabled="!newPhoto ||savingPhoto">
                    {{ savingPhoto ? 'saving...' : 'save photo' }}
                </button>
            </template>

            <h2 v-if="profileType==='mine' || profileType==='user'">
                {{ profile.username }}
            </h2>
            <h2 v-else-if="profileType==='group'">
                {{ profile.chatname }}
            </h2>
            <template v-if="profileType==='mine'||profileType==='group'">
                <input v-model="newName" />
                <button class="botton" @click="newNameSave">update username</button>
            </template>

        </div>

    </div>
</template>
