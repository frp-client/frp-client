<template>
  <div>

    <v-form ref="formRef" style="width: 400px; margin: 0 auto; padding-top: 100px;">
      <v-container>

        <div v-if="formData.username.value" style="position: absolute; top: 120px;font-weight: bold">
          当前已登录账号：{{ formData.username.value.substring(0, 12) }}...
        </div>
        <v-row>
          <v-col cols="12" md="12">
            <v-text-field
                label="用户名"
                placeholder="请输入用户名"
                v-model="formData.username.value"
                :rules="formData.username.rule"
                variant="underlined"
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="12">
            <v-text-field
                label="密码"
                placeholder="请输入密码"
                v-model="formData.password.value"
                :rules="formData.password.rule"
                type="password"
                variant="underlined"
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="12" class="d-flex justify-end">
            <v-btn color="rgb(24, 103, 192)" @click="onClickSubmit">登 录</v-btn>
          </v-col>
        </v-row>

      </v-container>
    </v-form>

    <MySnackbar ref="refMySnackbar"></MySnackbar>
    <MyLoading ref="refMyLoading"></MyLoading>


  </div>
</template>

<script>
import {defineComponent, getCurrentInstance, onBeforeMount, onMounted, ref} from "vue";
import api from "../common/api.js";
import Snackbar from "../components/MySnackbar.vue";
import MySnackbar from "../components/MySnackbar.vue";
import MyLoading from "../components/MyLoading.vue";
import {LoginSuccess} from "../../wailsjs/go/main/App.js";
import {getSession} from "../common/vars.js";

const refForm = ref(false)
const refMyLoading = ref(MyLoading)
const refMySnackbar = ref(MySnackbar)

const formData = ref({
  username: {
    value: '',
    rule: [
      value => {
        if (!value || value.length < 4 || value.length > 20) {
          return '请输入用户名(4-20位)';
        }
        return true;
      },
    ],
  },
  password: {
    value: '',
    rule: [
      value => {
        if (!value || value.length < 6) {
          return '请输入密码(不少于6位)';
        }
        return true;
      },
    ],
  },
})

const onClickSubmit = async () => {

  const {valid} = await inst.$refs.formRef.validate()
  if (!valid) {
    return
  }

  inst.$refs.myLoading.show()
  api.login({
    username: formData.value.username.value,
    password: formData.value.password.value,
  }).then(resp => {
    LoginSuccess(resp.data.data).finally(() => {
      inst.$refs.mySnackbar.show('登录成功')
    })
  }).catch(err => {
    inst.$refs.mySnackbar.show(err)
  }).finally(() => {
    inst.$refs.myLoading.hide()
  })

}

const onBeforeMountHandler = () => {
  const session = getSession()
  formData.value.username.value = session.username
  formData.value.password.value = ''
}

export default defineComponent({
  components: {MyLoading, MySnackbar, Snackbar},
  setup() {
    onMounted(() => {
      inst = getCurrentInstance().ctx

      // inst.$refs.myLoading.show()
      // setTimeout(() => {
      //   inst.$refs.myLoading.hide()
      // }, 1000)
      // inst.$refs.mySnackbar.show('OK')

    })
    onBeforeMount(onBeforeMountHandler)
    return {
      formData,
      onClickSubmit,
      refMySnackbar,
      refMyLoading,
      refForm,

    }
  }
})

</script>

<style scoped lang="scss">


</style>
