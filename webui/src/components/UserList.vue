<script>
export default {
    props: ["username"],
	data: function() {
		return {
			errormsg: null,
			loading: false,
		}
	},
    methods: {
		async goToProfile() {
			this.$router.push("/profile/"+this.username);
		},
    }, 
}
</script>

<template>
	<div class="card" @click="goToProfile" style="cursor: pointer;">
	  <div class="card-body">
		<p class="username">{{ username }}</p>
	  </div>
	</div>
</template>
  
<style scoped>
.card {
	margin-bottom: 10px;
	height:60px;
}
.username {
	font-size: 18px;
}
</style>
