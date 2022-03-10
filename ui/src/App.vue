<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      clipped-left
      dark
    >
      <div class="d-flex align-center">
        <v-img 
          src="/Mozaik-logos_white.png" 
          width="48"
        />
      </div>
      <v-btn
        icon
        to="/"
        class="ml-6"        
      >
        <v-icon>mdi-home</v-icon>
      </v-btn>
      <v-btn
        text
        to="/create"
        @click="$store.commit('createRequested', Date.now())"
        :loading="isSelecting"
      >
        <v-icon>mdi-plus</v-icon> New mozaik
      </v-btn>
      <input 
        ref="uploader"
        type="file"
        class="d-none"
        accept="image/*"
        @change="onFileChanged"
      />
      <v-spacer></v-spacer>
    </v-app-bar>
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: 'App',
  data: () => ({
    isSelecting: false,
  }),
  computed: {
    createRequested: {
      get() {
        return this.$store.state.createRequested
      },
      set(v) {
        this.$store.commit("createRequested", v)
      }
    },
    imageToUpload: {
      get() {
        return this.$store.state.imageToCrop
      },
      set(v) {
        this.$store.commit("imageToCrop", v)
      },
    },
    imageToUploadName: {
      get() {
        return this.$store.state.imageToCropName
      },
      set(v) {
        this.$store.commit("imageToCropName", v.replace(/\.[^.$]+$/, ''))
      },
    }
  },
  watch: {
    createRequested(nV) {
      if (!nV) {
        return
      }
      this.isSelecting = true
      window.addEventListener('focus', () => {
        this.isSelecting = false
      }, { once: true })

      this.$refs.uploader.click()
    }
  },
  mounted() {
    this.$store.dispatch("getColors")
  },  
  methods: {
    onFileChanged(e) {
      const selectedFile = e.target.files[0]
      this.imageToUploadName = selectedFile.name
      this.imageToUpload = URL.createObjectURL(selectedFile)
    },
  },
};
</script>
