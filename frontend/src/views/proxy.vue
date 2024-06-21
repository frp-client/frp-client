<template>
  <div class="">

    <div class="ml-2 mr-2">
      <div class="flex-items-between ml-2 mr-2">
        <v-chip color="primary" label>
          <v-icon icon="view_list" start></v-icon>
          规则列表
        </v-chip>
        <router-link to="/proxy-edit">
          <v-btn color="primary" variant="tonal">创建规则</v-btn>
        </router-link>
      </div>
    </div>

    <v-item-group selected-class="bg-primary" v-if="proxies.list">
      <v-container>
        <v-row>
          <v-col v-for="(proxy,index) in proxies.list" :key="index" cols="12" md="6" @click="onClickFrpcStart">
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
                      <span class="dashed">{{ proxy.proxy_local_addr }}</span>
                      <!--
                      <v-icon icon="md:content_copy"></v-icon>
                      -->
                    </div>
                    <div class="ma-2 item-line">
                      <span class="label">公网：</span>
                      <span class="dashed overflow-hidden">{{ handleProxyDomain(proxy) }}</span>
                      <!--
                      <v-icon icon="md:content_copy"></v-icon>
                      -->
                    </div>
                    <div class="ma-2 item-line">
                      <span class="label">状态：</span>
                      <span>{{ handleProxyStatusName(proxy.status) }}</span>
                    </div>
                    <div class="ma-2 item-line" style="font-size: 12px; color: grey">
                      <span class="label">创建时间：</span>
                      <span>{{ timeFormat(proxy.created_at) }}</span>
                    </div>
                    <div class="ma-2 flex-items-end">
                      <router-link :to="{path:'/proxy-edit', query:{id: proxy.id}}">
                        <v-chip color="orange" label>
                          <v-icon icon="view_list" start></v-icon>
                          修改
                        </v-chip>
                      </router-link>
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

    <MySnackbar ref="mySnackbar"></MySnackbar>
    <MyLoading ref="myLoading"></MyLoading>
    <MyConfirm ref="myConfirm"></MyConfirm>

  </div>
</template>

<script>
import {defineComponent, getCurrentInstance, onMounted, ref} from "vue"
import {FrpcLogin} from '../../wailsjs/go/main/App.js'
import api from "../common/api.js";
import MyLoading from "../components/MyLoading.vue";
import MySnackbar from "../components/MySnackbar.vue";
import {handleProxyStatusName, handleProxyTypeName} from "../common/types.js";
import MyEmpty from "../components/MyEmpty.vue";
import {timeFormat} from "../common/helper.js";
import {useRoute} from "vue-router";
import {handleProxyDomain} from '../common/proxy.js'
import MyConfirm from "../components/MyConfirm.vue";

let inst = null
const proxies = ref({})
const route = ref({})

const onClickFrpcStart = () => {
  console.log('[onClickFrpcStart]')

  FrpcLogin().then(resp => {
    console.log('[FrpcLoginOk]', resp)
  }).catch(err => {
    console.log('[FrpcLoginErr]', err)
  })
}

const loadProxies = () => {
  inst.$refs.myLoading.show()
  api.getProxies({limit: 20, page: 1}).then(resp => {
    console.log('[getProxies]', resp)
    proxies.value = resp.data.data

    // inst.$refs.mySnackbar.show('登录成功' + resp)
  }).catch(err => {
    inst.$refs.mySnackbar.show(err)
  }).finally(() => {
    inst.$refs.myLoading.hide()
  })
}

const deleteProxies = (id, callback = null) => {
  inst.$refs.myLoading.show()
  api.deleteProxy(id).then(resp => {
  }).catch(err => {
    inst.$refs.mySnackbar.show(err)
  }).finally(() => {
    inst.$refs.myLoading.hide()
    if (callback) {
      callback()
    }
  })
}

const onMountedHandler = () => {
  console.log('[onMountedHandler]', route.value)
  inst = getCurrentInstance().ctx
  loadProxies()
}

const onClickDeleteProxy = (proxy) => {
  inst.$refs.myConfirm.show(`确定删除规则 <b>[${proxy.proxy_name}]</b> 吗？`, {
    confirmCallback: () => {
      deleteProxies(proxy.id, () => {
        loadProxies()
      })
    }, cancelCallback: () => {
    }
  })
}

export default defineComponent({
  components: {MyConfirm, MyEmpty, MySnackbar, MyLoading},
  setup() {
    route.value = useRoute()
    onMounted(onMountedHandler)
    return {
      onClickFrpcStart,
      proxies,
      handleProxyTypeName,
      handleProxyStatusName,
      timeFormat,
      handleProxyDomain,
      onClickDeleteProxy,
    }
  }
})

</script>

<style scoped>
.card-item {
  padding: 8px 10px;

  .item {
    width: 100%;
  }

  .item-line {
    display: flex;
    flex-direction: row;
    justify-content: start;
    overflow: hidden;
    white-space: nowrap;


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
