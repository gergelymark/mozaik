<template>
  <v-container>
    <v-card
      :disabled="processing"
      :loading="processing"
    >
      <v-card-title>
        Create new mozaik
      </v-card-title>
      <v-card-text>
        <v-container>     
          <v-row>
            <v-col cols="6">
              <v-text-field
                persistent-hint
                outlined
                label="Name"
                v-model="name"
              />
            </v-col>
            <v-col cols="3">
              <v-text-field
                persistent-hint
                outlined
                label="Width"
                v-model="width"
                type="number"
              />
            </v-col>
            <v-col cols="3">
              <v-text-field
                persistent-hint
                outlined
                label="Height"
                v-model="height"
                type="number"
              />
            </v-col>
          </v-row>           
          <v-form>
            <cropper
              class="cropper"
              :src="$store.state.imageToCrop"
              :stencil-props="{
                aspectRatio: width/height
              }"
              @change="change"
            />
          
          <v-btn
            block 
            text
            class="mt-6"     
            @click="createMozaik"
          >
            <v-icon>mdi-upload</v-icon>upload
          </v-btn>
          </v-form>
        </v-container>
      </v-card-text>
    </v-card>   
  </v-container> 
</template>

<script>
  import { Cropper } from 'vue-advanced-cropper'
  import 'vue-advanced-cropper/dist/style.css';
  export default {
    name: 'create-mozaik',
    components: {
      Cropper
    },
    data: () => ({
      image: null,
      name: "",
      width: 48,
      height: 48,
      processing: false,
    }),    
    computed: {
      imageToUploadName() {
        return this.$store.state.imageToCropName
      }
    },
    watch: {
      imageToUploadName(nV) {
        this.name = nV
      },
    },
    methods: {
      change({ canvas }) {
        this.image = canvas.toDataURL();
      },
      async createMozaik() {
        this.processing = true
        let mozaikData = {
          name: this.name,
          dataUrl: this.image,
          width: parseInt(this.width),
          height: parseInt(this.height),
        }
        await this.$store.dispatch("createMozaik", mozaikData)        
        this.processing = false
        this.$router.push({name: "mozaik", params: {"name": this.name}})
      }
    }
  }
</script>
<style scoped>
  .cropper {
    min-height: 300px;
		max-height: 500px;
  }
</style>