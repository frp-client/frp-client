<template>
  <div>

    <v-item-group selected-class="bg-primary">
      <v-container>
        <v-row>
          <v-col
              v-for="n in 3"
              :key="n"
              cols="12"
              md="4"
              @click="onClickFrpcStart"
          >
            <v-item v-slot="{ isSelected, selectedClass, toggle }">
              <v-card
                  :class="['d-flex align-center', selectedClass]"
                  height="200"
                  dark
                  @click="toggle"
              >
                <div
                    class="text-h3 flex-grow-1 text-center"
                >
                  {{ isSelected ? 'Selected' : 'Click Me!' }}
                </div>
              </v-card>
            </v-item>
          </v-col>
        </v-row>
      </v-container>
    </v-item-group>

  </div>
</template>

<script>
import {defineComponent} from "vue"
import {FrpcStart} from '../../wailsjs/go/main/App.js'

const onClickFrpcStart = () => {
  console.log('[onClickFrpcStart]')
  FrpcStart().then(resp => {
    console.log('[FrpcStartOk]', resp)
  }).catch(err => {
    console.log('[FrpcStartError]', err)
  })
}

export default defineComponent({
  setup() {
    return {
      onClickFrpcStart,
    }
  }
})

</script>

<style scoped>

</style>
