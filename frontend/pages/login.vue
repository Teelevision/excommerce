<template>
  <v-row justify="center">
    <v-col cols="12" sm="10" md="8" lg="6">
      <v-card ref="form">
        <v-card-text>
          <v-text-field
            ref="name"
            v-model="name"
            :rules="[(v) => !!v || 'The username is required.']"
            :error-messages="apiErrors['/name'] || ''"
            label="Username"
            required
          ></v-text-field>
          <v-text-field
            ref="password"
            v-model="password"
            type="password"
            :rules="[(v) => !!v || 'The password is required.']"
            :error-messages="apiErrors['/password'] || ''"
            label="Password"
            required
          ></v-text-field>
        </v-card-text>
        <v-divider class="mt-12"></v-divider>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="default" @click="register">Register</v-btn>
          <v-btn color="primary" @click="doLogin">Login</v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import { mapActions } from 'vuex'
import { UsersApi } from '../client'
import { User } from '~/models'

export default {
  data: () => ({
    name: '',
    password: '',
    apiErrors: {}
  }),
  methods: {
    ...mapActions(['login']),
    doLogin() {
      this.apiErrors = {}
      this.$refs.name.validate(true)
      this.$refs.password.validate(true)

      this.login(new User(this.name, this.password))
        .then(() => this.$router.push('/'))
        .catch(({ response: { status } }) => {
          switch (status) {
            case 404:
              this.apiErrors = {
                '/name': 'Cannot login. Name or password invalid.'
              }
              break
          }
        })
    },
    register() {
      this.apiErrors = {}
      this.$refs.name.validate(true)
      this.$refs.password.validate(true)

      const user = new User(this.name, this.password)

      new UsersApi()
        .register(user)
        .then(({ data: { id } }) => ({ ...user, id }))
        .then((user) => this.login(user))
        .then(() => this.$router.push('/'))
        .catch(({ response: { status, data } }) => {
          switch (status) {
            case 409:
              this.apiErrors = {
                '/name': 'Cannot register name. Name is already taken.'
              }
              break
            case 422:
              this.apiErrors = { [data.pointer]: data.message }
              break
          }
        })
    }
  }
}
</script>
