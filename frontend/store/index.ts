import { Store, ActionTree, ActionContext } from 'vuex'
import { v4 as uuidv4 } from 'uuid'
import { Product, Cart, Position, User } from '~/models'
import { ProductsApi, UsersApi, CartsApi, Configuration } from '~/client'

interface State {
  products: Product[]
  cart: Cart
  user: User
}

export const state: () => State = () => ({
  products: [],
  cart: new Cart(),
  user: new User('', '')
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
  },
  loggedIn(state: State, user: User) {
    state.user = user
  },
  loggedOut(state: State) {
    state.user = initialState().user
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
    dispatch('storeCartOnServer')
  },
  logout({ commit }) {
    commit('loggedOut')
  },
  async storeCartOnServer({ commit, state: { cart, user } }) {
    if (!user.id) {
      return
    }
    const uuid = cart.id || uuidv4()
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
          price: 0
        }))
    })
    commit('updateCart', {
      id,
      positions: positions.map(({ quantity, product, price }) => ({
        quantity,
        product,
        price,
        productId: product.id
      }))
    })
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
