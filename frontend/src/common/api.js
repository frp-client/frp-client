import request from "./request.js";

const login = (data) => {
    return request.post('/api/user/login', data)
}

const getProxies = (params = null) => {
    return request.get('/api/frpc/proxies', {params: params})
}

const createProxy = (data) => {
    return request.post('/api/frpc/proxy', data)
}

export default {
    login,
    getProxies,
    createProxy,
}
