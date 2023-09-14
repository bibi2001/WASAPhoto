<!-- 
  ChangeUsername.vue lets the user change their username.
  If the new username does not match the expected format then a card with the error will appear.

  Endpoints called:
  PUT("/user/:username/username")
-->

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
        // Here we catch if the new username did not have the format expected
        if (e.response && e.response.status === 400) {
          // Is the length between 3 and 16?
          if (this.newUsername.length < 3 || this.newUsername.length > 16) {
            this.errormsg = "The username has to be between 3 and 16 characters long...";
          // Is it not different from the previous username?
          } else if (this.newUsername === this.oldUsername) {
            this.errormsg = "The username has to be different from your last one!";
          // Are there characters that were not expected?
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

    <div class="page-header">
      <h1 class="h2">Edit username</h1>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <LoadingSpinner v-if="loading"></LoadingSpinner>
    
    <div class="my-3">
      <label for="description" class="form-label">
        Pick a better username this time.</label>
      <input type="string" class="form-control" id="newUsername" v-model="this.newUsername" placeholder="newusername">
    </div>

    <div>
      <button v-if="!loading" type="button" class="btn btn-sm btn-outline-secondary" @click="edit">
        OK
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
