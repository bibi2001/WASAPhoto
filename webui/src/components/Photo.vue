
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
			showComments: false,

			isDeleted: false,
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
			this.showComments = true
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
				await this.$axios.post("/photo/" + this.photoId +"/comments", {
					text: this.textComment,}, {
					headers: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${getAuthToken()}`,
					},
					});
				await listComments();
				this.nComments = this.nComments + 1 ;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async deleteBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/photo/" + this.photoId , {
					headers: { 'Authorization': `Bearer ${getAuthToken()}`}
				});
				this.isDeleted = true;
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
	<div v-if="!isDeleted">
	  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"></div>
  
	  <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  
	  <div class="mb-3">
		<a @click="deleteBtn" class="photo-delete" v-if="this.isAuthor" style="cursor: pointer;">
		  Delete this photo
		</a>
  
		<div class="photo-container">
		  <img v-if="imageURL" :src="imageURL" alt="Photo" class="img-fluid" />
		</div>
		<div class="photo-caption">
		  <p>{{ caption }}</p>
		</div>
		<div class="photo-info">
		  <p>Posted by {{ username }} on {{ date }}</p>
		  <p>{{ nLikes }} Likes {{ nComments }} Comments</p>
		</div>
	  </div>
  
	 
	  <div class="photo-actions">
      <button @click="likeUnlikeBtn" :disabled="loading" class="btn btn-primary">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
        {{ isLiked ? 'Unlike' : 'Like' }}
      </button>

      <!-- Reduce the margin-left to make buttons closer together -->
      <button @click="listComments" :disabled="loading" class="btn btn-secondary" style="margin-left: 10px;">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
        Show comments
      </button>
    </div>
  
	  <div v-if="showComments">
		<div v-if="comments">
		  <div v-for="comment in comments" :key="comment.commentId">
			<Comment :commentId="comment.commentId" photoId="this.photoId"></Comment>
		  </div>
		</div>
  
		<div class="mt-3">
		  <label for="commentInput" class="form-label">Add a Comment</label>
		  <input type="text" class="form-control" id="commentInput" v-model="this.commentText" placeholder="What do you have to comment?" />
		  <button @click="commentBtn" :disabled="loading" class="btn btn-primary mt-2">Comment</button>
		</div>
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

  .photo-delete {
	font-size: 12px;
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
  