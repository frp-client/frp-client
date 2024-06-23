const handleProxyDomain = (data) => {
    let domain = ''
    const _domain = data.domain
    if (data.proxy_extra.subdomain) {
        switch (data.proxy_type) {
            case 1://http
                domain = `http://${data.proxy_extra.subdomain}.${_domain}`
                break
            case 2://https
                domain = `https://${data.proxy_extra.subdomain}.${_domain}`
                break
            case 3://tcp
                domain = `tcp://${data.proxy_extra.subdomain}.${_domain}:${data.proxy_remote_port}`
                break
            case 4://udp
                domain = `udp://${data.proxy_extra.subdomain}.${_domain}:${data.proxy_remote_port}`
                break
        }
    }
    return domain
}

export {
    handleProxyDomain,
}
