<template>
  <div>

    <v-form ref="refForm">
      <v-container>
        <v-row>
          <v-col cols="12" md="12">
            <v-chip class="" color="primary" label>
              <v-icon icon="tune" start></v-icon>
              <span>基础配置</span>
            </v-chip>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="6" md="6">
            <v-select
                label="代理类型"
                placeholder="请选择代理类型"
                clearable
                v-model="formData.proxyType.select"
                :items="formData.proxyType.items"
                :rules="formData.proxyType.rule"
                item-title="label"
                item-value="value"
                variant="underlined"
            >
            </v-select>
          </v-col>
          <v-col cols="6" md="6">
            <v-text-field
                label="代理名称"
                placeholder="请输入代理名称"
                v-model="formData.proxyName.value"
                :rules="formData.proxyName.rule"
                variant="underlined"
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="6" md="6">
            <v-text-field
                label="本地地址"
                placeholder="请输入本地地址(ip+端口格式)"
                v-model="formData.localAddr.value"
                :rules="formData.localAddr.rule"
                variant="underlined"
            ></v-text-field>
          </v-col>
          <v-col cols="6" md="6">
            <v-select
                label="状态"
                v-model="formData.proxyStatus.select"
                :items="formData.proxyStatus.items"
                :rules="formData.proxyStatus.rule"
                item-title="label"
                item-value="value"
                variant="underlined"
            ></v-select>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="6" md="6">
            <v-chip color="rgb(92, 187, 246)" label>
              <v-icon icon="more_horiz" start></v-icon>
              <span>更多配置</span>
            </v-chip>
          </v-col>
          <v-col cols="6" md="6" class="flex-items-end flex-align-items-center">
            <div @click="onChangeCustomDomain" class="pointer">
              <v-icon
                  v-if="isCustomDomain" icon="select_check_box" style="justify-content: flex-end">
              </v-icon>
              <v-icon v-else icon="check_box_outline_blank"></v-icon>
              自定义配置域名
            </div>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="6" md="6">
            <v-text-field
                label="公网域名"
                placeholder="默认自动分配"
                v-model="formData.domain.value"
                :rules="formData.domain.rule"
                variant="underlined"
                :disabled="isCustomDomain===false"
            ></v-text-field>
          </v-col>
          <v-col cols="6" md="6" v-if="[3,4].includes(formData.proxyType.select)">
            <v-text-field
                label="公网端口"
                placeholder="默认自动分配"
                v-model="formData.remotePort.value"
                :rules="formData.remotePort.rule"
                type="number"
                variant="underlined"
                disabled
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row v-if="formData.proxyType.select===2">
          <v-col cols="6" md="6">
            <v-text-field
                label="证书文件"
                placeholder="默认自动分配"
                v-model="formData.sslCrt.value"
                :rules="formData.sslCrt.rule"
                variant="underlined"
                disabled
            ></v-text-field>

          </v-col>
          <v-col cols="6" md="6">
            <v-text-field
                label="证书密钥"
                placeholder="默认自动分配"
                v-model="formData.sslKey.value"
                :rules="formData.sslKey.rule"
                variant="underlined"
                disabled
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="12" class="d-flex justify-end">
            <v-btn color="primary" @click="onClickSubmit">提交</v-btn>
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
import api from "../common/api.js";
import MyLoading from "../components/MyLoading.vue";
import MySnackbar from "../components/MySnackbar.vue";
import {useRoute, useRouter} from "vue-router";
import {handleProxyDomain} from '../common/proxy.js'
import MyConfirm from "../components/MyConfirm.vue";

let route = null
let router = null
let id = null

const showTipsModal = ref(false)
const refForm = ref(false)
const refMySnackbar = ref(MySnackbar)
const refMyLoading = ref(MyLoading)
const isCustomDomain = ref(false)

const formData = ref({
  expandModel: [0],
  proxyType: {
    select: null,
    items: [
      {label: 'http', value: 1},
      {label: 'https', value: 2},
      {label: 'tcp', value: 3},
      {label: 'udp', value: 4},
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
        if (!value || value.length < 2 || value.length > 20) {
          return '请输入代理名称(2-20位)';
        }
        return true;
      },
    ],
  },
  localAddr: {
    value: '127.0.0.1:8080',
    rule: [
      value => {
        const tmpArr = value.split(':')
        if (tmpArr.length !== 2) {
          return '本地地址格式错误(ip+端口格式)';
        }
        if (!(parseInt(tmpArr[1]) > 1)) {
          return '本地地址端口错误(ip+端口格式)';
        }
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
  proxyStatus: {
    select: 1,
    items: [
      {label: '启用', value: 1},
      {label: '禁用', value: 2},
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

const resetFormData = () => {
  formData.value.proxyType.select = null
  formData.value.proxyName.value = null
  formData.value.localAddr.value = '127.0.0.1:8080'
  formData.value.remotePort.value = null
  formData.value.domain.value = null
  formData.value.proxyStatus.select = 1
  formData.value.sslCrt.value = null
  formData.value.sslKey.value = null
}

const formatProxyForm = () => {
  return {
    proxy_type: +formData.value.proxyType.select,
    proxy_name: formData.value.proxyName.value,
    proxy_remote_port: +formData.value.remotePort.value,
    proxy_local_addr: formData.value.localAddr.value,
    proxy_extra: {
      subdomain: '',
      ssl_crt: formData.value.sslCrt.value,
      ssl_key: formData.value.sslKey.value,
    },
    status: +formData.value.proxyStatus.select,
  }
}

const onClickSubmit = async () => {
  const {valid} = await refForm.value.validate()
  if (!valid) {
    return
  }
  if (id) {
    updateProxy(id, formatProxyForm())
  } else {
    createProxy(formatProxyForm())
  }
}

const createProxy = (data) => {
  refMyLoading.value.show()
  api.createProxy(data).then(resp => {
    console.log('[api.createProxy]', resp)
    router.push('/proxy?reload=1')
  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
  })
}

const updateProxy = (id, data) => {
  refMyLoading.value.show()
  api.updateProxy(id, data).then(resp => {
    console.log('[api.createProxy]', resp)
    router.push('/proxy?reload=1')
  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
  })
}

const loadProxy = (id) => {
  refMyLoading.value.show()
  api.getProxy(id).then(resp => {
    const tmpForm = resp.data.data
    formData.value.proxyType.items.filter(item => {
      if (item.value === tmpForm.proxy_type) {
        formData.value.proxyType.select = tmpForm.proxy_type
      }
    })
    formData.value.proxyName.value = tmpForm.proxy_name
    formData.value.localAddr.value = tmpForm.proxy_local_addr
    formData.value.proxyStatus.select = tmpForm.status
    formData.value.remotePort.value = tmpForm.proxy_remote_port
    formData.value.sslCrt.value = tmpForm.proxy_extra.ssl_crt
    formData.value.sslKey.value = tmpForm.proxy_extra.ssl_key
    // 公网域名
    formData.value.domain.value = handleProxyDomain(tmpForm)
  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
  })
}

const onMountedHandler = () => {
  resetFormData()
  id = route.query['id']
  if (id) {
    loadProxy(id)
  }
}

const onChangeCustomDomain = () => {
  refMySnackbar.value.show('暂未开放')
  // return
  isCustomDomain.value = !isCustomDomain.value
}

export default defineComponent({
  components: {MyConfirm, MySnackbar, MyLoading},
  setup() {
    route = useRoute()
    router = useRouter()
    onMounted(onMountedHandler)
    return {
      formData,
      onClickSubmit,
      showTipsModal,
      refMySnackbar,
      refMyLoading,
      refForm,
      isCustomDomain,
      onChangeCustomDomain,
    }
  }
})

</script>

<style scoped lang="scss">

</style>
