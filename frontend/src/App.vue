<template>
  <div>
    <Head style="--wails-draggable:drag"></Head>
    <div style="--wails-draggable:no-drag;">
      <div class="container">
        <div class="main">
          <Sidebar/>
          <router-view class="plr20 ptb10" />
        </div>
        <div class="footer flex-items-center plr20">
          <span class="plr2">FOOTER</span>
          <a class="plr2" target="_blank" href="https://github.com">github</a>
          <a class="plr2" target="_blank" href="https://github.com">github</a>
          <a class="plr2" target="_blank" href="https://github.com">github</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {defineComponent, onBeforeMount} from "vue";
import Head from './components/Head.vue'
import {useRoute, useRouter} from "vue-router";
import Sidebar from "./components/Sidebar.vue";

let route = null
let router = null

const onBeforeMountHandler = () => {

  console.log('[onBeforeMountHandler1]', route)
  console.log('[onBeforeMountHandler2]', router)
  console.log('[onBeforeMountHandler3]', router.getRoutes())
}

export default defineComponent({
  components: {Sidebar, Head},
  setup() {
    route = useRoute()
    router = useRouter()
    onBeforeMount(onBeforeMountHandler)
    return {}
  },
})

</script>

<style scoped>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}

.container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 42px);

  .main {
    display: flex;
    flex-direction: row;
    flex: 1;
  }
  .footer{
    height: 38px;
    font-size: 12px;
    background-color: ghostwhite;
  }
}

</style>
