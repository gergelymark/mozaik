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
                show-expand
                item-key="ID"
                single-expand
                :items-per-page="-1"
                :hide-default-footer="true"
                :headers="headers"
                :items="parts"
                :search="search"
            >
            <template
                v-slot:item="{ item, expand, isExpanded }"
            >
                <tr>
                    <td style="text-align:center"><v-img width="64" :src="`https:${item.ImgURL}`" /></td>
                    <td v-html="item.WantedQty" />
                    <td v-html="item.InStock" />
                    <td v-html="item.ColorName" />
                    <td v-html="item.ItemName" />
                    <td>
                        <v-btn 
                            icon
                            text
                            @click="expand(!isExpanded)"
                        >
                            <v-icon v-html="isExpanded ? 'mdi-chevron-up' : 'mdi-chevron-down'" />
                        </v-btn>                            
                    </td>
                </tr>
            </template>
            <template v-slot:expanded-item="{ headers, item }">
                <td :colspan="headers.length">
                    <div class="inset">
                        <table>
                            <tr v-for="legoSet in item.LegoSets"
                                :key="`set-${legoSet.ID}`">
                                <td v-html="legoSet.Name" />
                                <td v-html="legoSet.PartQty" class="text-right"/>
                            </tr>
                        </table>
                    </div>                    
                </td>
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
          { text: 'In Stock', value: 'InStock', width: "150px" },
          { text: 'Color', value: 'ColorName' },
          { text: 'Name', value: 'ItemName' },
          { text: '', value: 'data-table-expand' },
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
<style scoped>
    .inset {
        box-shadow: inset 0 4px 8px -5px rgb(50 50 50 / 75%), inset 0 -4px 8px -5px rgb(50 50 50 / 75%) ! important;
        padding: 10px;
    }
    .inset td {
        padding: 5px;
        padding-left: 20px;
    }
</style>