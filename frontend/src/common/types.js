const _proxyTypeMap = [
    {label: 'http', 'value': 1},
    {label: 'https', 'value': 2},
    {label: 'tcp', 'value': 3},
    {label: 'udp', 'value': 4},
]

const _proxyStatusMap = [
    {label: '启用', 'value': 1},
    {label: '禁用', 'value': 2},
]

const _labelFilter = (value, mapValue) => {
    let label = ''
    mapValue.filter(item => {
        if (item.value === value) {
            label = item.label
        }
    })
    return label
}

const handleProxyTypeName = (proxyType) => {
    return _labelFilter(proxyType, _proxyTypeMap)
}

const handleProxyStatusName = (proxyStatus) => {
    return _labelFilter(proxyStatus, _proxyStatusMap)
}

export {
    handleProxyTypeName,
    handleProxyStatusName,
}
