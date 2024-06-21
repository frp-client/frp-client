import {OpenUrl} from "../../wailsjs/go/main/App.js";

const openUrl = (url) => {
    OpenUrl(url)
}

const timeFormat = (timestamp) => {
    if (+timestamp > 0) {
        return (new Date(timestamp * 1000)).toLocaleString()
    }
    return ''
}

export {
    openUrl,
    timeFormat,
}
