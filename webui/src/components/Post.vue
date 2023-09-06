
<script>
import { getAuthToken, getAuthUsername } from '../services/tokenService';

export default {
    props: ["photoId"],
	data: function() {
		return {
			errormsg: null,
			loading: false,
            
	        username: null,
            image: 0,
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
            	this.image = response.data.image;
            	this.date = response.data.date;
				this.caption = response.data.caption;
            	this.nComments = response.data.nComments;
            	this.nLikes = response.data.nLikes;
				this.isLiked = response.data.isLiked;
                this.isAuthor = this.username === getAuthUsername();
				
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
		async listBans(){
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/user/" + this.username + "/bans");
				this.following = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async likeUnlikeBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.isFollowed) {
					await this.$axios.delete("/user/" + this.photoId + "/likes" + getAuthUsername());
					this.isLikes = false;
				} else{ 
					await this.$axios.put("/user/" + this.photoId + "/likes" + getAuthUsername());
					this.isLikes = true;
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
				await this.$axios.post("/user/" + this.photoId + "/comments", {
					text: this.commentText, headers: { 'Authorization': `Bearer ${getAuthToken()}`}
				});
				await listComments();
				
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

		<LoadingSpinner v-if="loading"></LoadingSpinner>

		<div class="d-flex align-items-center">
			<h1 class="h2">{{ username }}</h1>
				<div class="ms-5">
					<div class="btn-group me-4">
						<button type="button" class="btn btn-sm btn-outline-secondary" @click="followUnfollowBtn">
							Follow
						</button>
					</div>
				</div>
				<div class="ms-0">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="banUnbanBtn">
						Block
					</button>
				</div>
			</div>
		<div class="d-flex">
			<p class="me-5">{{ nPosts }} Posts</p>
			<p class="me-5">{{ nFollowers }} Followers</p>
			<p>{{ nFollowing }} Following</p>
		</div>

		<div class="card" v-if="!photos">
			<div class="card-body">
				<p>No photos to show.</p>
			</div>
		</div>

		
	</div>
</template>

<style scoped>
.card {
	margin-bottom: 20px;
}
</style>
