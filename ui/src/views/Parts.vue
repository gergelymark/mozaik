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
  import lunr from 'lunr'
  export default {
    name: 'PartsView',
    data: () => ({
        loading: true,
        search: "",
        searchIndex: {},
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
        index() {
            const that = this
            return lunr(function () {
                this.ref('ID')                
                this.field('color')
                this.field('ItemName')
                this.field('size')
                this.field('type')
                this.field('subtype')
                that.allParts.forEach(function (doc) {
                    let type, size, subtype = ""
                    const res = doc.ItemName.match(/(.*)(\d+ x \d+)(.*)/)
                    if (res) {
                        type = res[1].trim()
                        size = res[2].trim().replaceAll(" ", "")
                        subtype = res[3].trim()
                    } 
                    console.log(doc.ItemName.match(/(.*)(\d+ x \d+)(.*)/))
                    this.add({
                        ID: doc.ID,
                        color: doc.ColorName,
                        ItemName: doc.ItemName,
                        type: type,
                        size: size,
                        subtype: subtype,
                    })
                }, this)
            })
        },
        allParts() {
            return this.$store.state.parts
        },
        parts() {                                    
            if (!this.search) {
                return this.allParts
            }
            return this.index.search(
                `${this.search}`
                ).map(
                    result => this.allParts.find(
                        element => element.ID === result.ref
                )
            )            
            // return this.index.search(
            //     `ColorName:${this.search}^10 `+
            //     `Type:${this.search}^6` +
            //     `Subtype:${this.search}^4` +
            //     `Size:${this.search}^12` +
            //     `${this.search}^0.1`
            //     ).map(
            //         result => this.allParts.find(
            //             element => element.ID === result.ref
            //     )
            // )            
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