<script>

export default {
  data: function() {
    return {
      errormsg: null,
      loading: false,
      username: this.$route.params.username,
      following: [],
    }
  },
  methods: {
    async refresh() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get("/user/" + this.username + "/following");
        console.log(response.data);
        this.following = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
  },
  mounted() {
    // Load data automatically when the component is mounted.
    this.refresh();
  },
}

</script>

<template>
    <div>
	  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
		<h1 class="h2">Following</h1>
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


        <div v-if="this.following">
            <div v-for="username in this.following" :key="username">
                <UserList :username="username"></UserList>
            </div>
        </div>
        <div class="card" v-else>
            <div class="card-body">
                <p>No followings here!</p>
            </div>
        </div>
    </div>
</template>

<style scoped>
.card {
  margin-bottom: 20px;
}
</style>
