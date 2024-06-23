const handleProxyDomain = (data) => {
    let domain = ''
    const _domain = data.domain
    if (data.proxy_extra.subdomain) {
        switch (data.proxy_type) {
            case 1:
                domain = `http://${data.proxy_extra.subdomain}.${_domain}:${data.proxy_remote_port}`
                break
            case 2:
                domain = `https://${data.proxy_extra.subdomain}.${_domain}:${data.proxy_remote_port}`
                break
            case 3:
                domain = `tcp://${data.proxy_extra.subdomain}.${_domain}:${data.proxy_remote_port}`
                break
            case 4:
                domain = `udp://${data.proxy_extra.subdomain}.${_domain}:${data.proxy_remote_port}`
                break
        }
    }
    return domain
}

export {
    handleProxyDomain,
}
