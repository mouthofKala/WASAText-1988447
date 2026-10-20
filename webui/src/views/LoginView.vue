
<script>
import "../assets/login-page.css";
export default {
    data: function(){
        return {
            username: "",
            errormsg: null,
            loading: false,
            userID: "",
        }
    },
    methods: {
        async doLogin() {
            this.loading = true;
            this.errormsg = null;

        
            try {
                const response = await this.$axios.post("/", {
                    name: this.username
                });
                console.log("dologin:\n", response.data)
                const userID = response.data.identifier;
                sessionStorage.setItem("userID", userID);
                this.$axios.defaults.headers.common.Authorization = `Bearer ${userID}`;
                console.log("login success", response.data);
                this.$router.push(`/chats`);
            }catch(e){
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
    },
}
</script>

<template>
    <div class="login-page">
        <div class = "login-box">
            <h1> WASAText</h1>
            <p>Welcome! Please enter your (new or current) username to log in.</p>
            <form @submit.prevent="doLogin">
                <div class="mb-3">
                    <input
                        type="text"
                        class="form-control"
                        v-model="username"
                        placeholder="enter your username..."
                        :disabled="loading" />
                </div>
                <button
                    type="submit"
                    class="buttn"
                    :disabled="loading||!username.trim()"
                >
                {{ loading ? "logging in..." : "login"}}
                </button>
            </form>

            <ErrorMsg v-if="errormsg" :msg="errormsg" />
        </div>
        <div class = "version">
            version 1.0.0
        </div>
    </div>    
</template>