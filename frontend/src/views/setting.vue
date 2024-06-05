<template>
  <div>
    <v-form v-model="valid">
      <v-container>
        <v-row>
          <v-col cols="12" md="12">
            <v-text-field
                v-model="firstname"
                :counter="10"
                :rules="nameRules"
                label="First name"
                hide-details
                required
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="12">
            <v-text-field
                v-model="lastname"
                :counter="10"
                :rules="nameRules"
                label="Last name"
                hide-details
                required
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="12">
            <v-text-field
                v-model="email"
                :rules="emailRules"
                label="E-mail"
                hide-details
                required
            ></v-text-field>
          </v-col>
        </v-row>

      </v-container>
    </v-form>
  </div>
</template>

<script >
console.log('[SETTING]')
</script>

<style scoped>

</style>