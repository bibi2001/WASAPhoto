<!-- 
  HomeView.vue shows the feed of the user.
  The feed is composed by the photos of the followed users of the user.
  If the user doesn't follow anyone, it shows a card with the text "No photos to show."
  
  Endpoints called:
  GET("/home")
-->

<script>
import { getAuthToken } from '../services/tokenService';

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			
			photos: [],
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
				const response = await this.$axios.get("/home", {
					headers: { Authorization: `Bearer ${getAuthToken()}` }
				});
				this.photos = response.data
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
		<div class="page-header">
			<h1 class="h2">Your feed</h1>
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

		<div v-if="photos">
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
</style>
