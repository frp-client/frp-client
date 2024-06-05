import 'material-design-icons-iconfont/dist/material-design-icons.css' // Ensure your project is capable of handling css files

// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import { aliases, md } from 'vuetify/iconsets/md'

export default createVuetify({
    components,
    directives,

    // 字体参考：https://vuetifyjs.com/zh-Hans/features/icon-fonts/
    // https://fonts.google.com/icons
    icons: {
        defaultSet: 'md',
        aliases,
        sets: {
            md,
        },
    },

})
