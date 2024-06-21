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
    {
        path: '/proxy-edit',
        name: 'proxy-edit',
        component: () => import('../views/proxy-edit.vue')
    },
    {
        path: '/login',
        name: 'login',
        component: () => import('../views/login.vue')
    },
    {
        path: '/debug',
        name: 'debug',
        component: () => import('../views/debug.vue')
    },
]

export default createRouter({
    history: createWebHashHistory(),
    routes,
})
