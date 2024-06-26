<template>
  <div class="flex-column flex-items-center flex-justify-content-center">
    <p class="p-gray"></p>

    <router-link to="/login">登录</router-link>
    <v-btn @click="onClickReloadApp" class="mt-2 mb-2">WindowReloadApp</v-btn>

  </div>
</template>

<script>
import {defineComponent} from "vue";
import {WindowReloadApp} from "../../wailsjs/runtime/runtime.js";

const onClickReloadApp = () => {
  console.log('[onClickReloadApp]')
  WindowReloadApp()
}

export default defineComponent({
  setup() {
    return {
      onClickReloadApp,
    }
  }
})
</script>

<style scoped lang="scss">


</style>
