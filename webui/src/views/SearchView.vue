<!-- 
  SearchView.vue lets the user search for other user profiles.
  If there are no matches, it shows a card with the text "No user match."
  This search does not show the banned profiles!

  Endpoints called:
  GET("/users/search")
-->

<script>
import { getAuthToken } from '../services/tokenService';

export default {
  data: function () {
    return {
      errormsg: null,
      loading: false,
      
      query: "",
      usernames: [],
    }
  },
  methods: {
    async search() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get("/users/search?q=" + this.query, {
          headers: { Authorization: `Bearer ${getAuthToken()}` }
        });
        this.usernames = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
  },
}
</script>

<template>
  <div>
    <div class="page-header">
      <h1 class="h2">Search</h1>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <div class="my-3">
      <label for="description" class="form-label"> Look someone up!</label>
      <input type="string" class="form-control" id="query" v-model="query" placeholder="username" @input="search">
    </div>

    <LoadingSpinner v-if="loading"></LoadingSpinner>


    <div v-if="usernames">
      <div v-for="username in usernames" :key="username">
        <UserList :username="username"></UserList>
      </div>
    </div>
    <div class="card" v-else>
      <div class="card-body">
        <p>No users match.</p>
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
