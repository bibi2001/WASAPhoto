<!-- 
  Bans.vue gets the bans of the user with the username in the parameter.
  If there are no bans, it shows a card with the text "You need to ban some people!".
  If there are any errors when getting the list, it shows a card with the error (stored in errormsg).

  Endpoints called:
  GET("/user/:username/bans")
-->

<script>
import { getAuthToken } from '../services/tokenService';

export default {
  data: function() {
    return {
      errormsg: null,
      loading: false,
      
      username: this.$route.params.username,
      bans: [],
    }
  },
  methods: {
    async refresh() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get("/user/" + this.username + "/bans", {
					headers: { Authorization: `Bearer ${getAuthToken()}` }
				});
        this.bans = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
  },
  mounted() {
    this.refresh();
  },
}
</script>

<template>
  <div>

	  <div class="page-header">
		  <h1 class="h2">Bans</h1>
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
    
    <div v-if="this.bans">
      <div v-for="username in this.bans" :key="username">
        <UserList :username="username"></UserList>
      </div>
    </div>
    <div class="card mt-3" v-else>
      <div class="card-body">
        <p>You need to ban some people!</p>
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
