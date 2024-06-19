<template>

  <v-overlay
      v-model="data.show"
      :persistent="data.persistent"
      style="--wails-draggable:drag"
      class="flex-items-center flex-justify-content-center">

    <v-progress-circular
        color="primary"
        :size="60"
        indeterminate
    ></v-progress-circular>

  </v-overlay>

</template>

<script>
import {defineComponent, ref} from "vue";

const data = ref({
  show: false,
  persistent: true,
})

const resetSnackbar = () => {
  data.value.show = false
}

const show = (persistent = true) => {
  data.value.show = true
  data.value.persistent = persistent

}

const hide = () => {
  data.value.show = false
}

export default defineComponent({
  setup() {
    return {
      data,
      show,
      hide,
    }
  }
})

</script>

<style scoped>

</style>
