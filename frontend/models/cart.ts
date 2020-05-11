import { Position } from './'

export class Cart {
  id?: string
  positions: Position[]

  constructor() {
    this.positions = []
  }
}
