<template>    
    <v-container>
        <v-card
            :loading="loading"
            :disabled="loading"
        >
            <v-card-title>
                <v-text-field
                    v-model="search"
                    append-icon="mdi-magnify"
                    label="Search"
                    single-line
                    hide-details
                ></v-text-field>
            </v-card-title>
            <v-data-table
                :items-per-page="-1"
                :hide-default-footer="true"
                :headers="headers"
                :items="parts"
                :search="search"
            >
            <template
                v-slot:body="{ items }"
            >
                <tbody>
                    <tr
                        v-for="item in items"
                        :key="item.name"
                    >
                        <td style="text-align:center"><v-avatar><v-img :src="`https:${item.ImgURL}`" /></v-avatar></td>
                        <td v-html="item.WantedQty" />
                        <td v-html="item.ColorName" />
                        <td v-html="item.ItemName" />
                    </tr>
                </tbody>
            </template>

            </v-data-table>
        </v-card>
    </v-container>
</template>
<script>
  export default {
    name: 'PartsView',
    data: () => ({
        loading: true,
        search: "",
        headers: [
          {
            text: '',
            align: 'start',
            filterable: false,
            value: 'ImgURL',
            width: "150px"
          },          
          { text: 'Amount', value: 'WantedQty', width: "150px" },
          { text: 'Color', value: 'ColorName' },
          { text: 'Name', value: 'ItemName' },
        ],
    }),
    computed: {
        parts() {
            return this.$store.state.parts
        }
    },
    async mounted() {
        await this.$store.dispatch("fetchParts")
        this.loading = false
    },
    methods: {
    }
  }
</script>