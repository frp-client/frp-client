import request from "./request.js";

const login = (data) => {
    return request.post('/api/user/login', data)
}

const getProxies = (params = null) => {
    return request.get('/api/frpc/proxies', {params: params})
}

const getProxy = (id) => {
    return request.get(`/api/frpc/proxy/${id}`)
}

const createProxy = (data) => {
    return request.post('/api/frpc/proxy', data)
}

const updateProxy = (id, data) => {
    return request.put(`/api/frpc/proxy/${id}`, data)
}

const deleteProxy = (id) => {
    return request.delete(`/api/frpc/proxy/${id}`)
}

export default {
    login,
    getProxies,
    getProxy,
    createProxy,
    updateProxy,
    deleteProxy,
}
