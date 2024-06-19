<template>
  <div>

    <v-form ref="formRef" style="width: 400px; margin: 0 auto; padding-top: 100px;">
      <v-container>

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

    <MySnackbar ref="mySnackbar"></MySnackbar>
    <MyLoading ref="myLoading"></MyLoading>


  </div>
</template>

<script>
import {defineComponent, getCurrentInstance, onMounted, ref} from "vue";
import api from "../common/api.js";
import Snackbar from "../components/MySnackbar.vue";
import MySnackbar from "../components/MySnackbar.vue";
import MyLoading from "../components/MyLoading.vue";

let inst = null


const formRef = ref(null)
const formData = ref({
  username: {
    value: '',
    rule: [
      value => {
        if (!value || value.length < 6 || value.length > 20) {
          return '请输入用户名(6-20位)';
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

  api.login({username: formData.value.username.value, password: formData.value.password.value}).then(resp => {
    console.log('[api.login]', resp)

    console.log('[]', inst.$refs.mySnackbar.show('OK'))
  }).catch(err => {
    console.log('[api.login]', err)
    console.log('[]', inst.$refs.mySnackbar.show(err))
  })
  // axios.hasOwnProperty()

}

export default defineComponent({
  components: {MyLoading, MySnackbar, Snackbar},
  setup() {
    onMounted(() => {
      inst = getCurrentInstance().ctx

      inst.$refs.myLoading.show()
      setTimeout(() => {
        inst.$refs.myLoading.hide()
      }, 1000)
      // inst.$refs.mySnackbar.show('OK')

    })
    return {
      formRef,
      formData,
      onClickSubmit,
    }
  }
})

</script>

<style scoped>


</style>
