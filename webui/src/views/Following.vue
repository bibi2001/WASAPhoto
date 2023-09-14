<!-- 
  Following.vue gets the followed users of the user with the username in the parameter.
  If there are no followed users, it shows a card with the text "No followings here!".
  If there are any errors when getting the list, it shows a card with the error (stored in errormsg).

  Endpoints called:
  GET("/user/:username/following")
-->

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
        this.following = response.data;
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
    <div class="card mt-3" v-else>
      <div class="card-body">
        <p>No followings here!</p>
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

