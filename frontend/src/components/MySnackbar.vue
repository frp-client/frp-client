<template>

  <v-overlay v-model="data.show" :persistent="data.persistent" style="--wails-draggable:drag">
    <v-snackbar v-model="data.show" multi-line>
      {{ data.content }}
      <template v-slot:actions>
        <v-btn color="red" variant="text" @click="resetSnackbar">关 闭</v-btn>
      </template>
    </v-snackbar>
  </v-overlay>

</template>

<script>
import {defineComponent, ref} from "vue";

const data = ref({
  show: false,
  content: '',
  timeout: 1000 * 3,
  persistent: false,
})

const resetSnackbar = () => {
  data.value.show = false
  data.value.content = ''
  data.value.timeout = 1000 * 3
  data.value.persistent = false
}

const show = (content, config = {timeout: null, persistent: false}) => {
  data.value.show = true
  data.value.content = content
  data.value.persistent = config.persistent
  if (config.timeout) {
    data.value.timeout = config.timeout
  }

}

export default defineComponent({
  setup() {
    return {
      data,
      resetSnackbar,
      show,
    }
  }
})

</script>

<style scoped lang="scss">

.v-snackbar:not(.v-snackbar--center):not(.v-snackbar--top) {
  align-items: center;
}

</style>
