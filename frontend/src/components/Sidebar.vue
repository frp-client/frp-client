<template>
  <div class="sidebar flex-align-items-center text-align-center">

    <div class="flex-column flex-align-items-center ">
      <v-avatar :size="80" class="avatar pointer" @click="onClickLogin">
        <v-img src="https://cdn.vuetifyjs.com/images/parallax/material.jpg"></v-img>
      </v-avatar>
      <div class="username mt-3 pointer overflow-hidden" @click="onClickLogin" v-if="username">
        {{ username.substring(0, 12) }}
      </div>
    </div>
    <div class="flex-column links">
      <router-link
          v-for="(item,idx) in menus" :key="idx"
          :to="item.to"
          :class="[item.activeClass]"
          @click="onClickMenu(item)">
        {{ item.name }}
      </router-link>

    </div>

  </div>
</template>

<script>
import {defineComponent, onBeforeMount, onMounted, ref} from "vue";
import {Icon} from '@vicons/utils'
import {CloseFilled, HorizontalRuleFilled, MenuFilled, MinusFilled, SettingsFilled} from "@vicons/material";
import router from "../router/index.js";
import {EventsOn} from "../../wailsjs/runtime/runtime.js";

const username = ref(null)

const onBeforeMountHandler = () => {
  EventsOn('onStartUpEvent', (data) => {
    username.value = data.username
  })
}

const onMountedHandler = () => {
}

const menus = ref([
  {
    id: 1,
    name: '首页',
    to: '/',
    activeClass: 'active',
  }, {
    id: 2,
    name: '规则',
    to: '/proxy',
    activeClass: '',
  }, {
    id: 4,
    name: '设置',
    to: '/setting',
    activeClass: '',
  }, {
    id: 5,
    name: '关于',
    to: '/about',
    activeClass: '',
  }
])

const onClickMenu = (data) => {
  console.log('[onClickMenu]', data)
  menus.value.map(item => {
    if (item.id === data.id) {
      item.activeClass = 'active'
    } else {
      item.activeClass = ''
    }
    return item
  })
}

const onClickLogin = () => {
  router.push('/login')
}

export default defineComponent({
  components: {Icon, CloseFilled, HorizontalRuleFilled, MinusFilled, SettingsFilled, MenuFilled},
  setup() {
    onBeforeMount(onBeforeMountHandler)
    onMounted(onMountedHandler)
    return {
      menus,
      onClickMenu,
      onClickLogin,
      username,
    }
  }
})

</script>

<style scoped lang="scss">

.sidebar {
  display: flex;
  flex-direction: column;


  .avatar {
    margin-top: 40px;
  }

  .username {
    width: 120px;
  }

  .links {
    margin-top: 30px;
    width: 100%;

    a {
      display: block;
      text-decoration: none;
      color: rgb(31, 34, 37);
      padding: 15px 0 15px 4px;
    }

    a:hover {
      background-color: #f7fbff !important;
      border-left: 4px solid rgba(0, 188, 213, 0.68);
      padding: 15px 0 15px 0;
    }

    a.active {
      background-color: ghostwhite;
      border-left: 4px solid rgba(0, 188, 213, 0.68);
      padding: 15px 0 15px 0;
    }
  }
}


</style>
