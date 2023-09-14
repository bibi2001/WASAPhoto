<!-- 
  ProfileView.vue shows the profile of a user.
  The feed is composed by the photos of the profile owner.
  If the owner does not have any photos or is banned, it shows a card with the text "No photos to show."
  If the profile is not owned by the user, the user can follow the owner, ban the owner and interact with the 
  photos. Otherwise,the user can change their username and check the list of users that they banned.
  
  Endpoints called:
  GET("/user/:username")
  PUT("/user/:username/followers/:authUser")
  DELETE("/user/:username/followers/:authUser")
  PUT("/user/:username/bans/:bannedUser")
  DELETE("/user/:username/bans/:bannedUser")
-->

<script>
import { getAuthToken, getAuthUsername } from '../services/tokenService';

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,

			username: this.$route.params.username,
			userId: 0,
			photos: [],

			nPosts: 0,
			nFollowers: 0,
			nFollowing: 0,

			isFollowed: false,
			isBanned: false,
			isAuth: false,
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
				const response = await this.$axios.get("/user/" + this.username, {
					headers: { 'Authorization': `Bearer ${getAuthToken()}` }
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
		listFollowers() {
			this.$router.push("/profile/" + this.username + "/followers");
		},
		listFollowing() {
			this.$router.push("/profile/" + this.username + "/following");
		},
		listBans() {
			this.$router.push("/profile/" + this.username + "/bans");
		},
		changeUsername() {
			this.$router.push("/profile/" + this.username + "/username");
		},
		async followUnfollowBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.isFollowed) {
					await this.$axios.delete("/user/" + this.username + "/followers/" + getAuthUsername());
					this.isFollowed = false;
					this.nFollowers = this.nFollowers - 1;
				} else {
					await this.$axios.put("/user/" + this.username + "/followers/" + getAuthUsername());
					this.isFollowed = true;
					this.nFollowers = this.nFollowers + 1;
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
			this.refresh();
		},
		async banUnbanBtn() {
			this.loading = true;
			this.errormsg = null;
			try {
				if (this.isBanned) {
					await this.$axios.delete("/user/" + getAuthUsername() + "/bans/" + this.username);
					this.isBanned = false;
				} else {
					await this.$axios.put("/user/" + getAuthUsername() + "/bans/" + this.username);
					this.isBanned = true;
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
			this.refresh();
		},
	},
	mounted() {
		this.create();
	}
}
</script>

<template>
	<div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<LoadingSpinner v-if="loading"></LoadingSpinner>

		<div class="page-header">
			<h1 class="h2">{{ username }}</h1>
			<div v-if="!this.isOwner" class="ms-5">
				<div class="btn-group me-4">
					<div cv-if="!isBanned">
						<button type="button" class="btn btn-sm btn-outline-secondary" @click="followUnfollowBtn">
							{{ this.isFollowed ? "Unfollow" : "Follow" }}
						</button>
					</div>
					<button type="button" class="btn  ms-4 btn-sm btn-outline-secondary" @click="banUnbanBtn">
						{{ this.isBanned ? "Unban" : "Ban" }}
					</button>
				</div>
			</div>
			<div v-else class="d-flex align-items-center ms-2">
				<a @click="changeUsername" class="nav-link" style="cursor: pointer;">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#edit" />
					</svg>
				</a>
				<button type="button" class="btn btn-sm btn-outline-secondary ms-5" @click="listBans">
					{{ "Bans" }}
				</button>
			</div>
		</div>

		<div class="d-flex mt-3">
			<p class="me-5">
				{{ nPosts }} Posts</p>
			<p class="me-5" @click="listFollowers" style="cursor: pointer;">
				{{ nFollowers }} Followers</p>
			<p @click="listFollowing" style="cursor: pointer;">
				{{ nFollowing }} Following</p>
		</div>

		<div v-if="photos && !isBanned">
			<div v-for="photo in photos" :key="photo.photoId">
				<Photo :photoId="photo.photoId"></Photo>
			</div>
		</div>
		<div class="card mt-3" v-else>
			<div class="card-body">
				<p>No photos to show.</p>
			</div>
		</div>

	</div>
</template>
  
<style scoped>
.page-header {
	display: flex;
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
</style>
