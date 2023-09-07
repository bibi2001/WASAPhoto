<script>
import { getAuthToken } from '../services/tokenService';

export default {
  data() {
    return {
      errormsg: null,
      loading: false,
      uploadFile: null, 
    };
  },
  methods: {
    // This function is called when a file is selected
    onFileSelected(event) {
      let files = event.target.files || event.dataTransfer.files;
      if (!files.length) return;
      this.uploadFile = files[0]; // Store the selected file
    },

    async uploadImage() {
      if (!this.uploadFile) {
        // Check if no file selected
        this.errormsg = "Please select an image to upload.";
        return;
      }

      try {
        const reader = new FileReader();

        // Read the selected file as a data URL
        reader.onload = async () => {
          const imageData = reader.result; 

          await this.$axios.post("/photo", {
            image: imageData, 
            caption: "caption!!",
          }, {
            headers: {
              'Content-Type': 'application/json', // Set content type to JSON
              'Authorization': `Bearer ${getAuthToken()}`,
            },
          });

          this.$router.push("/home");
        };

        reader.readAsDataURL(this.uploadFile); // Read the file as a data URL
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.errormsg = e.response.data.message; // Display error message from the server
        } else {
          this.errormsg = "An error occurred. Please try again later.";
        }
      }
    },
  },
};
</script>


<template>
  
  <div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Post photo</h1>
    </div>

    
    <input type="file" ref="fileInput" @change="onFileSelected" accept="image/*">
    <img v-if="selectedImage" :src="selectedImage" alt="Selected Image" />
    <button @click="uploadImage">Upload Image</button>
  </div>
</template>

<style scoped>
/* Add your CSS styles here */
</style>
