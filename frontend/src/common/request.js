import axios from "axios";

let _config = {baseURL: '', jwtToken: '', clientId: '',}

const request = axios.create({
    baseURL: '',
    timeout: 50000,
})

request.interceptors.request.use(value => {
    if (!value.headers['X_CLIENT_ID']) {
        value.headers['X_CLIENT_ID'] = _config.clientId
    }
    if (!value.headers['Authorization'] && _config.jwtToken) {
        value.headers['Authorization'] = `Bearer ${_config.jwtToken}`
    }

    // console.log('[request.data]', value)
    return value
}, error => {
    // console.log('[request.error]', error)
    return Promise.reject(error)
})


request.interceptors.response.use(resp => {
    // console.log('[response.data]', resp)
    if (resp.data.code === 200) {
        return resp
    }
    // 直接可以显示错误信息
    resp.data.toString = () => {
        return resp.data.msg
    }
    return Promise.reject(resp.data)
}, error => {
    // console.log('[response.error]', error)
    return Promise.reject(error)
})

request.config = (config = _config) => {
    // console.log('[request.config]', config)
    _config = Object.assign(_config, config)
    request.defaults.baseURL = _config.baseURL
}

export default request
