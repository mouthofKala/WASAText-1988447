
<script>
export default {
    data: function(){ // state
        return {
            username: "",
            errormsg: null,
            loading: false, // waiting for backend
        }
    },
    methods: {
        async doLogin() {
            this.loading = true;
            this.errormsg = null;
        
            try {
                let response = await this.$axios.post("/session", {
                    username: this.username
                });
                console.log("login success", response.data);
                // ?
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
                    <label for="username" class="form-label"> username</label>
                    <input
                        id="username"
                        type="text"
                        class="form-control"
                        v-model="username"
                        placeholder="enter your username..."
                        :disabled="loading"
                </div>
                <button
                    type="submit"
                    class="btn btn-primary w-100"
                    :disabled="loading||!username"
                >
                {{ loading ? "logging in..." : "login"}}
                </button>
            </form>

            <ErrorMsg
                v-if="errormsg"
                :msg="errormsg"
            ></ErrorMsg>
        </div>
        <div class = "version">
            version 1.0.0
        </div>
    </div>    
</template>

<style scoped>
.login-page {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    position: relative;
}

.login-box {
    width: 100%;
    max-width: 400px;
    padding: 2rem;
}

.version {
    position: absolute;
    bottom: 1rem;
    font-size: 0.8rem;
    color: --bs-gray-500
}


</style>