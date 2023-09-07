
<script>
import { getAuthToken, getAuthUsername } from '../services/tokenService';

export default {
    props: ["photoId"],
	data: function() {
		return {
			errormsg: null,
			loading: false,
            
	        username: null,
            imageURL: null,
            date: null,
            caption: null,

            nComments: 0,
            comments: [],
            commentText: null,
            nLikes: 0,
            likes: [],

            isLiked: false,
            isAuthor: false,
            
		}
	},
	methods: {
		load() {
			return load
		},
		async create() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/photo/" + this.photoId, { 
					headers: { 'Authorization': `Bearer ${getAuthToken()}`}
                });
				this.username = response.data.username;
            	this.date = response.data.date;
				this.caption = response.data.caption;
            	this.nComments = response.data.nComments;
            	this.nLikes = response.data.nLikes;
				this.isLiked = response.data.isLiked;
                this.isAuthor = this.username === getAuthUsername();

    			this.imageURL = `/storage/${this.photoId}.jpg`;
				console.log(this.imageURL);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

		async listComments() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/photo/" + this.photoId + "/comments");
				this.comments = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async listLikes() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/photo/" + this.photoId + "/likes");
				this.likes = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async likeUnlikeBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.isLiked) {
					await this.$axios.delete("/photo/" + this.photoId + "/likes/" + getAuthUsername());
					this.isLiked = false;
					this.nLikes = this.nLikes - 1 ;
				} else{ 
					await this.$axios.put("/photo/" + this.photoId + "/likes/" + getAuthUsername());
					this.isLiked = true;
					this.nLikes = this.nLikes + 1 ;
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async commentBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.post("/photo/" + this.photoId + "/comments", {
					text: this.commentText, headers: { 'Authorization': `Bearer ${getAuthToken()}`}
				});
				await listComments();
				this.nComments = this.nComments + 1 ;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.create()
	}
}
</script>
<template>
	<div>
	  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
		
	  </div>
  
	  <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  
	  <div class="mb-3">
		<div class="photo-container">
			<img v-if="imageURL" :src="imageURL" alt="Photo" class="img-fluid" />
		</div>
		<div class="photo-caption">
		  <p>{{ caption }}</p>
		</div>
		<div class="photo-info">
		  <p>Posted by {{ username }} on {{ date }}</p>
		  <p>{{ nLikes }} Likes</p>
		</div>
	  </div>
  
	  <div class="photo-actions">
		<button @click="likeUnlikeBtn" :disabled="loading" class="btn btn-primary">
		  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
		  {{ isLiked ? 'Unlike' : 'Like' }}
		</button>
  
		<button @click="listComments" :disabled="loading" class="btn btn-secondary">
		  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
		  Show comments
		</button>
	  </div>
  
	</div>
  </template>
  
  <style scoped>
  .photo-container {
	text-align: center;
  }
  
  .photo-caption {
	margin-top: 20px;
  }
  
  .photo-info {
	font-size: 14px;
	color: #777;
  }
  
  .photo-actions {
	display: flex;
	justify-content: space-between;
	margin-top: 20px;
  }
  
  .btn {
	display: flex;
	align-items: center;
  }
  
  .btn svg {
	margin-right: 5px;
  }
  </style>
  