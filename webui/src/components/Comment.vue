<script>
import { getAuthUsername, getAuthToken } from '../services/tokenService';
export default {
    props: ["photoId","comment"],
	data: function() {
		return {
			errormsg: null,
			loading: false,

			commentPhotoId: this.photoId,
			commentId: this.comment.commentId,
            username: this.comment.username,
    		text: this.comment.text,

			isAuthor: this.comment.username === getAuthUsername(),
			isDeleted: false,
		}
	},
    methods: {
		goToProfile() {
			this.$router.push("/profile/"+this.username);
		},
		async deleteBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/photo/" + this.commentPhotoId +"/comments/" + this.commentId, {
					headers: { 'Authorization': `Bearer ${getAuthToken()}`}
				});
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.isDeleted = true;
			this.loading = false;
		},
        
},

}
</script>

<template>
	<div class="card-item" v-if="!this.isDeleted">
		<a class="comment-delete" @click="deleteBtn" v-if="this.isAuthor">
		  Delete this comment
		</a>
		<div class="card-title">
			<p class="username" @click="goToProfile">{{ username }}</p>
		</div>
		<div class="card-text">
			<p class="text">{{ text }}</p>
	  	</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>
  
<style scoped>
.card-item {
	height: 50px;
	-ms-wrap-margin: 10px;
}
  
.username {
	cursor: pointer;
	font-size: 13px;
	font-weight: bold;
	height: 0px;
}
  
.text {
	font-size: 15px;
}

.comment-delete {
	font-size: 12px;
	color: #777;
	cursor: pointer;
  }
 
</style>
  