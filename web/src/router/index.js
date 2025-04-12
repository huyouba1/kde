import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/cluster',
    name: 'ClusterManagement',
    component: () => import('../views/cluster/ClusterManagement.vue')
  },
  {
    path: '/deploy',
    name: 'ApplicationDeployment',
    component: () => import('../views/deploy/ApplicationDeployment.vue')
  },
  {
    path: '/plugins',
    name: 'PluginManagement',
    component: () => import('../views/plugins/PluginManagement.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router