<template>
  <v-app dark>
    <v-app-bar fixed app>
      <v-toolbar-title>
        <nuxt-link to="/" style="color: white; text-decoration: none;">{{
          title
        }}</nuxt-link>
      </v-toolbar-title>
      <v-spacer />
      <nuxt-link to="/cart" style="text-decoration: none;">
        <v-btn icon>
          <v-badge :content="cartSize" bottom left :value="cartSize > 0">
            <v-icon>mdi-cart</v-icon>
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
import { mapState } from 'vuex'

export default {
  data() {
    return {
      title: 'ExCommerce'
    }
  },
  computed: mapState({
    cartSize: (state) =>
      state.cart.positions.reduce(
        (prev, position) => prev + position.quantity,
        0
      )
  })
}
</script>
