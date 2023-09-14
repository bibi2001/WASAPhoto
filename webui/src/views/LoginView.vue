<!-- 
  LoginView.vue lets the user sign in or log in.
  If the username already exists then it logs into that account, if not then it creates a new one.
  If the username does not match the expected format then a card with the error will appear.

  Endpoints called:
  POST("/session")
-->

<script>
import { setAuthToken, setAuthUsername } from '../services/tokenService';

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
        // Here we catch if the username did not have the format expected
        // Is the length between 3 and 16?
        if (this.username.length < 3 || this.username.length > 16) {
          this.errormsg = "The username has to be between 3 and 16 characters long.";
        } else {
          const usernameRegex = /^[a-zA-Z0-9_]+$/;
          // Are there characters that were not expected?
          if (!usernameRegex.test(this.username)) {
            this.errormsg = "The username has too many weird characters.";
          // If all the previous checks were passed we can send it to the backend
          } else {
            const response = await this.$axios.post("/session", {
              name: this.username,
            });

            // Save the bearer token and the username
            setAuthToken(response.data["identifier"]);
            setAuthUsername(this.username);

            // Redirect the user to their feed
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
    <div class="page-header">
      <h1 class="h2">Login</h1>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <LoadingSpinner v-if="loading"></LoadingSpinner>

    <div class="my-3">
      <label for="description" class="form-label">Try to remember your username!</label>
      <input type="text" class="form-control" id="username" v-model="username" placeholder="yourusername">
    </div>

    <div>
      <button v-if="!loading" type="button" class="btn btn-sm btn-outline-secondary" @click="login">
        Log in
      </button>
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
</style>
