<script>
export default{
    props:{
        mode:{type:String,required:true,validator:value=>["add","fwd"].includes(value),},
        chatid:{type:String,required:true},
        members:{type:Array, default:()=>[],},
        messageid:{type:String,default:""},
        isgroup:{type:Boolean, default:false},

    },
    data(){
        return{users:[],chats:[],searchtext:"",loading:false,errormsg:null,selected:[]};
    },
    computed:{
        availableUsers(){
            let users=this.users;
            if(this.mode==="add"){
                users=users.filter(user=>!this.members.includes(user.userid));
            }
            const search=this.searchtext.toLowerCase();
            return users.filter(user=>user.username?.toLowerCase().includes(search));

        },
        availableChats(){
            const search=this.searchtext.toLowerCase();
            return this.chats.filter(chat=>chat.chatname?.toLowerCase().includes(search));
        },
    },
    methods:{
        cancel(){
            this.$emit("close");
        },
        async loadItems(){
            this.loading=true;
            this.errormsg=null;
            try{
                if (this.mode==="add"){
                    const response = await this.$axios.get("/users");
                    this.users=response.data||[];
                }else{
                    const response = await this.$axios.get("/chats");
                    this.chats=response.data||[];
                }
            } catch(e){
                this.errormsg=e.response?.data?.error||e.message;
            } finally {this.loading=false;}
        },
        async addOrFwd(){
            if (this.selected.length===0) return;
            this.loading=true;
            this.errormsg=null;
            try{
                if (this.mode === 'fwd'){
                    await this.$axios.post(`chats/${this.chatid}/messages/${this.messageid}/forward`,
                    {targetchatIDs:this.selected,});
                }else{
                    await this.$axios.post(`chats/${this.chatid}/members`,
                    {targetuserIDs:this.selected,});
                }
                this.$emit("completed");
                this.$emit("close");
            } catch(e){
                this.errormsg=e.response?.data?.error || e.response?.data || e.message;
            } finally {this.loading= false;}

        },
        toggleselect(item, id){
            const alrselected=this.selected.includes(id);
            if (alrselected){
                this.selected=this.selected.filter(selectedid=>selectedid!==id);
            } else{this.selected.push(id)}
        }
    },
    mounted(){this.loadItems();},
};
</script>

<template>
    <div class="add-or-fwd">
        <h3>{{ mode==="add" ? "add members" : "forward message" }}</h3>
    <input v-model="searchtext" type="text" placeholder="search ..." />
    <p v-if="loading">Loading...</p>
    <p v-else-if="errormsg">{{ errormsg }}</p>    

    <template v-else>
        <template v-if="mode==='add'">
            <div v-for="user in availableUsers" :key="user.userid">
                <label>
                    <input type="checkbox"
                    :value="user.userid"
                    :checked="selected.includes(user.userid)"
                    @change="toggleselect(user,user.userid)" />
                    {{ user.username }} 
                </label>
            </div>
        </template>

        <template v-else>
            <div v-for="chat in availableChats" :key="chat.chatid">
                <label>
                    <input type="checkbox"
                    :value="chat.chatid"
                    :checked="selected.includes(chat.chatid)"
                    @change="toggleselect(chat,chat.chatid)" />
                    {{ chat.chatname || 'private chat' }}
                </label>
            </div>
        </template>

        <p v-if="mode==='add' && availableUsers.length===0">
            no available users to add!
        </p>
        <p v-if="mode==='fwd' && availableChats.length===0">
            no available chats to fwd to!
        </p>
        <button @click="cancel">✕</button>
        <!-- !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!-->
        <button @click="addOrFwd" :disabled="loading||selected.length===0||(mode==='add' && !isgroup)">
            {{ loading ? 'working...' : mode==="add" ? "add" : "forward" }}
        </button>
    </template>
    
    </div>
</template>