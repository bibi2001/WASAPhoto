<!-- 
  Comment.vue is the common component of the list of comments.
  If the user is the author of the comment, they can delete the comment.
  By clicking on the name of the author, the user is redirected to the author's profile.

  Endpoints called:
  DELETE("/photo/:photoId/comments/:commentId")
-->


<script>
import { getAuthUsername, getAuthToken } from '../services/tokenService';

export default {
	props: ["photoId", "comment"],
	data: function () {
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
			this.$router.push("/profile/" + this.username);
		},
		async deleteBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/photo/" + this.commentPhotoId + "/comments/" + this.commentId, {
					headers: { 'Authorization': `Bearer ${getAuthToken()}` }
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
		<div class="my-1">
			<v-bind class="username" @click="goToProfile">{{ username }}</v-bind>
			<v-bind class="text">{{ text }}</v-bind>

			<a class="comment-delete" @click="deleteBtn" v-if="this.isAuthor">
				Delete this comment
			</a>

		</div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>
  
<style scoped>
.username {
	cursor: pointer;
	font-size: 12px;
	font-weight: bold;
	height: 0px;
	margin-left: 5px;
	margin-top:1rem;
}

.text {
	font-size: 17x;
	height: 0px;
	margin-left: 5px;
	margin-top:1rem;
}

.comment-delete {
	font-size: 11px;
	color: #777;
	cursor: pointer;
	margin-left: 5px;
	height: 0;
}
</style>
  