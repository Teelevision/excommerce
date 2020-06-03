import { Store, ActionTree, ActionContext } from 'vuex'
import createPersistedState from 'vuex-persistedstate'
import { v4 as uuidv4 } from 'uuid'
import { Product, Cart, Position, User, Order } from '~/models'
import { ProductsApi, UsersApi, CartsApi, Configuration } from '~/client'

interface State {
  products: Product[]
  cart: Cart
  user: User
  order: Order | null
}

export const state: () => State = () => ({
  products: [],
  cart: new Cart(),
  user: new User('', ''),
  order: null
})

const initialState = state

export const mutations = {
  allProductsLoaded(state: State, products: Product[]) {
    state.products = products
  },
  updateCart(state: State, cart: Cart) {
    state.cart = cart
  },
  updateCartPositions(state: State, positions: Position[]) {
    const products: { [key: string]: Position } = {}
    for (const i in positions) {
      const { quantity, savedPrice, product: p } = positions[i]
      let { productId } = positions[i]
      if (productId === undefined) {
        productId = `x-${i}`
      }
      let product = products[productId]
      if (product === undefined) {
        product = { quantity: 0, savedPrice, productId, product: p }
      }
      product.quantity += quantity
      product.price = 0
      products[productId] = product
    }
    positions = []
    for (let product of Object.values(products)) {
      if (product.quantity <= 0) continue
      if (product.productId?.startsWith('x-')) {
        product = { ...product, productId: undefined }
      }
      positions.push(product)
    }
    state.cart.positions = positions
  },
  addProductToCart(state: State, productId: string) {
    mutations.updateCartPositions(state, [
      ...state.cart.positions,
      { productId, quantity: 1 }
    ])
  },
  loggedIn(state: State, user: User) {
    state.user = user
  },
  loggedOut(state: State) {
    state.user = initialState().user
    state.cart = initialState().cart
  },
  orderReceived(state: State, order: Order) {
    state.order = order
  },
  clearOrder(state: State) {
    state.order = initialState().order
  },
  orderPlaced(state: State, _order: Order) {
    mutations.clearOrder(state)
    state.cart = initialState().cart
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
  updateCartPositions({ commit, dispatch }, positions: Position[]) {
    commit('updateCartPositions', positions)
    dispatch('storeCartOnServer')
  },
  addToCart({ commit, dispatch }, productId: string) {
    commit('addProductToCart', productId)
    dispatch('storeCartOnServer')
  },
  async login({ commit, dispatch }, user: User) {
    const resp = await new UsersApi().login(user)
    commit('loggedIn', { ...user, id: resp.data.id })
    dispatch('syncCart')
  },
  logout({ commit }) {
    commit('loggedOut')
  },
  syncCart({ dispatch, state }) {
    if (state.cart.positions.length) {
      dispatch('storeCartOnServer')
    } else {
      dispatch('loadCartFromServer')
    }
  },
  async storeCartOnServer({ commit, dispatch, state: { cart, user } }) {
    if (!user.id) {
      return
    }
    const uuid = cart.id || uuidv4()
    try {
      const {
        data: { id, positions }
      } = await new CartsApi(
        new Configuration({ username: user.id, password: user.password })
      ).storeCart(uuid, {
        id: uuid,
        positions: cart.positions
          .filter(({ productId }) => productId !== undefined)
          .map(({ quantity, productId }) => ({
            quantity,
            product: { id: productId || '', name: '' },
            price: 0,
            savedPrice: 0
          }))
      })
      commit('updateCart', {
        id,
        positions: positions.map(
          ({ quantity, product, price, savedPrice }) => ({
            quantity,
            product,
            price,
            savedPrice,
            productId: product.id
          })
        )
      })
    } catch (e) {
      switch (e.response.status) {
        case 423:
        case 410:
          // try again with new cart id
          if (cart.id) {
            commit('updateCart', { ...cart, id: undefined })
            dispatch('storeCartOnServer')
            return
          }
      }
      throw e
    }
  },
  async loadCartFromServer({ commit, state: { user } }) {
    const { data: carts } = await new CartsApi(
      new Configuration({ username: user.id, password: user.password })
    ).getAllCarts(false)
    for (const cart of carts) {
      if (!cart.positions.length) {
        continue
      }
      commit('updateCart', {
        id: cart.id,
        positions: cart.positions.map(
          ({ quantity, product, price, savedPrice }) => ({
            quantity,
            product,
            price,
            savedPrice,
            productId: product.id
          })
        )
      })
      break
    }
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

// plugins

export const plugins = [
  createPersistedState({
    paths: ['user'],
    storage: sessionStorage
  }),
  createPersistedState({
    paths: ['cart']
  })
]
