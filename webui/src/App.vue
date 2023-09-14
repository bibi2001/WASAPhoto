<!-- 
  App.vue is responsible for the sidebar that is always visible throught the whole app.
  If the user is not authenticated (authToken == 0), the sidebar will show the log in item,
  otherwise this item is hidden.
  
  Items:
  Login - Default view of logged out user
  Home - Default view of logged in user/Redirects the user to the Home view
  Search - Redirects the user to the Search view
  Create - Redirects the user to the Upload view
  Profile - Redirects the user to their own Profile view
  Logout - Redirects the user to the Login view and sets the authToken to 0
-->

<script>
import { getAuthToken, setAuthToken, getAuthUsername } from './services/tokenService';

export default {
	computed: {
		// Retrieves wether or not the user has signed in
		isAuthenticated() {
			return getAuthToken() > 0;
		},
	},
	methods: {
		async logout() {
			setAuthToken(0);
			this.$router.push('/');
		},
		async openAuthProfile() {
			this.$router.push('/profile/' + getAuthUsername());
		},
	},
};
</script>

<template>
	<div>
		<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
			<span class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6">WASAPhoto</span>
		</header>

		<div class="container-fluid">
			<div class="row">
				<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
					<div class="position-sticky pt-3 sidebar-sticky">
						<ul class="nav flex-column">
							<div v-if="isAuthenticated">
								<li class="nav-item">
									<RouterLink to="/home" class="nav-link">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#home" />
										</svg>
										Home
									</RouterLink>
								</li>
								<li class="nav-item">
									<RouterLink to="/search" class="nav-link">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#search" />
										</svg>
										Search
									</RouterLink>
								</li>
								<li class="nav-item">
									<RouterLink to="/new" class="nav-link">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#plus-square" />
										</svg>
										Create
									</RouterLink>
								</li>
								<li class="nav-item">
									<a @click="openAuthProfile" class="nav-link" style="cursor: pointer;">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#user" />
										</svg>
										Profile
									</a>
								</li>
								<li class="nav-item">
									<a @click="logout" class="nav-link" style="cursor: pointer;">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#log-out" />
										</svg>
										Logout
									</a>
								</li>
							</div>
							<div v-else>
								<li class="nav-item">
									<RouterLink :to="'/'" class="nav-link">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#log-in" />
										</svg>
										Login
									</RouterLink>
								</li>
							</div>
						</ul>
					</div>
			</nav>
			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</div></template>

<style></style>