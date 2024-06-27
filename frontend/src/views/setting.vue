<template>
  <div>

    <div class="ml-2 mr-2">
      <div class="ml-2 mr-2">
        <v-alert text="配置修改后需要重启程序以生效" border="start" border-color="info"></v-alert>
      </div>
    </div>

    <v-form ref="refForm" class="form ml-2 mr-2">
      <v-container>
        <div class="ptb10"></div>
        <v-row>
          <v-col cols="12" md="12">
            <v-text-field
                label="API服务器"
                v-model="formData.apiServer.value"
                :rules="formData.apiServer.rule"
                variant="underlined"
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="2" md="2">
            <v-select
                label="本地WEB服务"
                v-model="formData.localWebServer.select"
                :items="formData.localWebServer.items"
                :rules="formData.localWebServer.rule"
                item-title="label"
                item-value="value"
                variant="underlined"
            ></v-select>
          </v-col>
          <v-col cols="2" md="2" v-if="formData.localWebServer.select">
            <v-text-field
                label="端口"
                v-model="formData.localWebServerPort.value"
                :rules="formData.localWebServerPort.rule"
                type="number"
                variant="underlined"
            ></v-text-field>
          </v-col>
          <v-col cols="5" md="5" class="flex-items-center" v-if="formData.localWebServer.select">
            <v-text-field
                label="本地WEB目录"
                placeholder="请复制粘贴本地WEB根目录地址"
                v-model="formData.localWebServerPath.value"
                :rules="formData.localWebServerPath.rule"
                variant="underlined"
            ></v-text-field>
            <div class="ml-2"></div>
            <v-icon
                icon="content_copy"
                class="pointer"
                @click="copyToClipboard(formData.localWebServerPath.value)">
            </v-icon>
          </v-col>
          <v-col cols="3" md="3" class="flex-items-center" v-if="formData.localWebServer.select">
            <text class="dashed-x pointer" @click="copyToClipboard(handleLocalUrl(formData.localWebServerPort.value))">
              {{ handleLocalUrl(formData.localWebServerPort.value) }}
            </text>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="3" md="3">
            <v-select
                label="日志"
                v-model="formData.log.select"
                :items="formData.log.items"
                :rules="formData.log.rule"
                item-title="label"
                item-value="value"
                variant="underlined"
            ></v-select>
          </v-col>
          <v-col cols="9" md="9" class="flex-items-center" v-if="formData.log.select">
            <text class="dashed-x pointer overflow-hidden" @click="copyToClipboard(formData.logPath.value)">
              {{ formData.logPath.value }}
            </text>
          </v-col>
        </v-row>


        <v-row>
          <v-col cols="12" md="12" class="d-flex justify-end">
            <v-btn color="rgb(24, 103, 192)" @click="onClickSubmit">保存</v-btn>
          </v-col>
        </v-row>

      </v-container>
    </v-form>

    <v-dialog v-model="showTipsModal" width="500">
      <v-card text="确定关闭程序？">
        <template v-slot:title>
          <div class="flex-items-center">
            <v-icon icon="info"/>&nbsp;提示
          </div>
        </template>
        <template v-slot:actions>
          <v-btn text="取消" @click="showTipsModal = false"></v-btn>
          <v-btn text="确定" @click=""></v-btn>
        </template>
      </v-card>
    </v-dialog>

    <MySnackbar ref="refMySnackbar"></MySnackbar>
    <MyLoading ref="refMyLoading"></MyLoading>

  </div>
</template>

<script>
import {defineComponent, onMounted, ref} from "vue";
import clipboard from "../common/clipboard.js";
import MySnackbar from "../components/MySnackbar.vue";
import {getAppConfig} from "../common/vars.js";
import {AppConfigUpdate} from "../../wailsjs/go/main/App.js";
import MyLoading from "../components/MyLoading.vue";

const refForm = ref(false)
const showTipsModal = ref(false)
const refMySnackbar = ref(MySnackbar)
const refMyLoading = ref(MyLoading)

const formData = ref({
  apiServer: {
    value: '',
    rule: [
      value => {
        return true;
      },
    ],
  },
  localWebServer: {
    select: {label: '关闭', value: 0},
    items: [
      {label: '开启', value: 1},
      {label: '关闭', value: 0},
    ],
    rule: [],
  },
  localWebServerPort: {
    value: null,
    rule: [],
  },
  localWebServerPath: {
    value: null,
    rule: [],
  },
  log: {
    select: {label: '关闭日志', value: 0},
    items: [
      {label: '开启日志', value: 1},
      {label: '关闭日志', value: 0},
    ],
    rule: [],
  },
  logPath: {
    value: null,
    rule: [],
  },
})

const onClickSubmit = async () => {
  console.log('[onClickSubmit]', formData.value)
  const {valid} = await refForm.value.validate()

  refMyLoading.value.show()
  AppConfigUpdate({
    api_server: formData.value.apiServer.value,
    local_server: !!formData.value.localWebServer.select,
    local_server_port: formData.value.localWebServerPort.value,
    local_server_path: formData.value.localWebServerPath.value,
    log: !!formData.value.log.select,
  }).then(resp => {
    console.log('[AppConfigUpdate] resp ', resp)
  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
  })

  // showTipsModal.value=true

}

const copyToClipboard = (url) => {
  clipboard.write(url, false)
  refMySnackbar.value.show('已复制：' + url, {timeout: 1000})
}

const handleLocalUrl = (port) => {
  if (port) {
    return `http://127.0.0.1:${port}`
  }
  return ''
}

const onMountedHandler = () => {
  const appConfig = getAppConfig()
  formData.value.apiServer.value = appConfig.api_server
  formData.value.localWebServer.select = +appConfig.local_server
  formData.value.localWebServerPort.value = +appConfig.local_server_port
  formData.value.localWebServerPath.value = appConfig.local_server_path
  formData.value.log.select = +appConfig.log
  formData.value.logPath.value = appConfig.log_path
}

export default defineComponent({
  components: {MyLoading, MySnackbar},
  setup() {

    onMounted(onMountedHandler)

    return {
      refForm,
      formData,
      onClickSubmit,
      showTipsModal,
      refMySnackbar,
      refMyLoading,
      copyToClipboard,
      handleLocalUrl,
    }
  }
})

</script>

<style scoped lang="scss">
.form {
}
</style>
