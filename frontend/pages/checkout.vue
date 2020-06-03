<template>
  <v-layout>
    <v-flex class="text-center">
      <v-container>
        <v-row>
          <v-col cols="12" sm="6" style="text-align: left;">
            <div style="min-height: 66px;">
              <small>Billing address</small>
            </div>
            <v-text-field
              v-model="buyer.name"
              :error-messages="apiErrors['/buyer/name'] || ''"
              label="Name"
              required
            ></v-text-field>
            <v-select
              v-model="buyer.country"
              :items="[
                { value: 'DE', text: 'Germany' },
                { value: 'US', text: 'USA' }
              ]"
              :error-messages="apiErrors['/buyer/country'] || ''"
              label="Country"
            ></v-select>
            <v-text-field
              v-model="buyer.postalCode"
              :error-messages="apiErrors['/buyer/postalCode'] || ''"
              label="Postal code"
              required
            ></v-text-field>
            <v-text-field
              v-model="buyer.city"
              :error-messages="apiErrors['/buyer/city'] || ''"
              label="City"
              required
            ></v-text-field>
            <v-text-field
              v-model="buyer.street"
              :error-messages="apiErrors['/buyer/street'] || ''"
              label="Street"
              required
            ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" style="text-align: left;">
            <v-switch v-model="differentRecipient">
              <template v-slot:label>
                <small>Different delivery address</small>
              </template>
            </v-switch>
            <v-text-field
              v-if="differentRecipient"
              v-model="recipient.name"
              :error-messages="apiErrors['/recipient/name'] || ''"
              label="Name"
              required
            ></v-text-field>
            <v-select
              v-if="differentRecipient"
              v-model="recipient.country"
              :items="[
                { value: 'DE', text: 'Germany' },
                { value: 'US', text: 'USA' }
              ]"
              :error-messages="apiErrors['/recipient/country'] || ''"
              label="Country"
            ></v-select>
            <v-text-field
              v-if="differentRecipient"
              v-model="recipient.postalCode"
              :error-messages="apiErrors['/recipient/postalCode'] || ''"
              label="Postal code"
              required
            ></v-text-field>
            <v-text-field
              v-if="differentRecipient"
              v-model="recipient.city"
              :error-messages="apiErrors['/recipient/city'] || ''"
              label="City"
              required
            ></v-text-field>
            <v-text-field
              v-if="differentRecipient"
              v-model="recipient.street"
              :error-messages="apiErrors['/recipient/street'] || ''"
              label="Street"
              required
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-combobox
              v-model="coupons"
              label="Coupons"
              multiple
              small-chips
              append-icon=""
              deletable-chips
              :delimiters="[' ', ',']"
              :error-messages="
                Object.entries(apiErrors)
                  .filter(([key]) => key.startsWith('/coupons'))
                  .map(([key, value]) => value) || []
              "
            ></v-combobox>
          </v-col>
        </v-row>
        <v-row>
          <v-col sm="4" offset-sm="8" class="text-right">
            <v-btn block large color="primary" @click="storeOrder">
              Proceed
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
    buyer: {
      name: '',
      country: '',
      postalCode: '',
      city: '',
      street: ''
    },
    differentRecipient: false,
    recipient: {
      name: '',
      country: '',
      postalCode: '',
      city: '',
      street: ''
    },
    coupons: [],
    apiErrors: {}
  }),
  computed: mapState({
    user: 'user',
    cartId: ({ cart }) => cart.id,
    cartEmpty: ({ cart }) => !cart.positions.length
  }),
  beforeMount() {
    if (!this.user.id) {
      this.$router.push('/login')
    }
    if (this.cartEmpty) {
      this.$router.push('/cart')
    }
  },
  methods: {
    ...mapMutations(['orderReceived']),
    storeOrder() {
      const buyer = { ...this.buyer }
      const recipient = this.differentRecipient ? { ...this.recipient } : buyer
      const coupons = [...this.coupons]

      this.apiErrors = {}

      new OrdersApi(
        new Configuration({
          username: this.user.id,
          password: this.user.password
        })
      )
        .createOrderFromCart(this.cartId, { buyer, recipient, coupons })
        .then(({ data }) => data)
        .then((order) => this.orderReceived(order))
        .then(() => this.$router.push('/place-order'))
        .catch(({ response: { status, data } }) => {
          switch (status) {
            case 422:
              this.apiErrors = { [data.pointer]: data.message }
              break
          }
        })
    }
  }
}
</script>
