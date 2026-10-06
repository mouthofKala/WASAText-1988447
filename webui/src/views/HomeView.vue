<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			chatlist: [],
		}
	},
	methods: {
		async getMyConversations() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/chats");
				this.chatlist = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.getMyConversations()
	}
}
</script>

<template>
	<div class="home-page">
		<div class="home-header">
			<h1 class="h2">Conversations</h1>
			<div class="search-bar">
				<input type="text" class="form-control" placeholder="search users...">
			</div>
			<button class="menu-button">⋮</button>
			<!--add click functions for ^v-->
			<button class="profile-button">🧑</button>
		</div>
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<LoadingSpinner v-if="loading" />
		<div v-else-if="chatlist.length === 0">
			<p>you have no chats :[ </p>
		</div>
		<div v-else class="chat-list">
			<ChatItem
				v-for="chat in chatlist" :key="chat.chatid" :chat="chat" />
		</div>
	</div>
</template>

<style scoped>
	.home-page{
		height: 100vh;
		display: flex;
		flex-direction: column;
	}

	.home-header{
		flex-shrink: 0;

	}

	.chat-list{
		flex: 1;
		overflow-y: auto;
	}
</style>
