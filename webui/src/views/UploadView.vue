<!-- 
  UploadView.vue lets the user choose a photo, write a caption and upload it to their profile.
  If no photo is selected and the user tries to post it, it shows a card with the text "Please select an image to upload."
  
  Endpoints called:
  POST("/photo")
-->

<script>
import { getAuthToken, getAuthUsername } from '../services/tokenService';

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
      this.uploadFile = files[0]; 
    },


    async uploadImage() {
      // Check if no file selected
      if (!this.uploadFile) {        
        this.errormsg = "Please select an image to upload.";
        return;
      }
      try {
        const reader = new FileReader();
        // Read the selected file as a data URL
        reader.readAsDataURL(this.uploadFile);
        reader.onload = async () => {
          const imageData = reader.result;
          await this.$axios.post("/photo", {
            image: imageData,
            caption: this.caption,
          }, {
            headers: {
              'Content-Type': 'application/json', // Set content type to JSON
              'Authorization': `Bearer ${getAuthToken()}`,
            },
          });

          // Redirect user to their own profile
          this.$router.push('/profile/' + getAuthUsername());
        };
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

    <div class="page-header">
      <h1 class="h2">Post photo</h1>
    </div>

    <div class="my-3">
      <input type="file" ref="fileInput" @change="onFileSelected" accept="image/*">
      <img v-if="selectedImage" :src="selectedImage" alt="Selected Image" />
    </div>

    <div class="my-3">
      <div class="input-group my-3">
        <input type="text" class="form-control" id="caption" v-model="this.caption"
          placeholder="Write your caption here" />
      </div>
      <button class="btn btn-sm btn-outline-secondary" @click="uploadImage">Upload Image</button>
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
