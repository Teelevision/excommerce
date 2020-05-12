import { Store, ActionTree, ActionContext } from 'vuex'
import { Product, Cart, Position } from '~/models'
import { ProductsApi } from '~/client'

interface State {
  products: Product[]
  cart: Cart
}

export const state: () => State = () => ({
  products: [],
  cart: new Cart()
})

export const mutations = {
  allProductsLoaded(state: State, products: Product[]) {
    state.products = products
  },
  updateCartPositions(state: State, positions: Position[]) {
    const products: { [key: string]: Position } = {}
    for (const { productId, quantity } of positions) {
      if (productId === undefined) {
        continue
      }
      let product = products[productId]
      if (product === undefined) {
        product = { quantity: 0, productId }
      }
      product.quantity += quantity
      product.price = 0
      products[productId] = product
    }
    positions = []
    for (const product of Object.values(products)) {
      if (product.quantity <= 0) continue
      positions.push(product)
    }
    state.cart.positions = positions
  },
  addProductToCart(state: State, productId: string) {
    mutations.updateCartPositions(state, [
      ...state.cart.positions,
      { productId, quantity: 1 }
    ])
  }
}

export const actions = <ActionTreeMutations>{
  async loadAllProducts({ commit }) {
    const api = new ProductsApi()
    commit(
      'allProductsLoaded',
      await api.getAllProducts().then((resp) => resp.data)
    )
  },
  updateCartPositions({ commit }, positions: Position[]) {
    commit('updateCartPositions', positions)
  },
  addToCart({ commit }, productId: string) {
    commit('addProductToCart', productId)
  }
}

// The code below is so that the commit function of the actions has type
// hinting.

interface ActionTreeMutations extends ActionTree<State, State> {
  [key: string]: ActionHandler
}

type ActionHandler = (
  this: Store<State>,
  injectee: ActionContextMutation,
  payload?: any
) => any

interface ActionContextMutation extends ActionContext<State, State> {
  commit: CommitMutation
}

interface CommitMutation {
  <K extends keyof typeof mutations>(
    type: K,
    payload?: typeof mutations[K] extends (
      state: State,
      payload: infer U
    ) => void
      ? U
      : never,
    options?: any
  ): void
}
