import { Position, Address } from './'

export class Order {
  id?: string
  price?: number
  status?: string
  buyer?: Address
  recipient?: Address
  positions?: Position[]
}
