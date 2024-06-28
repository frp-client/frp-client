<template>
  <div class="">

    <div class="ml-2 mr-2">
      <v-chip class="pointer" color="primary" label @click="reloadAppLogs">
        <v-icon icon="view_list" start></v-icon>
        日志列表
      </v-chip>
    </div>

    <v-list ref="refLogList" @scrollend="alert(22)" @scroll="alert(5656)">
      <v-list-item v-for="(item,index) in respLogs.logs" :key="index" @scrollend="alert(1)">
        <b>{{ index }}</b>、{{ item }}
      </v-list-item>
    </v-list>

    <div class="ptb20 flex-row flex-justify-content-center">
      <v-chip>每隔5秒钟将加载更多日志</v-chip>
    </div>

    <MySnackbar ref="refMySnackbar"></MySnackbar>
    <MyLoading ref="refMyLoading"></MyLoading>


  </div>
</template>

<script>
import {defineComponent, onMounted, ref} from "vue";
import {openUrl} from "../common/helper.js";
import {AppLogs} from "../../wailsjs/go/main/App.js";
import MySnackbar from "../components/MySnackbar.vue";
import MyLoading from "../components/MyLoading.vue";

const refMyLoading = ref(MyLoading)
const refMySnackbar = ref(MySnackbar)

const timer = ref(null)
const refLogList = ref(null)
const respLogs = ref({
  start_line: 0,
  last_line: 0,
  logs: [],
})

const getAppLogs = () => {
  refMyLoading.value.show()
  AppLogs(respLogs.value.start_line).then(resp => {
    console.log('[AppLogs.resp]', resp)
    if (!resp.logs || resp.logs.length === 0) {
      console.log('[clearInterval]', timer.value)
      return
    }

    respLogs.value.start_line = resp.last_line
    respLogs.value.last_line = resp.last_line
    respLogs.value.logs = respLogs.value.logs.concat(...resp.logs)

    // console.log('[refLogList]', refLogList.value)
    // console.log('[refLogList2]', refLogList.value.height)
    // console.log('[refLogList3]', refLogList.value.scrollHeight)

  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
  })
}

const reloadAppLogs = () => {
  respLogs.value.start_line = 0
  respLogs.value.last_line = 0
  respLogs.value.logs = []
  getAppLogs()
}

const onMountedHandler = () => {
  getAppLogs()

  timer.value = setInterval(() => {
    getAppLogs()
  }, 3000)

}

export default defineComponent({
  components: {MyLoading, MySnackbar},
  setup() {
    onMounted(onMountedHandler)
    return {
      openUrl,
      reloadAppLogs,
      respLogs,
      refMyLoading,
      refMySnackbar,
      refLogList,

    }
  }
})

</script>

<style scoped lang="scss">
.about {
  line-height: 250%;
}
</style>
