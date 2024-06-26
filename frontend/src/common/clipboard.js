const androidWrite = (text) => {
    let result = {ok: false, msg: ''}
    let textarea = document.createElement("textarea");
    textarea.style.position = 'fixed';
    textarea.style.border = 'none';
    textarea.style.outline = 'none';
    textarea.style.background = 'transparent';
    textarea.style.color = 'transparent';
    textarea.value = text;
    document.body.appendChild(textarea)
    textarea.select();
    try {
        if (document.execCommand("copy")) {
            result = {ok: true, value: text}
        } else {
            result = {ok: false, value: 'copy failed'}
        }
    } catch (err) {
        result = {ok: false, value: err}
    }
    document.body.removeChild(textarea);
    return result
}

const otherWrite = (text) => {
    let result = {ok: false, msg: ''}

    window.getSelection().removeAllRanges();//这段代码必须放在前面否则无效
    let input = document.createElement("input");
    document.body.appendChild(input);
    input.setAttribute("value", text);
    input.select();
    let range = document.createRange();
    // 选中需要复制的节点
    range.selectNode(input);
    // 执行选中元素
    window.getSelection().addRange(range);
    input.select();
    input.setSelectionRange(0, input.value.length);//适配高版本ios
    // 执行 copy 操作
    if (document.execCommand('copy')) {
        result = {ok: true, value: text}
    } else {
        result = {ok: false, value: 'copy failed'}
    }

    // 移除选中的元素
    window.getSelection().removeAllRanges();
    document.body.removeChild(input);

    return result;
}

// !!! 必须是点击事件触发才会生效，不然不会返回成功
export default {
    read(isAndroid) {
    },
    write(text, isAndroid) {
        return isAndroid ? androidWrite(text) : otherWrite(text)
    },
}
