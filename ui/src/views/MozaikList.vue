<template>
  <v-container>
    <v-card
      :loading="loading"
      :disabled="loading"
    >    
    <v-list two-line>
      <v-list-item
        v-for="mozaik in mozaiks" :key="`mozaik-${mozaik.Name}`"
        :to="`/mozaik/${mozaik.Name}`"
      >
        <v-list-item-avatar size="56">
          <v-img :src="`/api/images/${mozaik.Name}`" />
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title v-html="mozaik.Name" />
          <v-list-item-subtitle v-html="`${mozaik.Width} x ${mozaik.Height}`"/>
        </v-list-item-content>     
      </v-list-item>
      <v-list-item 
        to="/create"
        @click.native="$store.commit('createRequested', Date.now())"
      >
        <v-list-item-avatar size="56">
          <v-icon size="56">mdi-plus</v-icon>
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>Create new mozaik</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list> 
    </v-card>
  </v-container>
</template>

<script>
  export default {
    name: 'MozaikList',
    data: () => ({
      loading: true,
    }),
    components: {
    },
    computed: {
      mozaiks() {
        return this.$store.state.mozaiks
      }
    },
    async mounted() {
      await this.$store.dispatch("fetchAllMozaik")
      this.loading = false
    }
  }
</script>
