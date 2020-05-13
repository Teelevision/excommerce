import { Product } from '.'

export class Position {
  productId?: string
  product?: Product
  quantity!: number
  price?: number
}
