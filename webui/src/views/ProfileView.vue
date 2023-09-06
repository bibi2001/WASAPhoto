<script>
import { getAuthToken, getAuthUsername } from '../services/tokenService';

export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,

            username: this.$route.params.username,
			userId: 0,
            nPosts: 0,
            nFollowers: 0,
			nFollowing: 0,
            isFollowed: false,
            isBanned: false,
			isAuth: false,
			photos: [],

			followers: [],
			following: [],
			bans: [],
		}
	},
	methods: {
		load() {
			return load
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/user/" + this.username, { 
					headers: { 'Authorization': `Bearer ${getAuthToken()}`}
                });
				this.userId = response.data.userId;
            	this.nPosts = response.data.nPosts;
            	this.nFollowers = response.data.nFollowers;
				this.nFollowing = response.data.nFollowing;
            	this.isFollowed = response.data.isFollowed;
            	this.isBanned = response.data.isbanned;
				this.photos = response.data.photos;
				this.isOwner = this.userId === getAuthToken();
				
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async listFollowers() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/user/" + this.username + "/followers");
				this.followers = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async listFollowing() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get("/user/" + this.username + "/following");
				this.following = response.data;
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
		async followUnfollowBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.isFollowed) {
					await this.$axios.delete("/user/" + this.username + "/followers" + getAuthUsername());
					this.isFollowed = false;
				} else{ 
					await this.$axios.put("/user/" + this.username + "/followers" + getAuthUsername());
					this.isFollowed = true;
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async banUnbanBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.isBanned){
					await this.$axios.delete("/user/" + this.username + "/bans" + getAuthUsername());
					this.isBanned = false;
				}else{
					await this.$axios.put("/user/" + this.username + "/bans" + getAuthUsername());
					this.isBanned = true;
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">User Profile</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
				
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

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
        <div v-if="photos.length">
            <div v-for="photo in photos" :key="photo.id">
                <photo :photoId="photo.id"></photo>
            </div>
        </div>
        <div class="card" v-else>
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
