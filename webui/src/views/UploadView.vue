<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Post photo</h1>
    </div>

    <input type="file" ref="fileInput" @change="onFileSelected" accept="image/*">
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
    onFileSelected(event) {
      const file = event.target.files[0];
      const reader = new FileReader();

      if (file) {
        reader.onload = (evt) => {
          const binaryString = evt.target.result;
          this.selectedImage = URL.createObjectURL(file);
          // Convert binaryString to Array or send it directly
          this.fileArrayBuffer = new Uint8Array(binaryString.length);
          for (let i = 0; i < binaryString.length; i++) {
            this.fileArrayBuffer[i] = binaryString.charCodeAt(i);
          }
        };
        reader.readAsBinaryString(file);
      }
    },
    async uploadImage() {
      if (!this.fileArrayBuffer) {
        return;
      }

      try {
        await this.$axios.post("/photo", {
          image: Array.from(this.fileArrayBuffer),
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
