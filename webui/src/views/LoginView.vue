<script>
import { setAuthToken , getAuthToken, setAuthUsername, getAuthUsername} from '../services/tokenService';

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
					setAuthToken(response.data["identifier"]);
					setAuthUsername(this.username);

					// Debug 
					console.log("Received Authorization Header:", getAuthToken());
					console.log("Received Username:", getAuthUsername());

					this.$router.push("/home");
				}

			} catch (e) {
				if (response.status == 400)
					this.errormsg = "The username has to be between 3 and 16 characters long"
				else
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
			<label for="description" class="form-label">Try to remember your username!</label>
			<input type="string" class="form-control" id="username" v-model="username" placeholder="yourusername">
		</div>

		<div>
			<button v-if="!loading" type="button" class="btn btn-primary" @click="login">
				Sign in
			</button>
			<LoadingSpinner v-if="loading"></LoadingSpinner>
		</div>
	</div>
</template>

<style>
</style>
