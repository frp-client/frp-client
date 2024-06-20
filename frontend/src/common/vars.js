let _session = {}

const getSession = () => {
    return _session
}

const setSession = (session) => {
    _session = session
}

export {
    getSession,
    setSession,
}
