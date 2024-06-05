<template>
  <div>
    <Head style="--wails-draggable:drag"></Head>
    <div style="--wails-draggable:no-drag">
      <router-view />
    </div>
  </div>
</template>

<script>
import {defineComponent, onBeforeMount} from "vue";
import Head from './components/Head.vue'
import {useRoute, useRouter} from "vue-router";

let route = null
let router = null

const onBeforeMountHandler = () => {

  console.log('[onBeforeMountHandler1]', route)
  console.log('[onBeforeMountHandler2]', router)
  console.log('[onBeforeMountHandler3]', router.getRoutes())
}

export default defineComponent({
  components: {Head},
  setup() {
    route = useRoute()
    router = useRouter()
    onBeforeMount(onBeforeMountHandler)
    return {
    }
  },
})

</script>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
