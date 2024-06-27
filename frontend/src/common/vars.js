let _session = {}
let _appConfig = {}

const getSession = () => {
    return _session
}

const setSession = (session) => {
    _session = session
}

const getAppConfig = () => {
    return _appConfig
}

const setAppConfig = (appConfig) => {
    _appConfig = appConfig
}

export {
    getSession,
    setSession,
    getAppConfig,
    setAppConfig,
}
