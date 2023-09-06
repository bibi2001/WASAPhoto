<template>
    <div>
      <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">Post photo</h1>
      </div>
  
      <input type="file" ref="fileInput" @change="handleFileChange" accept="image/*">
      <img v-if="selectedImage" :src="selectedImage" alt="Selected Image" />
      <button @click="uploadImage">Upload Image</button>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        selectedImage: null,
        fileArrayBuffer: null,
      };
    },
    methods: {
      handleFileChange(event) {
        // Corrected: Use event.target.files[0] directly
        const file = event.target.files[0];
        const reader = new FileReader();
  
        reader.onload = (e) => {
          // Convert the file to an ArrayBuffer
          this.fileArrayBuffer = e.target.result;
          this.selectedImage = URL.createObjectURL(file); // Use 'file' here
        };
  
        reader.readAsArrayBuffer(file); // Use 'file' here
      },
      async uploadImage() {
        if (!this.fileArrayBuffer) {
          return;
        }
  
        try {
          await this.$axios.post("/photo", {
            image: Array.from(new Uint8Array(this.fileArrayBuffer)),
            caption: "caption!!",
          }, {
            headers: {
              'Authorization': `Bearer ${getAuthToken()}`,
            },
          });
  
          this.$router.push("/home");
        } catch (e) {
          this.errormsg = e.toString();
        }
      },
    },
  };
  </script>
  
  <style scoped>
  /* Add your CSS styles here */
  </style>
  