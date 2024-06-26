<template>
  <div>

    <v-form ref="formRef" class="form">
      <v-container>
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
          <v-col cols="3" md="3">
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
          <v-col cols="3" md="3" v-if="formData.localWebServer.select">
            <v-text-field
                label="端口"
                v-model="formData.localWebServerPort.value"
                :rules="formData.localWebServerPort.rule"
                type="number"
                variant="underlined"
            ></v-text-field>
          </v-col>
          <v-col cols="6" md="6" class="flex-items-center" v-if="formData.localWebServer.select">
            <text class="dashed-x pointer" @click="copyToClipboard(handleLocalUrl(formData.localWebServerPort.value))">
              {{ handleLocalUrl(formData.localWebServerPort.value) }}
            </text>
            &nbsp;
            <v-icon v-if="s" icon="circle" color="green"></v-icon>
            <v-icon v-else icon="circle" color="red"></v-icon>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="6" md="6">
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

  </div>
</template>

<script>
import {defineComponent, getCurrentInstance, onMounted, ref} from "vue";
import clipboard from "../common/clipboard.js";
import MySnackbar from "../components/MySnackbar.vue";

let inst = null

const showTipsModal = ref(false)
const refMySnackbar = ref(MySnackbar)


const formRef = ref(null)
const formData = ref({
  expandModel: [0],
  proxyType: {
    select: {label: 'http', value: 'http'},
    items: [
      {label: 'http', value: 'http'},
      {label: 'https', value: 'https'},
      {label: 'tcp', value: 'tcp'},
      {label: 'udp', value: 'udp'},
    ],
    rule: [
      value => {
        if (!value) {
          return '请选择代理类型';
        }
        return true;
      },
    ],

  },
  proxyName: {
    value: '',
    rule: [
      value => {
        if (!value || value.length < 2) {
          return '请输入代理名称(不少于2位)';
        }
        return true;
      },
    ],
  },
  apiServer: {
    value: 'https://tunnel-api.lixiang4u.xyz:3000',
    rule: [
      value => {
        return true;
      },
    ],
  },
  remotePort: {
    value: '',
    rule: [],
  },
  domain: {
    value: '',
    rule: [],
  },
  localWebServer: {
    select: {label: '开启', value: 1},
    items: [
      {label: '开启', value: 1},
      {label: '关闭', value: 0},
    ],
    rule: [],
  },
  localWebServerPort: {
    value: 8080,
    rule: [],
  },
  log: {
    select: {label: '开启日志', value: 1},
    items: [
      {label: '开启日志', value: 1},
      {label: '关闭日志', value: 0},
    ],
    rule: [],

  },
  sslCrt: {
    value: '',
    rule: [],
  },
  sslKey: {
    value: '',
    rule: [],
  }
})

const onClickSubmit = async () => {
  console.log('[onClickSubmit]', formData.value)


  console.log('[formRef1]', inst.$refs.formRef)

  const {valid} = await inst.$refs.formRef.validate()

  console.log('[valid]', valid)
  if (valid) {
    alert('Form is valid')
  }

  // showTipsModal.value=true

}

const copyToClipboard = (url) => {
  clipboard.write(url, false)
  refMySnackbar.value.show('已复制：' + url, {timeout: 1000})
}

const handleLocalUrl = (port) => {
  return `http://127.0.0.1:${port}`
}

export default defineComponent({
  components: {MySnackbar},
  setup() {

    onMounted(() => {
      inst = getCurrentInstance().ctx
    })

    return {
      formRef,
      formData,
      onClickSubmit,
      showTipsModal,
      refMySnackbar,
      copyToClipboard,
      handleLocalUrl,
    }
  }
})

</script>

<style scoped lang="scss">
.form {
  padding: 20px 70px;
}
</style>
