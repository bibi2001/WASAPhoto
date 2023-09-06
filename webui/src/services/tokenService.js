import { reactive } from 'vue';

// Create a reactive token service
const tokenService = reactive({
  token: 0, 
  username: "",
});

// Function to set the token
function setAuthToken(token) {
  tokenService.token = token;
}

// Function to set the username
function setAuthUsername(username) {
  tokenService.username = username
}

// Function to get the token
function getAuthToken() {
  return tokenService.token;
}

// Function to get the username
function getAuthUsername() {
  return tokenService.username;
}

export { setAuthToken, getAuthToken, setAuthUsername, getAuthUsername };
