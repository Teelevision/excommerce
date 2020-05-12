<template>
  <v-layout>
    <v-flex class="text-center">
      <v-container fluid>
        <v-row>
          <v-col v-for="product in products" :key="product.id" cols="6">
            <v-card>
              <v-img :src="product.img" max-height="300px" max-width="500px">
                <v-card-title v-text="product.name"></v-card-title>
              </v-img>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="primary" @click="addToCart(product.id)">
                  <v-icon left>mdi-cart-plus</v-icon>
                  {{ product.price }} EUR
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-flex>
  </v-layout>
</template>

<script>
import { mapState, mapActions } from 'vuex'

export default {
  computed: mapState({
    products: (state) =>
      state.products.map((p) => ({
        ...p,
        img: 'https://picsum.photos/seed/' + p.name + '/500/300'
      }))
  }),
  mounted() {
    this.loadAllProducts()
  },
  methods: mapActions(['loadAllProducts', 'addToCart'])
}
</script>
