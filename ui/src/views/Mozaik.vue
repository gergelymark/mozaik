<template>    
    <v-sheet>
        <v-row>
            <v-col cols="2">
                <v-sheet 
                    style="height: calc(100vh - 64px);overflow: auto"
                >
                    <v-text-field
                        outlined
                        dense
                        class="ma-2 ml-4"      
                        hide-details
                        prepend-inner-icon="mdi-magnify"          
                        v-model="searchColor"
                    />
                    <v-list>
                        <v-list-item-group v-model="selectedColor">
                            <v-list-item v-for="(color, hex) in colors" :key="`color-${hex}`">
                                <v-list-item-icon><v-icon :color="hex">mdi-toy-brick</v-icon></v-list-item-icon>
                                <v-list-item-content>{{color.name}}</v-list-item-content>
                            </v-list-item>
                        </v-list-item-group>
                    </v-list>
                </v-sheet>
            </v-col>                        
            <v-col cols="9">
                <v-app-bar
                    dense
                    flat
                    rounded
                    color="grey lighten-1"
                    class="mt-2"
                >
                    <v-row align="center" justify="center">
                        <v-col cols="6">
                            <v-btn icon
                                :disabled="lastChangeIndex === 0"
                                @click="save"
                            >
                                <v-icon>mdi-content-save</v-icon>
                            </v-btn>
                            <v-btn 
                                icon 
                                :disabled="lastChangeIndex === 0"
                                @click="undo"
                            >
                                <v-icon>mdi-undo</v-icon>
                            </v-btn>
                            <v-btn 
                                icon 
                                :disabled="lastChangeIndex == changes.length"
                                @click="redo"
                            >
                                <v-icon>mdi-redo</v-icon>
                            </v-btn>
                        </v-col>
                        <v-col cols="6">
                            <v-row align="center" justify="center">
                            <div>{{hoverColorName}} <v-icon v-if="hoverColorName" :color="hoverColor">mdi-toy-brick</v-icon></div>                            
                            <v-spacer/>
                            <v-btn icon
                                color="error"
                                @click="remove"
                                class="mr-2"
                            >
                                <v-icon>mdi-delete</v-icon>
                            </v-btn>
                            </v-row>
                        </v-col>
                    </v-row>
                </v-app-bar>
                <v-card
                    :loading="loading"
                    :disabled="loading"
                    class="mt-2"
                >
                    <v-card-text>
                        <v-row 
                            v-if="mozaik.Name" 
                            align="center" 
                            justify="center"
                        >
                            <v-col 
                                cols="6" 
                                align-self="center"
                                class="text-center"
                            >
                                <v-img 
                                    width="480" 
                                    height="480" 
                                    contain 
                                    :src="`/api/images/${mozaik.Name}`" 
                                />
                            </v-col>
                            <v-col 
                                cols="6"
                                align-self="center"
                            >
                                <svg 
                                    viewPort="0 0 48 48" 
                                    style="width:480px;height:480px"
                                >                        
                                    <g 
                                        v-for="(col, colIdx) in mozaik.Image" 
                                        :key="`mozaic-col-${colIdx}`"
                                    >
                                        <rect
                                            v-for="(row, rowIdx) in col"
                                            :key="`mozaic-col-${colIdx}-${rowIdx}`"
                                            :x="colIdx * 10"
                                            :y="rowIdx * 10"
                                            :id="`mozaic-col-${colIdx}-${rowIdx}`"
                                            :style="`fill: rgb(${row.r}, ${row.g}, ${row.b});`"
                                            :ref="`orig-${colIdx}-${rowIdx}`"
                                            width="10"
                                            height="10"
                                            class="pixel"
                                            v-on:pointermove="hover(rowIdx, colIdx)"
                                            @click="tileClicked(rowIdx, colIdx)"
                                        />
                                    </g>
                                </svg>
                            </v-col>
                        </v-row>                        
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>
    </v-sheet>
</template>
<script>
  export default {
    name: 'MozaikView',
    data: () => ({
        loading: true,
        searchColor: "",
        hoverColor: "",
        hoverColorName: "",
        selectedColor: "",
        changes: [],
        lastChangeIndex: 0,
    }),
    computed: {
        mozaik: {
            get() {
                return this.$store.state.mozaik
            },
            set(m) {
                this.$store.commit("mozaik", m)
            }
        },
        colors() {            
            return Object.fromEntries(
                Object.entries(this.$store.state.colors).filter(color => color[1].name.toLowerCase().includes(this.searchColor.toLowerCase()))
            )
        }
    },
    async mounted() {
        await this.$store.dispatch("fetchMozaik", this.$route.params.name)
        this.loading = false
    },
    methods: {
        hover(row, col) {            
            const pixel = this.mozaik.Image[col][row]
            const hexName = `#${pixel.r.toString(16).padStart(2,0)}${pixel.g.toString(16).padStart(2,0)}${pixel.b.toString(16).padStart(2,0)}`
            this.hoverColor = hexName
            this.hoverColorName = this.colors[hexName].name
        },
        tileClicked(row, col) {
            if (!this.selectedColor) {
                return
            }            
            const mozaik = {...this.mozaik}
            this.changes = this.changes.slice(0, this.lastChangeIndex)
            this.changes.push({
                col: col,
                row: row,
                from: mozaik.Image[col][row],
                to: Object.entries(this.colors)[this.selectedColor][1]
            })
            this.lastChangeIndex = this.changes.length
            mozaik.Image[col][row] = Object.entries(this.colors)[this.selectedColor][1]
            this.mozaik = mozaik
        },
        undo() {
            this.lastChangeIndex = this.lastChangeIndex - 1
            const change = this.changes[this.lastChangeIndex]
            const mozaik = {...this.mozaik}
            mozaik.Image[change.col][change.row] = change.from
            this.mozaik = mozaik
        },
        redo() {            
            const change = this.changes[this.lastChangeIndex]
            this.lastChangeIndex = this.lastChangeIndex + 1
            const mozaik = {...this.mozaik}
            mozaik.Image[change.col][change.row] = change.to
            this.mozaik = mozaik
        },
        async save() {
            this.loading = true
            await this.$store.dispatch("save")
            this.loading = false
            this.lastChangeIndex = 0
            this.changes = []
        },
        async remove() {
            this.loading = true
            await this.$store.dispatch("remove")
            this.loading = false
            this.$router.push({name:"mozaik-list"})
        }
    }
  }
</script>