<template>
  <v-layout>
    <v-flex class="text-center">
      <v-container>
        <v-row class="d-none d-sm-flex">
          <v-col cols="12" sm="3" md="6" class="text-left">
            <small>Product</small>
          </v-col>
          <v-col cols="12" sm="3" md="2" class="text-right">
            <small>Price per item</small>
          </v-col>
          <v-col cols="12" sm="3" md="2" class="text-right">
            <small>Price</small>
          </v-col>
        </v-row>
        <v-row v-if="!positions.length">
          <v-col cols="12">
            Your cart is empty. Go
            <nuxt-link to="/">here</nuxt-link>
            to add delicious fruits.
          </v-col>
        </v-row>
        <v-row
          v-for="(position, i) in positions"
          :key="i"
          :style="{
            'background-color': i % 2 ? 'transparent' : 'rgba(255,255,255,.1)'
          }"
        >
          <v-col cols="6" sm="3" md="6" class="text-left">
            {{ position.product.name }}
          </v-col>
          <v-col cols="6" sm="3" md="2" order-sm="4">
            <v-text-field
              v-if="position.product.id"
              v-model.number="quantity[i]"
              type="number"
              label="Quantity"
              outlined
              dense
              hide-details="auto"
              style="max-width: 120px;"
              class="float-right"
              :rules="[(v) => v > 0 && v <= 99]"
              append-outer-icon="mdi-cart-off"
              @input="quantityChanged"
              @click:append-outer="() => removePosition(i)"
            />
          </v-col>
          <v-col cols="6" sm="3" md="2" class="text-left text-sm-right">
            <span class="d-none d-sm-inline">
              EUR {{ position.product.price.toFixed(2) }}
            </span>
            <small class="d-sm-none">
              EUR {{ position.product.price.toFixed(2) }} / item
            </small>
          </v-col>
          <v-col cols="6" sm="3" md="2" class="text-right">
            EUR
            {{
              (position.price || position.product.price * quantity[i]).toFixed(
                2
              )
            }}
          </v-col>
        </v-row>
        <v-row
          class="font-weight-black"
          :style="{
            'background-color':
              positions.length % 2 ? 'transparent' : 'rgba(255,255,255,.1)'
          }"
        >
          <v-col cols="6" sm="6" md="8" class="text-right">
            TOTAL
          </v-col>
          <v-col cols="6" sm="3" md="2" class="text-right">
            EUR
            {{
              positions
                .reduce(
                  (prev, p) => prev + (p.price || p.product.price * p.quantity),
                  0
                )
                .toFixed(2)
            }}
          </v-col>
        </v-row>
        <v-row
          v-if="totalSaved > 0"
          class="font-weight-black"
          :style="{
            'background-color':
              positions.length % 2 ? 'rgba(255,255,255,.1)' : 'transparent'
          }"
        >
          <v-col cols="6" sm="6" md="8" class="text-right">
            SAVED
          </v-col>
          <v-col cols="6" sm="3" md="2" class="text-right">
            EUR {{ totalSaved.toFixed(2) }}
          </v-col>
        </v-row>
        <v-row>
          <v-col sm="6" offset-sm="6" class="text-right">
            <v-btn
              block
              large
              color="primary"
              :disabled="!positions.length"
              @click="goToCheckout"
              >Proceed to checkout</v-btn
            >
          </v-col>
        </v-row>
      </v-container>
    </v-flex>
  </v-layout>
</template>

<script>
import { mapState, mapActions } from 'vuex'

export default {
  data: () => ({
    quantity: []
  }),
  computed: {
    ...mapState({
      positions: ({ cart, products }) =>
        cart.positions.map((p) => ({
          ...p,
          product: p.product ||
            Object.values(products).find((pr) => pr.id === p.productId) || {
              name: 'unknown',
              price: 0
            }
        }))
    }),
    totalSaved() {
      return this.positions.reduce((prev, p) => prev + (p.savedPrice || 0), 0)
    }
  },
  watch: {
    positions() {
      this.updateQuantity()
    }
  },
  mounted() {
    this.updateQuantity()
  },
  methods: {
    ...mapActions(['updateCartPositions']),
    updateQuantity() {
      this.quantity = this.positions.map((p) => p.quantity)
    },
    removePosition(i) {
      this.updateCartPositions([
        ...this.positions.slice(0, i),
        ...this.positions.slice(i + 1)
      ]).then(this.updateQuantity)
    },
    quantityChanged() {
      for (const i in this.quantity) {
        if (this.quantity[i] <= 0) {
          this.quantity[i] = 1
        } else if (this.quantity[i] > 99) {
          this.quantity[i] = 99
        }
      }
      this.updateCartPositions(
        this.positions.map((p, i) => ({
          ...p,
          quantity: this.quantity[i]
        }))
      ).then(this.updateQuantity)
    },
    goToCheckout() {
      this.$router.push('/checkout')
    }
  }
}
</script>
