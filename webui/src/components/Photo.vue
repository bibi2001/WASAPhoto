<!-- 
  Photo.vue is the common component of the list of photos.
  If the user is the author of the photo, they can delete the photo.
  In this file we fetch the image that corresponds to the photo post and show it to the user,
  along with the number of likes, comments and list of comments of the photo.

  Endpoints called:
  GET("/photo/:photoId")
  GET("/photo/:photoId/comments")
  PUT("/photo/:photoId/likes/:username")
  DELETE("/photo/:photoId/likes/:username")
  POST("/photo/:photoId/comments")
  DELETE("/photo/:photoId")
-->

<script>
import { getAuthToken, getAuthUsername } from '../services/tokenService';

// Auxiliar function to convert the dataURL string to a blob
function dataURLtoBlob(dataURL) {
	const parts = dataURL.split(',');
	const mimeType = parts[0].match(/:(.*?);/)[1];
	const b64Data = atob(parts[1]);
	const byteArray = new Uint8Array(b64Data.length);

	for (let i = 0; i < b64Data.length; i++) {
		byteArray[i] = b64Data.charCodeAt(i);
	}

	return new Blob([byteArray], { type: mimeType });
}

export default {
	props: ["photoId"],
	data: function () {
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
					headers: { 'Authorization': `Bearer ${getAuthToken()}` }
				},);
				this.username = response.data.photo.username;
				this.date = response.data.photo.date;
				this.caption = response.data.photo.caption;
				this.nComments = response.data.photo.nComments;
				this.nLikes = response.data.photo.nLikes;
				this.isLiked = response.data.photo.isLiked;
				this.isAuthor = this.username === getAuthUsername();

				// Convert Blob to a data URL
				this.imageURL = URL.createObjectURL(dataURLtoBlob(response.data.dataURL));

			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

		async listComments() {
			this.loading = true;
			this.errormsg = null;
			try {
				if (!this.showComments) {
					const response = await this.$axios.get("/photo/" + this.photoId + "/comments");
					this.comments = response.data;
					this.showComments = true;
				} else {
					this.showComments = false
				}
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
					this.nLikes = this.nLikes - 1;
				} else {
					await this.$axios.put("/photo/" + this.photoId + "/likes/" + getAuthUsername());
					this.isLiked = true;
					this.nLikes = this.nLikes + 1;
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
					text: this.commentText
				}, {
					headers: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${getAuthToken()}`,
					},
				});

				this.showComments = false;
				this.nComments = this.nComments + 1;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async deleteBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.delete("/photo/" + this.photoId, {
					headers: { 'Authorization': `Bearer ${getAuthToken()}` }
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
		<div
			class="photo-separator">
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="mb-3">
			<a @click="deleteBtn" class="photo-delete" v-if="this.isAuthor" style="cursor: pointer;">
				Delete this photo
			</a>

			<div class="photo-container">
				<img v-if="imageURL" :src="imageURL" alt="Photo" class="img-fluid" />
			</div>
			<div class="my-1">
				<p class="photo-info">{{ nLikes }} Likes {{ nComments }} Comments</p>
				<p class="photo-caption">{{ caption }}</p>
				<p class="photo-info">Posted by {{ username }} on {{ date }}</p>
			</div>
			
		</div>


		<div class="photo-actions">
			<button @click="likeUnlikeBtn" :disabled="loading" 
				class="btn btn-sm btn-outline-secondary">
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#heart" />
				</svg>
				{{ isLiked ? 'Unlike' : 'Like' }}
			</button>

			<!-- Reduce the margin-left to make buttons closer together -->
			<button @click="listComments" :disabled="loading" 
				class="btn btn-sm btn-outline-secondary">
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#message-circle" />
				</svg>
				{{ showComments ? 'Hide comments' : 'Show comments' }}
			</button>
		</div>

		<div v-if="showComments" class="mt-3">
			<div class="card" v-if="comments">
				<div v-for="comment in comments" :key="comment">
					<Comment :comment="comment" :photoId=this.photoId></Comment>
				</div>
			</div>

			<div class="mt-1">
				<label for="commentText" class="form-label">What do you have to comment?</label>
				<div class="input-group">
					<input type="text" class="form-control" id="commentText" v-model="this.commentText"
						placeholder="Cool pic!" />
					<button @click="commentBtn" :disabled="loading" 
						class="btn btn-sm btn-outline-secondary">Comment</button>
				</div>
			</div>
		</div>

	</div>
</template>
  
  
  
<style scoped>
.photo-separator { 
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  padding-top: 1rem;
  padding-bottom: 2px;
  margin-bottom: 3px;
  border-bottom: 1px solid #ccc;
}
.card {
  margin-bottom: 20px;
}
.photo-container {
	text-align: center;
}

.photo-caption {
	margin-top: 34px;
	margin-bottom: 30px;
	height: 0;
}

.photo-info {
	font-size: 14px;
	color: #777;
	height: 0px;
	margin-bottom: 30px;
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

.username {
	cursor: pointer;
	font-size: 13px;
	font-weight: bold;
	height: 0px;
}</style>
  