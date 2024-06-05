import {createRouter, createWebHashHistory} from 'vue-router'

const routes = [
    {
        path: '/',
        name: 'home',
        component: () => import('../views/home.vue')
    },
    {
        path: '/about',
        name: 'about',
        component: () => import('../views/about.vue')
    },
    {
        path: '/setting',
        name: 'setting',
        component: () => import('../views/setting.vue')
    },
    {
        path: '/proxy',
        name: 'proxy',
        component: () => import('../views/proxy.vue')
    },
]

export default createRouter({
    history: createWebHashHistory(),
    routes,
})