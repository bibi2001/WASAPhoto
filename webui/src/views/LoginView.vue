<script>
import { setAuthToken , getAuthToken} from '../services/tokenService';

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			
			// field 
			username: "",
		}
	},
	methods: {
		async login () {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.post("/session", {
					name: this.username,
				});

				if (response.status == 201) {
					setAuthToken(response.data["identifier"])

					// Log the received Authorization header
					console.log("Received Authorization Header:", getAuthToken());
				}
				this.$router.push("/");

			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		}
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Login</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="mb-3">
			<label for="description" class="form-label">Username</label>
			<input type="string" class="form-control" id="username" v-model="username" placeholder="yourusername">
		</div>

		<div>
			<button v-if="!loading" type="button" class="btn btn-primary" @click="login">
				SIGN IN
			</button>
			<LoadingSpinner v-if="loading"></LoadingSpinner>
		</div>
	</div>
</template>

<style>
</style>
