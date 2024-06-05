<template>
  <div class="flex-column flex-items-center flex-justify-content-center about">
    <p>
      <text class="plr2">frp客户端</text>
      <text class="plr2">v1.0.0</text>
    </p>

    <router-link class="alink" @click="openUrl('https://github.com/')" to="">
      github
    </router-link>

    <p class="p-gray">通过将内网端口暴露到公网，实现从互联网访问内网资源</p>
    <p class="p-gray">请勿使用该软件从事违法用途</p>

  </div>
</template>

<script>
import {defineComponent} from "vue";
import {openUrl} from "../common/helper.js";

export default defineComponent({
  setup() {
    return {
      openUrl,
    }
  }
})

</script>

<style scoped>
.about {
  line-height: 250%;
}
</style>