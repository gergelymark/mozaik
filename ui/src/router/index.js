import Vue from 'vue'
import VueRouter from 'vue-router'
import MozaikList from '../views/MozaikList.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'mozaik-list',
    component: MozaikList
  },
  {
    path: '/create',
    name: 'create',
    component: () => import(/* webpackChunkName: "about" */ '../views/CreateMozaik.vue')
  },
  {
    path: '/mozaik/:name',
    name: 'mozaik',
    component: () => import(/* webpackChunkName: "about" */ '../views/Mozaik.vue')
  },
  {
    path: '/parts',
    name: 'parts',
    component: () => import(/* webpackChunkName: "about" */ '../views/Parts.vue')
  },

]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
