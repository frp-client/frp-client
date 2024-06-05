<template>
  <div class="sidebar flex-align-items-center text-align-center">

    <div class="flex-column">
      <v-avatar :size="80" class="avatar">
        <v-img src="https://cdn.vuetifyjs.com/images/parallax/material.jpg"></v-img>
      </v-avatar>

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
import {defineComponent, onBeforeMount, ref} from "vue";
import {Icon} from '@vicons/utils'
import {CloseFilled, HorizontalRuleFilled, MenuFilled, MinusFilled, SettingsFilled} from "@vicons/material";

const onBeforeMountHandler = () => {

}

const menus = ref([
  {
    id: 1,
    name: '首页',
    to: '/',
    activeClass: 'active',
  }, {
    id: 2,
    name: '转发',
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

export default defineComponent({
  components: {Icon, CloseFilled, HorizontalRuleFilled, MinusFilled, SettingsFilled, MenuFilled},
  setup() {
    onBeforeMount(onBeforeMountHandler)
    return {
      menus,
      onClickMenu,
    }
  }
})

</script>

<style scoped>

.sidebar {
  display: flex;
  flex-direction: column;
  width: 160px;

  .avatar {
    margin-top: 40px;
  }

  .links {
    margin-top: 40px;
    width: 100%;

    a {
      display: block;
      text-decoration: none;
      color: rgb(31, 34, 37);
      padding: 15px 0 15px 4px;
    }

    a:hover, a.active {
      background-color: ghostwhite;
      border-left: 4px solid rgba(0, 188, 213, 0.68);
      padding: 15px 0 15px 0;
    }
  }
}


</style>
