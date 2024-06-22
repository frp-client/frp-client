<template>
  <div>

    <v-form ref="formRef">
      <v-container>
        <v-expansion-panels v-model="formData.expandModel" class="mb-5">
          <v-expansion-panel>
            <v-expansion-panel-title :readonly="true">
              <v-chip class="" color="primary" label>
                <v-icon icon="tune" start></v-icon>
                <span>基础配置</span>
              </v-chip>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
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

            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>

        <v-expansion-panels class="mb-5">
          <v-expansion-panel>
            <v-expansion-panel-title>
              <v-chip class="" color="primary" label>
                <v-icon icon="more_horiz" start></v-icon>
                <span>更多配置</span>
              </v-chip>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <v-row>
                <v-col cols="6" md="6">
                  <v-text-field
                      label="公网域名"
                      placeholder="默认自动分配"
                      v-model="formData.domain.value"
                      :rules="formData.domain.rule"
                      variant="underlined"
                  ></v-text-field>
                </v-col>
                <v-col cols="6" md="6">
                  <v-text-field
                      label="公网端口"
                      placeholder="默认自动分配"
                      v-model="formData.remotePort.value"
                      :rules="formData.remotePort.rule"
                      type="number"
                      variant="underlined"
                  ></v-text-field>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="6" md="6">
                  <v-file-input
                      label="证书文件"
                      placeholder="默认自动分配"
                      v-model="formData.sslCrt.rule"
                      :rules="formData.sslCrt.rule"
                      variant="underlined"
                  ></v-file-input>

                </v-col>
                <v-col cols="6" md="6">
                  <v-file-input
                      label="证书密钥"
                      placeholder="默认自动分配"
                      v-model="formData.sslKey.rule"
                      :rules="formData.sslKey.rule"
                      variant="underlined"
                  ></v-file-input>

                </v-col>
              </v-row>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>

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

    <MySnackbar ref="mySnackbar"></MySnackbar>
    <MyLoading ref="myLoading"></MyLoading>


  </div>
</template>

<script>
import {defineComponent, getCurrentInstance, onMounted, ref} from "vue";
import api from "../common/api.js";
import MyLoading from "../components/MyLoading.vue";
import MySnackbar from "../components/MySnackbar.vue";
import {useRoute, useRouter} from "vue-router";
import {handleProxyDomain} from '../common/proxy.js'
import MyConfirm from "../components/MyConfirm.vue";

let route = null
let router = null
let inst = null
let id = null

const showTipsModal = ref(false)

const formRef = ref(null)
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
    },
    status: +formData.value.proxyStatus.select,
  }
}

const onClickSubmit = async () => {
  const {valid} = await inst.$refs.formRef.validate()
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
  inst.$refs.myLoading.show()
  api.createProxy(data).then(resp => {
    console.log('[api.createProxy]', resp)
    // inst.$refs.mySnackbar.show('创建完成，即将跳转...')
    router.push('/proxy')
  }).catch(err => {
    inst.$refs.mySnackbar.show(err)
  }).finally(() => {
    inst.$refs.myLoading.hide()
  })
}

const updateProxy = (id, data) => {
  inst.$refs.myLoading.show()
  api.updateProxy(id, data).then(resp => {
    console.log('[api.createProxy]', resp)
    // inst.$refs.mySnackbar.show('修改完成，即将跳转...')
    router.push('/proxy')
  }).catch(err => {
    inst.$refs.mySnackbar.show(err)
  }).finally(() => {
    inst.$refs.myLoading.hide()
  })
}

const loadProxy = (id) => {
  inst.$refs.myLoading.show()
  api.getProxy(id).then(resp => {
    const tmpForm = resp.data.data
    formData.value.proxyType.items.filter(item => {
      if (item.value === tmpForm.proxy_type) {
        formData.value.proxyType.select = tmpForm.proxy_type
      }
    })
    formData.value.proxyName.value = tmpForm.proxy_name
    formData.value.localAddr.value = tmpForm.proxy_local_addr
    formData.value.proxyStatus.value = tmpForm.status
    formData.value.remotePort.value = tmpForm.proxy_remote_port
    formData.value.sslCrt.value = tmpForm.proxy_extra.ssl_crt
    formData.value.sslKey.value = tmpForm.proxy_extra.ssl_key
    // 公网域名
    formData.value.domain.value = handleProxyDomain(tmpForm)
  }).catch(err => {
    inst.$refs.mySnackbar.show(err)
  }).finally(() => {
    inst.$refs.myLoading.hide()
  })
}

const onMountedHandler = () => {
  inst = getCurrentInstance().ctx

  resetFormData()
  id = route.query['id']
  if (id) {
    loadProxy(id)
  }
}

export default defineComponent({
  components: {MyConfirm, MySnackbar, MyLoading},
  setup() {
    route = useRoute()
    router = useRouter()
    onMounted(onMountedHandler)
    return {
      formRef,
      formData,
      onClickSubmit,
      showTipsModal,
    }
  }
})

</script>

<style scoped lang="scss">

</style>
