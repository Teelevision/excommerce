<template>
  <v-app dark>
    <v-app-bar fixed app>
      <v-toolbar-title v-text="title" />
      <v-spacer />
      <v-btn icon>
        <v-badge :content="cartSize" bottom left :value="cartSize > 0">
          <v-icon>mdi-cart</v-icon>
        </v-badge>
      </v-btn>
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
