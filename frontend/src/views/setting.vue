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
          <v-col cols="12" md="12"><b>Api服务器</b></v-col>
        </v-row>
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
          <v-col cols="12" md="12"><b>本地Web服务</b></v-col>
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
          <v-col cols="12" md="12"><b>本地Tcp服务</b></v-col>
        </v-row>
        <v-row>
          <v-col cols="2" md="2">
            <v-select
                label="本地TCP服务"
                v-model="formData.localTcpServer.select"
                :items="formData.localTcpServer.items"
                :rules="formData.localTcpServer.rule"
                item-title="label"
                item-value="value"
                variant="underlined"
            ></v-select>
          </v-col>
          <v-col cols="2" md="2" v-if="formData.localTcpServer.select">
            <v-text-field
                label="端口"
                v-model="formData.localTcpServerPort.value"
                :rules="formData.localTcpServerPort.rule"
                type="number"
                variant="underlined"
            ></v-text-field>
          </v-col>
          <v-col cols="5" md="5" class="flex-items-center" v-if="formData.localTcpServer.select">
            <v-text-field
                label="TCP回复数据"
                placeholder=""
                v-model="formData.localTcpServerResponse.value"
                :rules="formData.localTcpServerResponse.rule"
                variant="underlined"
            ></v-text-field>
            <div class="ml-2"></div>
          </v-col>
          <v-col cols="3" md="3" class="flex-items-center" v-if="formData.localTcpServer.select">
            <text class="dashed-x pointer" @click="copyToClipboard(handleLocalTcp(formData.localTcpServerPort.value))">
              {{ handleLocalTcp(formData.localTcpServerPort.value) }}
            </text>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="12"><b>本地Udp服务</b></v-col>
        </v-row>
        <v-row>
          <v-col cols="2" md="2">
            <v-select
                label="本地UDP服务"
                v-model="formData.localUdpServer.select"
                :items="formData.localUdpServer.items"
                :rules="formData.localUdpServer.rule"
                item-title="label"
                item-value="value"
                variant="underlined"
            ></v-select>
          </v-col>
          <v-col cols="2" md="2" v-if="formData.localUdpServer.select">
            <v-text-field
                label="端口"
                v-model="formData.localUdpServerPort.value"
                :rules="formData.localUdpServerPort.rule"
                type="number"
                variant="underlined"
            ></v-text-field>
          </v-col>
          <v-col cols="5" md="5" class="flex-items-center" v-if="formData.localUdpServer.select">
            <v-text-field
                label="UDP回复数据"
                placeholder=""
                v-model="formData.localUdpServerResponse.value"
                :rules="formData.localUdpServerResponse.rule"
                variant="underlined"
            ></v-text-field>
            <div class="ml-2"></div>
          </v-col>
          <v-col cols="3" md="3" class="flex-items-center" v-if="formData.localUdpServer.select">
            <text class="dashed-x pointer" @click="copyToClipboard(handleLocalUdp(formData.localUdpServerPort.value))">
              {{ handleLocalUdp(formData.localUdpServerPort.value) }}
            </text>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="12"><b>本地Shodowsocks服务</b></v-col>
        </v-row>
        <v-row>
          <v-col cols="3" md="3">
            <v-select
                label="本地shadowsocks服务"
                v-model="formData.localSsServer.select"
                :items="formData.localSsServer.items"
                :rules="formData.localSsServer.rule"
                item-title="label"
                item-value="value"
                variant="underlined"
            ></v-select>
          </v-col>
          <v-col cols="4" md="4" class="flex-items-center" v-if="formData.localSsServer.select">
            <v-select
                label="加密方式"
                v-model="formData.localSsCipher.select"
                :items="formData.localSsCipher.items"
                :rules="formData.localSsCipher.rule"
                item-title="label"
                item-value="value"
                variant="underlined"
            ></v-select>
          </v-col>
          <v-col cols="2" md="2" v-if="formData.localSsServer.select">
            <v-text-field
                label="端口"
                v-model="formData.localSsPort.value"
                :rules="formData.localSsPort.rule"
                type="number"
                variant="underlined"
            ></v-text-field>
          </v-col>
          <v-col cols="3" md="3" class="flex-items-center" v-if="formData.localSsServer.select">
            <v-text-field
                label="密码"
                placeholder=""
                v-model="formData.localSsPassword.value"
                :rules="formData.localSsPassword.rule"
                variant="underlined"
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="12"><b>日志</b></v-col>
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
            <v-icon icon="visibility" class="mlr6 pointer" @click="$router.push('/log')"></v-icon>
          </v-col>
        </v-row>


        <v-row>
          <v-col cols="12" md="12" class="d-flex justify-end">
            <v-btn class="mlr20" color="orange" @click="onClickReset">重置</v-btn>
            <v-btn class="" color="rgb(24, 103, 192)" @click="onClickSubmit">保存</v-btn>
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
import {setAppConfig} from "../common/vars.js";
import {AppConfig, AppConfigReset, AppConfigUpdate} from "../../wailsjs/go/main/App.js";
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
  localTcpServer: {
    select: {label: '关闭', value: 0},
    items: [
      {label: '开启', value: 1},
      {label: '关闭', value: 0},
    ],
    rule: [],
  },
  localTcpServerPort: {
    value: null,
    rule: [],
  },
  localTcpServerResponse: {
    value: null,
    rule: [],
  },
  localUdpServer: {
    select: {label: '关闭', value: 0},
    items: [
      {label: '开启', value: 1},
      {label: '关闭', value: 0},
    ],
    rule: [],
  },
  localUdpServerPort: {
    value: null,
    rule: [],
  },
  localUdpServerResponse: {
    value: null,
    rule: [],
  },

  localSsServer: {
    select: 1,
    items: [
      {label: '开启', value: 1},
      {label: '关闭', value: 0},
    ],
    rule: [],
  },
  localSsPort: {
    value: null,
    rule: [],
  },
  localSsCipher: {
    // https://github.com/softwaredownload/openwrt-fanqiang/blob/master/ebook/03.8.md
    select: 'AES-256-GCM',
    items: [
      {label: '(荐) CHACHA20-IETF-POLY1305', value: 'CHACHA20-IETF-POLY1305'},
      {label: '(荐) AES-128-GCM', value: 'AES-128-GCM'},
      {label: '(荐) AES-256-GCM', value: 'AES-256-GCM'},

      {label: '(弱) bf-cfb', value: 'bf-cfb'},
      {label: '(弱) chacha20', value: 'chacha20'},
      {label: '(弱) salsa20', value: 'salsa20'},
      {label: '(弱) rc4-md5', value: 'rc4-md5'},

      {label: '(危) aes-128-ctr', value: 'aes-128-ctr'},
      {label: '(危) aes-192-ctr', value: 'aes-192-ctr'},
      {label: '(危) aes-256-ctr', value: 'aes-256-ctr'},
      {label: '(危) aes-128-cfb', value: 'aes-128-cfb'},
      {label: '(危) aes-192-cfb', value: 'aes-192-cfb'},
      {label: '(危) aes-256-cfb', value: 'aes-256-cfb'},
      {label: '(危) chacha20-ietf', value: 'chacha20-ietf'},

    ],
    rule: [],
  },
  localSsPassword: {
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
    local_web_server: !!formData.value.localWebServer.select,
    local_web_server_port: formData.value.localWebServerPort.value,
    local_web_server_path: formData.value.localWebServerPath.value,
    local_tcp_server: !!formData.value.localTcpServer.select,
    local_tcp_server_port: formData.value.localTcpServerPort.value,
    local_tcp_server_response: formData.value.localTcpServerResponse.value,
    local_udp_server: !!formData.value.localUdpServer.select,
    local_udp_server_port: formData.value.localUdpServerPort.value,
    local_udp_server_response: formData.value.localUdpServerResponse.value,
    local_ss_server: !!formData.value.localSsServer.select,
    local_ss_port: +formData.value.localSsPort.value,
    local_ss_cipher: formData.value.localSsCipher.select,
    local_ss_password: formData.value.localSsPassword.value,
    log: !!formData.value.log.select,
  }).then(resp => {
    refMySnackbar.value.show('保存完成')
  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
    reloadAppConfig()
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

const handleLocalTcp = (port) => {
  if (port) {
    return `tcp://127.0.0.1:${port}`
  }
  return ''
}

const handleLocalUdp = (port) => {
  if (port) {
    return `udp://127.0.0.1:${port}`
  }
  return ''
}

const reloadAppConfig = () => {
  refMyLoading.value.show()
  AppConfig().then(resp => {
    const appConfig = resp

    setAppConfig(appConfig)

    formData.value.apiServer.value = appConfig.api_server
    formData.value.localWebServer.select = +appConfig.local_web_server
    formData.value.localWebServerPort.value = +appConfig.local_web_server_port
    formData.value.localWebServerPath.value = appConfig.local_web_server_path
    formData.value.localTcpServer.select = +appConfig.local_tcp_server
    formData.value.localTcpServerPort.value = +appConfig.local_tcp_server_port
    formData.value.localTcpServerResponse.value = appConfig.local_tcp_server_response
    formData.value.localUdpServer.select = +appConfig.local_udp_server
    formData.value.localUdpServerPort.value = +appConfig.local_udp_server_port
    formData.value.localUdpServerResponse.value = appConfig.local_udp_server_response
    formData.value.localSsServer.select = +appConfig.local_ss_server
    formData.value.localSsPort.value = +appConfig.local_ss_port
    formData.value.localSsCipher.value = appConfig.local_ss_cipher
    formData.value.localSsPassword.value = appConfig.local_ss_password
    formData.value.log.select = +appConfig.log
    formData.value.logPath.value = appConfig.log_path

  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
  })
}

const onMountedHandler = () => {
  reloadAppConfig()
}

const onClickReset = () => {
  AppConfigReset().then(resp => {
    refMySnackbar.value.show('重置完成')
  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
    reloadAppConfig()
  })
}

export default defineComponent({
  components: {MyLoading, MySnackbar},
  setup() {

    onMounted(onMountedHandler)

    return {
      refForm,
      formData,
      onClickSubmit,
      onClickReset,
      showTipsModal,
      refMySnackbar,
      refMyLoading,
      copyToClipboard,
      handleLocalUrl,
      handleLocalTcp,
      handleLocalUdp,
    }
  }
})

</script>

<style scoped lang="scss">
.form {
}
</style>
