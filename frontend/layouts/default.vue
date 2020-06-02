<template>
  <v-app dark>
    <v-app-bar fixed app>
      <v-toolbar-title>
        <nuxt-link to="/" style="color: white; text-decoration: none;">{{
          title
        }}</nuxt-link>
      </v-toolbar-title>
      <v-spacer />
      <nuxt-link
        v-if="!user.id"
        to="/login"
        style="text-decoration: none; margin-right: 24px;"
      >
        <v-btn text><v-icon left>mdi-account-circle</v-icon> Login</v-btn>
      </nuxt-link>
      <v-btn v-if="user.id" text @click="doLogout">
        <v-icon left>mdi-account-circle</v-icon>
        Logout ({{ user.name }})
      </v-btn>
      <nuxt-link to="/cart" style="text-decoration: none;">
        <v-btn icon>
          <v-badge :content="cartSize" bottom left :value="cartSize > 0">
            <v-icon>mdi-{{ cartSize ? 'cart' : 'cart-outline' }}</v-icon>
          </v-badge>
        </v-btn>
      </nuxt-link>
    </v-app-bar>
    <v-content>
      <v-container>
        <nuxt />
      </v-container>
    </v-content>
    <v-footer app>
      <span>&copy; {{ new Date().getFullYear() }} Marius Neugebauer</span>
    </v-footer>
  </v-app>
</template>

<script>
import { mapState, mapActions } from 'vuex'

export default {
  data() {
    return {
      title: 'ExCommerce'
    }
  },
  computed: mapState({
    user: (state) => state.user,
    cartSize: (state) =>
      state.cart.positions.reduce(
        (prev, position) => prev + position.quantity,
        0
      )
  }),
  methods: {
    ...mapActions(['logout']),
    doLogout() {
      this.logout().then(() => this.$router.push('/'))
    }
  }
}
</script>
