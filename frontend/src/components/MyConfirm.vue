<template>

  <v-dialog v-model="data.show" max-width="400" :persistent="data.persistent">
    <v-card prepend-icon="info" title="提示">
      <template v-slot:text>
        <div v-html="data.msg"></div>
      </template>
      <template v-slot:actions>
        <v-spacer></v-spacer>
        <v-btn @click="onClickCancel">取消</v-btn>
        <v-btn @click="onClickConfirm">确认</v-btn>
      </template>
    </v-card>
  </v-dialog>

</template>

<script>
import {defineComponent, ref} from "vue";

const data = ref({
  show: false,
  msg: '',
  persistent: false,
  confirmCallback: null,
  cancelCallback: null,
})


const show = (msg, {confirmCallback, cancelCallback, persistent = true}) => {
  data.value.show = true
  data.value.msg = msg
  data.value.persistent = persistent
  data.value.confirmCallback = confirmCallback
  data.value.cancelCallback = cancelCallback
}

const hide = () => {
  data.value.show = false
}

const onClickConfirm = () => {
  if (data.value.confirmCallback) {
    data.value.confirmCallback()
  }
  data.value.show = false
}

const onClickCancel = () => {
  if (data.value.cancelCallback) {
    data.value.cancelCallback()
  }
  data.value.show = false
}

export default defineComponent({
  setup() {
    return {
      data,
      show,
      hide,
      onClickConfirm,
      onClickCancel,
    }
  }
})

</script>

<style scoped  lang="scss">

</style>
