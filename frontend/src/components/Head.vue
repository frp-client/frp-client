<template>
  <div>
    <div class="head plr20 flex-items-between">
      <div class="flex-items-center">
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#5f6368">
          <path
              d="M200-160v-280h-80v-80h240v80h-80v280h-80Zm0-440v-200h80v200h-80Zm160 0v-80h80v-120h80v120h80v80H360Zm80 440v-360h80v360h-80Zm240 0v-120h-80v-80h240v80h-80v120h-80Zm0-280v-360h80v360h-80Z"/>
        </svg>
        <span class="title">frp客户端</span>
      </div>
      <div class="flex-items-center">
        <!--
        <Icon class="icon menu" size="24">
          <MenuFilled/>
        </Icon>
        -->
        <Icon class="icon minus" size="24" @click="Hidden()">
          <HorizontalRuleFilled/>
        </Icon>
        <Icon class="icon close" size="24" @click="showCloseModal=true">
          <CloseFilled/>
        </Icon>
      </div>
    </div>

    <v-dialog v-model="showCloseModal" width="500">
      <v-card text="确定关闭程序？">
        <template v-slot:title>
          <div class="flex-items-center">
            <v-icon icon="info"/>&nbsp;提示
          </div>
        </template>
        <template v-slot:actions>
          <v-btn text="取消" @click="showCloseModal = false"></v-btn>
          <v-btn text="确定" @click="onClickQuit"></v-btn>
        </template>
      </v-card>
    </v-dialog>

  </div>
</template>

<script>
import {defineComponent, onBeforeMount, ref} from "vue";
import {Icon} from '@vicons/utils'
import {CloseFilled, HorizontalRuleFilled, MenuFilled, MinusFilled, SettingsFilled} from "@vicons/material";
import {Hidden, Quit} from "../../wailsjs/go/main/App.js";

const showCloseModal = ref(null)

const onBeforeMountHandler = () => {

}

const onClickQuit = () => {
  Quit()
}

export default defineComponent({
  components: {Icon, CloseFilled, HorizontalRuleFilled, MinusFilled, SettingsFilled, MenuFilled},
  setup() {
    onBeforeMount(onBeforeMountHandler)
    return {
      showCloseModal,
      onClickQuit,
      Hidden,
    }
  }
})

</script>

<style scoped  lang="scss">
.head {
  height: 42px;
  line-height: 42px;
  cursor: grab;
  background-color: ghostwhite;

  .title {
    margin-left: 6px;
  }
}

.icon {
  cursor: pointer;
  margin: 0 5px;

  transition: transform 0.5s ease;
}

.icon:hover {
  margin-top: 3px;
}

.icon.menu:hover {
  color: #444444;
}

.icon.minus:hover {
  color: grey;
}

.icon.close:hover {
  color: red;
  transform: rotate(180deg);
}


</style>
