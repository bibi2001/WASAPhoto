<script>
import { getAuthToken } from '../services/tokenService';

export default {
  data: function() {
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
        console.log(response.data);
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
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
		<h1 class="h2">Search</h1>
	</div>

	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

	<div class="mb-3">
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
.card {
  margin-bottom: 20px;
}
</style>
