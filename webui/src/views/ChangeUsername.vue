<script>
import { getAuthUsername, getAuthToken, setAuthUsername } from '../services/tokenService';

export default {
  data: function () {
    return {
      errormsg: null,
      loading: false,
      oldUsername: getAuthUsername(),
      newUsername: "",
    }
  },

  methods: {
    async edit() {
      this.loading = true;
      this.errormsg = null;
      try {
        setAuthUsername(this.newUsername);
        await this.$axios.put("/user/" + this.oldUsername + "/username", {
          username: this.newUsername,
        }, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getAuthToken()}`,
          },
        });

        this.$router.push("/profile/" + this.newUsername);
      } catch (e) {
          if (e.response && e.response.status === 400) {
            if (this.newUsername.length < 3 || this.newUsername.length > 16) {
              this.errormsg = "The username has to be between 3 and 16 characters long...";
            } else if (this.newUsername === this.oldUsername) {
              this.errormsg = "The username has to be different from your last one!";
            } else {
              const usernameRegex = /^[a-zA-Z0-9_]+$/;
              if (!usernameRegex.test(this.newUsername)) {
                this.errormsg = "The username has too many weird characters.";
              }
            }
          } else {
            this.errormsg = e.toString();
          }
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
			<h1 class="h2">Edit username</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="mb-3">
			<label for="description" class="form-label">Pick a better username this time.</label>
			<input type="string" class="form-control" id="newUsername" v-model="this.newUsername" placeholder="yournewusername">
		</div>

		<div>
			<button v-if="!loading" type="button" class="btn btn-primary" @click="edit">
				OK
			</button>
			<LoadingSpinner v-if="loading"></LoadingSpinner>
		</div>
	</div>
</template>

<style>
</style>
