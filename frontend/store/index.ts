import { Store, ActionTree, ActionContext } from 'vuex'
import Product from '~/models/product'
import { ProductsApi } from '~/client'

interface State {
  products: Product[]
}

export const state: () => State = () => ({
  products: []
})

export const mutations = {
  allProductsLoaded(state: State, products: Product[]) {
    state.products = products
  }
}

export const actions = <ActionTreeMutations>{
  async loadAllProducts({ commit }) {
    const api = new ProductsApi()
    commit(
      'allProductsLoaded',
      await api
        .getAllProducts()
        .then((resp) => resp.data)
        .then((products) =>
          products.map(
            (p) => <Product>{ ID: p.id, Name: p.name, Price: p.price }
          )
        )
    )
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
