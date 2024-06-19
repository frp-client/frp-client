<template>
  <div>
    <Head style="--wails-draggable:drag"></Head>
    <div style="--wails-draggable:no-drag;">
      <div class="container">
        <div class="main">
          <Sidebar style="width: 160px;"/>
          <router-view class="plr20 ptb20" style="flex: 1; overflow-y: auto;"/>
        </div>
        <div class="footer flex-items-center flex-items-between plr20">
          <div>
            <router-link class="alink" @click="openUrl('https://github.com/frp-client/frp-client')" to="">
              frp-client
            </router-link>
            <span class="mlr2">v1.0.0</span>

            <span v-if="clientId" class="mlr6">
              客户端ID: {{ clientId.substring(0, 12) }}
            </span>

          </div>
          <div>{{ clock }}</div>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import {defineComponent, onBeforeMount, onMounted, ref} from "vue";
import Head from './components/Head.vue'
import {useRoute, useRouter} from "vue-router";
import Sidebar from "./components/Sidebar.vue";
import {openUrl} from "./common/helper.js";
import {ClientId} from "../wailsjs/go/main/App.js";
import {EventsOn} from "../wailsjs/runtime/runtime.js";
import request from "./common/request.js";

let route = null
let router = null
let clock = ref(null)
let clientId = ref(null)

const onBeforeMountHandler = () => {
  EventsOn('onStartUpEvent', (data) => {
    console.log('[onStartUpEvent]', data)
    request.setBaseURL(data.baseURL)
    request.setClientId(data.clientId)
  })
  timeClock()
}

const onMountedHandler = () => {
  ClientId().then(resp => {
    console.log('[onMountedHandler]', resp)
    clientId.value = resp
  }).catch(err => {
    console.log('[onMountedHandler]', err)
  })
  // alert('[ClientId]'+a)
}

const timeClock = () => {
  clock.value = (new Date()).toLocaleString()
  setInterval(() => {
    clock.value = (new Date()).toLocaleString()
  }, 2000)
}

export default defineComponent({
  methods: {openUrl},
  components: {Sidebar, Head},
  setup() {
    route = useRoute()
    router = useRouter()
    onBeforeMount(onBeforeMountHandler)
    onMounted(onMountedHandler)
    return {
      clock,
      clientId,
    }
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
    overflow: hidden;
  }

  .footer {
    height: 38px;
    font-size: 12px;
    background-color: ghostwhite;
  }
}

</style>
