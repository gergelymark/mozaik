import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    imageToCrop: "",
    imageToCropName: "",
    colors: {},
    mozaiks: [],
    createRequested: "",
    mozaik: {},
  },
  getters: {
  },
  mutations: {
    colors(state, colors) {
      state.colors = colors
    },
    imageToCrop(state, image) {
      state.imageToCrop = image
    },
    imageToCropName(state, name) {
      state.imageToCropName = name
    },
    mozaiks(state, mozaiks) {
      state.mozaiks = mozaiks
    },
    createRequested(state, time) {
      state.createRequested = time
    },
    mozaik(state, mozaik) {
      state.mozaik = mozaik
    }
  },
  actions: {
    async getColors(context) {
      const resp = await axios.get("/api/colors/")
      context.commit("colors", resp.data)
    },
    async fetchAllMozaik(context) {
      const resp = await axios.get("/api/mozaik/")
      context.commit("mozaiks", resp.data)
    },
    async createMozaik(_, mozaikData) {
      await axios.post("/api/mozaik/", mozaikData)      
    },
    async fetchMozaik(context, id) {
      const resp = await axios.get(`/api/mozaik/${id}`)
      context.commit("mozaik", resp.data)
    },
    async save(context) {
      await axios.put("/api/mozaik/", context.state.mozaik)
    },
    async remove({state, dispatch}) {
      await axios.delete(`/api/mozaik/${state.mozaik.Name}`)
      await dispatch("fetchAllMozaik")
    }
  },
  modules: {
  }
})
