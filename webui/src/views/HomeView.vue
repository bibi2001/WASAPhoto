<script>
import { getAuthToken } from '../services/tokenService';

export default {
	data: function() {
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
  					headers: { Authorization: `Bearer ${getAuthToken()}`}
				});
				this.photos = response
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

		<div class="card" v-if="photos">
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
