<template>
  <div class="">

    <div class="ml-2 mr-2">
      <div class="flex-items-between ml-2 mr-2">
        <div>
          <v-chip class="pointer" color="primary" label @click="onClickLoadProxies">
            <v-icon icon="view_list" start></v-icon>
            规则列表
          </v-chip>

          <!--
          <v-chip class="pointer ml-3" color="primary" label>
            启用
          </v-chip>

          <v-chip class="pointer ml-3" color="primary" label>
            禁用
          </v-chip>
          -->

        </div>
        <div class="flex-items-between">
          <v-chip class="pointer" color="primary" label @click="onClickReload" >
            <v-icon icon="refresh" start></v-icon>
            重载服务
          </v-chip>
          <div class="ma-2"></div>
          <v-chip class="pointer" color="primary" label @click="$router.push('/proxy-edit')" >
            <v-icon icon="add" start></v-icon>
            创建规则
          </v-chip>
        </div>
      </div>
    </div>

    <v-item-group selected-class="bg-primary" v-if="proxies.list">
      <v-container>
        <v-row>
          <v-col
              v-for="(proxy,index) in proxies.list"
              :key="index"
              cols="12"
              md="6"
              @click="onClickFrpcStart"
          >
            <v-item v-slot="{ isSelected, selectedClass, toggle }">
              <v-card
                  :class="['d-flex card-item', 'selectedClass']"
                  height="238"
                  dark
                  @click="toggle"
              >
                <div class="flex-grow-1 item">
                  <div>
                    <div class="overflow-hidden">
                      <v-chip class="ma-2" color="primary" label>
                        <v-icon icon="sync_alt" start></v-icon>
                        <b>{{ handleProxyTypeName(proxy.proxy_type) }}</b>
                      </v-chip>
                      <b>{{ proxy.proxy_name }}</b>
                    </div>
                    <div class="ma-2 item-line">
                      <span class="label">内网：</span>
                      <span class="dashed small-font overflow-hidden">{{ proxy.proxy_local_addr }}</span>
                      <!--
                      <v-icon icon="md:content_copy"></v-icon>
                      -->
                    </div>
                    <div class="ma-2 item-line">
                      <span class="label">公网：</span>
                      <span class="dashed small-font overflow-hidden">{{ handleProxyDomain(proxy) }}</span>
                      <!--
                      <v-icon icon="md:content_copy"></v-icon>
                      -->
                    </div>
                    <div class="ma-2 item-line">
                      <span class="label">状态：</span>
                      <span class="small-font">{{ handleProxyStatusName(proxy.status) }}</span>&nbsp;
                      <span class="small-font">
                        <v-icon v-if="proxy.status===1" icon="sync" color="green" start></v-icon>
                        <v-icon v-else icon="sync_disabled" color="orange" start></v-icon>
                      </span>
                    </div>
                    <div class="ma-2 item-line" style="font-size: 12px; color: grey">
                      <span class="label">创建时间：</span>
                      <span>{{ timeFormat(proxy.created_at) }}</span>
                    </div>
                    <div class="ma-2 flex-items-end">
                      <v-chip color="orange" label class="ml-3" @click="onClickUpdateProxy(proxy)">
                        <v-icon icon="view_list" start></v-icon>
                        修改
                      </v-chip>
                      <v-chip color="red" label class="ml-3" @click="onClickDeleteProxy(proxy)">
                        <v-icon icon="view_list" start></v-icon>
                        删除
                      </v-chip>
                    </div>
                  </div>
                </div>
              </v-card>
            </v-item>
          </v-col>
        </v-row>
      </v-container>
    </v-item-group>

    <div v-else>
      <MyEmpty/>
    </div>

    <MySnackbar ref="refMySnackbar"></MySnackbar>
    <MyLoading ref="refMyLoading"></MyLoading>
    <MyConfirm ref="refMyConfirm"></MyConfirm>

  </div>
</template>

<script>
import {defineComponent, getCurrentInstance, onMounted, ref} from "vue"
import {FrpcLogin, FrpcStart} from '../../wailsjs/go/main/App.js'
import api from "../common/api.js";
import MyLoading from "../components/MyLoading.vue";
import MySnackbar from "../components/MySnackbar.vue";
import {handleProxyStatusName, handleProxyTypeName} from "../common/types.js";
import MyEmpty from "../components/MyEmpty.vue";
import {timeFormat} from "../common/helper.js";
import {useRoute, useRouter} from "vue-router";
import {handleProxyDomain} from '../common/proxy.js'
import MyConfirm from "../components/MyConfirm.vue";

let inst = null
let route = null
let router = null

const proxies = ref({})
const refMyLoading = ref(MyLoading)
const refMySnackbar = ref(MySnackbar)
const refMyConfirm = ref(MyConfirm)

const onClickFrpcStart = () => {
  console.log('[onClickFrpcStart]')

  FrpcLogin().then(resp => {
    console.log('[FrpcLoginOk]', resp)
  }).catch(err => {
    console.log('[FrpcLoginErr]', err)
  })
}

const loadProxies = () => {
  refMyLoading.value.show()
  api.getProxies({limit: 20, page: 1}).then(resp => {
    console.log('[getProxies]', resp)
    proxies.value = resp.data.data
  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
  })
}

const deleteProxies = (id, callback = null) => {
  refMyLoading.value.show()
  api.deleteProxy(id).then(resp => {
  }).catch(err => {
    refMySnackbar.value.show(err)
  }).finally(() => {
    refMyLoading.value.hide()
    if (callback) {
      callback()
    }
  })
}

const onMountedHandler = () => {
  console.log('[onMountedHandler]', route)
  inst = getCurrentInstance().ctx
  loadProxies()
}

const onClickDeleteProxy = (proxy) => {
  refMyConfirm.value.show(`确定删除规则 <b>[${proxy.proxy_name}]</b> 吗？`, {
    confirmCallback: () => {
      deleteProxies(proxy.id, () => {
        loadProxies()
      })
    }, cancelCallback: () => {
    }
  })
}

const onClickUpdateProxy = (proxy) => {
  router.push({path: '/proxy-edit', query: {id: proxy.id}})
}

const onClickReload = () => {
  refMyLoading.value.show()
  FrpcStart().then(resp => {
    console.log('[FrpcStart] then ', resp)
  }).catch(err => {
    console.log('[FrpcStart] err', err)
    refMySnackbar.value.show(err)
  }).finally(() => {
    console.log('[FrpcStart] finally')
    refMyLoading.value.hide()
  })
}

const onClickLoadProxies = () => {
  loadProxies()
}

export default defineComponent({
  components: {MyConfirm, MyEmpty, MySnackbar, MyLoading},
  setup() {
    route = useRoute()
    router = useRouter()
    onMounted(onMountedHandler)
    return {
      onClickFrpcStart,
      proxies,
      handleProxyTypeName,
      handleProxyStatusName,
      timeFormat,
      handleProxyDomain,
      onClickDeleteProxy,
      onClickUpdateProxy,
      refMyLoading,
      refMySnackbar,
      refMyConfirm,
      onClickReload,
      onClickLoadProxies,
    }
  }
})

</script>

<style scoped lang="scss">
.card-item {
  padding: 8px 10px;

  .item {
    width: 100%;
  }

  .item-line {
    display: flex;
    flex-direction: row;
    justify-content: start;
    align-items: center;
    overflow: hidden;
    white-space: nowrap;

    .small-font {
      font-size: 14px;
      color: grey;
    }

    .dashed {
      border-bottom: 1px dashed #c3c3c3;
    }
  }

  .selectedClass {
    background-color: #c3c3c3 !important;
    color: #ffffff !important;
  }
}
</style>
