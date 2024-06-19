import axios from "axios";

let _client_id = ''

const request = axios.create({
    baseURL: '',
    timeout: 50000,
})

request.interceptors.request.use(value => {
    if (!value.headers['X_CLIENT_ID']) {
        value.headers['X_CLIENT_ID'] = _client_id
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
    return Promise.reject(resp.data)
}, error => {
    // console.log('[response.error]', error)
    return Promise.reject(error)
})

request.setBaseURL = (_baseUrl) => {
    request.defaults.baseURL = _baseUrl
}
request.setClientId = (clientId) => {
    _client_id = clientId
}

export default request
