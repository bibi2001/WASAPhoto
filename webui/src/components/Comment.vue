<script>
import { getAuthUsername, getAuthToken } from '../services/tokenService';
export default {
    props: ["photoId","commentId"],
	data: function() {
		return {
			errormsg: null,
			loading: false,

            username: null,
    		text: null,
			isAuthor: false,
    
		}
	},
    methods: {
		async goToProfile() {
			this.$router.push("/profile/"+this.username);
		},
        async create() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/photo/" + this.photoId + "/comments/" + this.commentId, {
					headers: { 'Authorization': `Bearer ${getAuthToken()}`}
				});
				this.username = response.data.username;
                this.text = response.data.text;
				this.isAuthor = this.username === getAuthUsername();
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
        async delete() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/photo/" + this.photoId + "/comments/" + this.commentId, {
					headers: { 'Authorization': `Bearer ${getAuthToken()}`}
				});
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

    }, 
}
</script>

<template>
  <div class="card">
    <div class="card-header" @click="goToProfile" style="cursor: pointer;">
      <p class="username">{{ username }}</p>
    </div>
    <div class="card-body">
      <p class="text">{{ text }}</p>
      <span v-if="this.isAuthor" @click="deleteComment" class="delete-icon">
        <i data-feather="trash"></i>
      </span>
    </div>
  </div>
</template>

<style scoped>
.card {
	margin-bottom: 10px;
	height:60px;
}
.username {
	font-size: 8px;
}
.text {
	font-size: 15px;
}
</style>
