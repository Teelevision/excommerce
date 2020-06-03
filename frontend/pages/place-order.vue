<template>
  <v-layout>
    <v-flex class="text-center">
      <v-container v-if="order !== null">
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
          <v-col cols="12" sm="3" md="2" class="text-right">
            <small>Quantity</small>
          </v-col>
        </v-row>
        <v-row
          v-for="(position, i) in order.positions"
          :key="i"
          :style="{
            'background-color': i % 2 ? 'transparent' : 'rgba(255,255,255,.1)'
          }"
        >
          <v-col cols="6" sm="3" md="6" class="text-left">
            {{ position.product.name }}
          </v-col>
          <v-col cols="6" sm="3" md="2" order-sm="4" class="text-right">
            <span class="d-none d-sm-inline">
              {{ position.quantity }}
            </span>
            <span class="d-sm-none">x {{ position.quantity }}</span>
          </v-col>
          <v-col cols="6" sm="3" md="2" class="text-left text-sm-right">
            <span v-if="position.quantity != 1" class="d-none d-sm-inline">
              EUR {{ position.product.price.toFixed(2) }}
            </span>
            <small v-if="position.quantity != 1" class="d-sm-none">
              EUR {{ position.product.price.toFixed(2) }} / item
            </small>
          </v-col>
          <v-col cols="6" sm="3" md="2" class="text-right">
            EUR {{ position.price.toFixed(2) }}
          </v-col>
        </v-row>
        <v-row
          class="font-weight-black"
          :style="{
            'background-color':
              positions.length % 2 ? 'transparent' : 'rgba(255,255,255,.1)'
          }"
        >
          <v-col cols="6" sm="6" md="8" style="text-align: right;">
            TOTAL
          </v-col>
          <v-col cols="6" sm="3" md="2" style="text-align: right;">
            EUR {{ order.price.toFixed(2) }}
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" sm="6" style="text-align: left;">
            <small>Billing address</small>
            <v-container>
              <v-row style="background-color: rgba(255,255,255,.1);">
                <v-col cols="4">
                  Name
                </v-col>
                <v-col cols="8">
                  {{ order.buyer.name }}
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="4">
                  Country
                </v-col>
                <v-col cols="8">
                  {{ order.buyer.country }}
                </v-col>
              </v-row>
              <v-row style="background-color: rgba(255,255,255,.1);">
                <v-col cols="4">
                  Postal code
                </v-col>
                <v-col cols="8">
                  {{ order.buyer.postalCode }}
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="4">
                  City
                </v-col>
                <v-col cols="8">
                  {{ order.buyer.city }}
                </v-col>
              </v-row>
              <v-row style="background-color: rgba(255,255,255,.1);">
                <v-col cols="4">
                  Street
                </v-col>
                <v-col cols="8">
                  {{ order.buyer.street }}
                </v-col>
              </v-row>
            </v-container>
          </v-col>
          <v-col cols="12" sm="6" style="text-align: left;">
            <small>Delivery address</small>
            <v-container>
              <v-row style="background-color: rgba(255,255,255,.1);">
                <v-col cols="4">
                  Name
                </v-col>
                <v-col cols="8">
                  {{ order.recipient.name }}
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="4">
                  Country
                </v-col>
                <v-col cols="8">
                  {{ order.recipient.country }}
                </v-col>
              </v-row>
              <v-row style="background-color: rgba(255,255,255,.1);">
                <v-col cols="4">
                  Postal code
                </v-col>
                <v-col cols="8">
                  {{ order.recipient.postalCode }}
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="4">
                  City
                </v-col>
                <v-col cols="8">
                  {{ order.recipient.city }}
                </v-col>
              </v-row>
              <v-row style="background-color: rgba(255,255,255,.1);">
                <v-col cols="4">
                  Street
                </v-col>
                <v-col cols="8">
                  {{ order.recipient.street }}
                </v-col>
              </v-row>
            </v-container>
          </v-col>
        </v-row>
        <v-row v-if="success">
          <v-col>
            <v-alert type="success">
              Your order was placed.
            </v-alert>
          </v-col>
        </v-row>
        <v-row v-if="errorMsg != ''">
          <v-col>
            <v-alert type="error">
              {{ errorMsg }}
            </v-alert>
          </v-col>
        </v-row>
        <v-row>
          <v-col sm="4" offset-sm="8" class="text-right">
            <v-btn
              block
              large
              :disabled="success || errorMsg != ''"
              color="primary"
              @click="placeOrder"
            >
              Buy
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-flex>
  </v-layout>
</template>

<script>
import { mapState, mapMutations } from 'vuex'
import { OrdersApi, Configuration } from '~/client'

export default {
  data: () => ({
    quantity: [],
    order: null,
    products: {},
    success: false,
    errorMsg: ''
  }),
  computed: {
    ...mapState({
      user: 'user',
      orderFromState: 'order',
      productsFromState: 'products',
      cartEmpty: ({ cart }) => !cart.positions.length
    }),
    positions() {
      if (this.order == null) {
        return []
      }
      return this.order.positions.map((p) => ({
        ...p,
        product: p.product ||
          Object.values(this.products).find((pr) => pr.id === p.productId) || {
            name: 'unknown'
          }
      }))
    }
  },
  beforeMount() {
    if (!this.user.id) {
      this.$router.push('/login')
    }
    if (this.cartEmpty) {
      this.$router.push('/cart')
    }
  },
  mounted() {
    // init only once, don't sync after that
    this.order = this.orderFromState
    this.products = this.productsFromState
  },
  methods: {
    ...mapMutations(['orderPlaced', 'clearOrder']),
    placeOrder() {
      new OrdersApi(
        new Configuration({
          username: this.user.id,
          password: this.user.password
        })
      )
        .placeOrder(this.order.id)
        .then(({ data }) => data)
        .then((order) => this.orderPlaced(order))
        .then(() => (this.success = true))
        .catch(({ response: { status } }) => {
          switch (status) {
            case 410:
              this.errorMsg =
                'The order is not fully available anymore. Please retry ordering your cart.'
              this.clearOrder()
              break
          }
        })
    }
  }
}
</script>
