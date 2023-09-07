<script>
import { setAuthToken, getAuthToken, setAuthUsername } from '../services/tokenService';

export default {
  data: function () {
    return {
      errormsg: null,
      loading: false,
      username: '',
    }
  },

  methods: {
    async login() {
      this.loading = true;
      this.errormsg = null;
      try {
        if (this.username.length < 3 || this.username.length > 16) {
          this.errormsg = "The username has to be between 3 and 16 characters long.";
        } else {
          const usernameRegex = /^[a-zA-Z0-9_]+$/;
          if (!usernameRegex.test(this.username)) {
            this.errormsg = "The username has too many weird characters.";
          } else {
            const response = await this.$axios.post("/session", {
              name: this.username,
            });

            setAuthToken(response.data["identifier"]);
            setAuthUsername(this.username);

            this.$router.push("/home");
          }
        }
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
      <label for="description" class="form-label">Try to remember your username!</label>
      <input type="text" class="form-control" id="username" v-model="username" placeholder="yourusername">
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
