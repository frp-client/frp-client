import request from "./request.js";

const login = (data) => {
    return request.post('/api/user/login',data)
}

export default {
    login,
}
