<script>
import { getAuthUsername } from '../services/tokenService';
export default {
    props: ["photoId","comment"],
	data: function() {
		return {
			errormsg: null,
			loading: false,

			commentId: this.comment.commentId,
            username: this.comment.username,
    		text: this.comment.text,
			isAuthor: this.comment.username === getAuthUsername(),
    
		}
	},
    methods: {
		async goToProfile() {
			this.$router.push("/profile/"+this.username);
		},
		async deleteBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/photo/" + this.photoId +"/comments/" + this.commentId, {
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
	  <div class="card-body">
		<p class="username" @click="goToProfile">{{ username }}</p>
		<p class="text">{{ text }}</p>
		<span v-if="isAuthor" @click="deleteBtn" class="delete-icon">
		  <i data-feather="trash"></i>
		</span>
	  </div>
	</div>
  </template>
  
  <style scoped>
  .card {
	margin-bottom: 10px;
	padding: 10px; /* Add padding to the card body for spacing */
	border: 1px solid #ddd; /* Add a border for separation */
  }
  
  .username {
	cursor: pointer;
	font-size: 10px;
	margin-bottom: 3px;
	color: #007bff; /* Change the username color to a more prominent color */
	font-weight: bold; /* Make the username bold */
  }
  
  .text {
	font-size: 12px;
	color: #333; /* Darken the text color for better readability */
	line-height: 1.4; /* Add line-height for better text spacing */
  }
  
  /* Additional styling for the delete icon */
  .delete-icon {
	cursor: pointer;
	float: right;
	font-size: 14px;
	color: #777;
  }
  </style>
  