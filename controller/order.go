package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Teelevision/excommerce/authentication"
	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/google/uuid"
)

// Order is the controller that handles orders.
type Order struct {
	OrderRepository       persistence.OrderRepository
	CartRepository        persistence.CartRepository
	ProductRepository     persistence.ProductRepository
	CouponRepository      persistence.CouponRepository
	PlacedOrderRepository persistence.PlacedOrderRepository
}

// CreateAndGet creates the given order. The order is returned with a unique id.
func (c *Order) CreateAndGet(ctx context.Context, order *model.Order) (*model.Order, error) {
	// create id
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	id := uuid.String()

	// prepare positions
	positions := generateOrderPositions(order.Cart.Positions, order.Coupons)

	// hash
	hash := hashPositions(positions)

	// coupon codes
	couponCodes := make([]string, len(order.Coupons))
	for i, coupon := range order.Coupons {
		couponCodes[i] = coupon.Code
	}

	err = c.OrderRepository.CreateOrder(ctx,
		authentication.AuthenticatedUser(ctx).ID,
		id,
		persistence.OrderAttributes{
			Hash:      hash,
			CartID:    order.CartID,
			Buyer:     persistence.OrderAddress(order.Buyer),
			Recipient: persistence.OrderAddress(order.Recipient),
			Coupons:   couponCodes,
		},
	)
	switch {
	case errors.Is(err, persistence.ErrConflict):
		panic(err)
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, err
	case err == nil:
		return &model.Order{
			ID:        id,
			Hash:      hash,
			Buyer:     order.Buyer,
			Recipient: order.Recipient,
			Cart:      order.Cart,
			CartID:    order.CartID,
			Coupons:   order.Coupons,
			Positions: positions,
			Price:     calculatePositionSum(positions),
		}, nil
	default:
		panic(err)
	}
}

// Place places the order with the given id. ErrNotFound is returned if the
// order does not exist. ErrDeleted is returned if the order used to exist, but
// is deleted. ErrForbidden is returned if the order exists, but is not owned by
// the current user. ErrLocked is returned if the order is already placed.
func (c *Order) Place(ctx context.Context, orderID string) (*model.Order, error) {
	// First call checks and locks the order and cart. This ensures that the
	// order did not change and the cart cannot be updated anymore.
	_, err := c.preparePlace(ctx, orderID, false)
	if err != nil {
		return nil, err
	}

	// Second call checks order and cart again after they were locked. This
	// ensures that between checking and locking nothing happened that changed
	// the order, like a product became unavailable.
	order, err := c.preparePlace(ctx, orderID, true)
	if err != nil {
		return nil, err
	}

	// Place order. Just locking it is not enough, because there is the edge
	// case that a locked order did change between checking and locking it,
	// which is why we have the second call above. Placing the order also saves
	// a the current state of the cart and products, which we all got from the
	// second call.
	placedOrder := persistence.PlacedOrder{
		UserID:    authentication.AuthenticatedUser(ctx).ID,
		Buyer:     persistence.OrderAddress(order.Buyer),
		Recipient: persistence.OrderAddress(order.Recipient),
		Coupons:   make(map[string]persistence.OrderCoupon, len(order.Coupons)),
		Products:  make(map[string]persistence.OrderProduct),
		Price:     order.Price,
		Positions: make([]persistence.OrderPosition, len(order.Positions)),
	}
	for _, coupon := range order.Coupons {
		placedOrder.Coupons[coupon.Code] = persistence.OrderCoupon{
			ProductID: coupon.ProductID,
			Name:      coupon.Name,
			Discount:  coupon.Discount,
		}
	}
	for i, position := range order.Positions {
		placedOrder.Positions[i] = persistence.OrderPosition{
			ProductID:  position.ProductID,
			CouponCode: position.CouponCode,
			Quantity:   position.Quantity,
			Price:      position.Price,
		}
		if position.ProductID != "" {
			placedOrder.Products[position.ProductID] = persistence.OrderProduct{
				Name:  position.Product.Name,
				Price: position.Product.Price,
			}
		}
	}
	err = c.PlacedOrderRepository.PlaceOrder(ctx, placedOrder)
	switch {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, err
	case err == nil:
		return order, nil
	default:
		panic(err)
	}
}

// Must be called twice. First time it expects the order and cart to be
// unlocked. It locks and returns both, including products. Second time it
// expects the order and cart to be locked. Both times it checks that the order
// did not change in any way, like containing a product that changed its price.
// Both calls together ensure that the resources were locked and that we were
// the caller that locked them. Worst case we do not place the order but lock it
// and the cart. Not placing the order is correct, because something changed.
// Locking the cart just means that the client has to store it under a new id.
// But this is an edge case. Notice that in that case we would return ErrLocked
// instead of ErrDeleted, because we cannot delete the order that we just
// locked.
func (c *Order) preparePlace(ctx context.Context, orderID string, expectLocked bool) (*model.Order, error) {
	userID := authentication.AuthenticatedUser(ctx).ID

	// load order
	order, err := c.OrderRepository.FindOrderOfUser(ctx, userID, orderID)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, ErrNotFound
	case errors.Is(err, persistence.ErrDeleted):
		return nil, ErrDeleted
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return nil, ErrForbidden
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, err
	case err == nil:
		if !expectLocked && order.Locked {
			return nil, ErrLocked
		}
	default:
		panic(err)
	}

	// if the cart changed somehow, we delete the outdated order
	deleteOrder := func() error {
		if err := c.Delete(ctx, orderID); err != nil {
			return err
		}
		return ErrDeleted
	}

	// load cart
	order.Cart, err = c.CartRepository.FindCartOfUser(ctx, userID, order.CartID)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, deleteOrder()
	case errors.Is(err, persistence.ErrDeleted):
		return nil, deleteOrder()
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return nil, deleteOrder()
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, err
	case err == nil:
		if !expectLocked && order.Cart.Locked {
			return nil, deleteOrder()
		}
	default:
		panic(err)
	}

	// load products
	for i, position := range order.Cart.Positions {
		if product := getSpecialProduct(position.ProductID); product != nil {
			order.Cart.Positions[i].Product = product
			continue
		}
		product, err := c.ProductRepository.FindProduct(ctx, position.ProductID)
		switch {
		case errors.Is(err, persistence.ErrNotFound):
			return nil, deleteOrder()
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return nil, err
		case err == nil:
			order.Cart.Positions[i].Product = product
		default:
			panic(err)
		}
	}

	// load coupons
	for i, coupon := range order.Coupons {
		coupon, err := c.CouponRepository.FindValidCoupon(ctx, coupon.Code)
		switch {
		case errors.Is(err, persistence.ErrNotFound):
			return nil, deleteOrder()
		case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
			return nil, err
		case err == nil:
			order.Coupons[i] = coupon
		default:
			panic(err)
		}
	}

	// prepare positions
	positions := generateOrderPositions(order.Cart.Positions, order.Coupons)

	// hash
	hash := hashPositions(positions)
	if !bytes.Equal(hash, order.Hash) {
		// The hash changed. This means that the cart changed, maybe indirectly,
		// like a product that changed its price.
		return nil, deleteOrder()
	}

	// second call ends here
	if expectLocked {
		return order, nil
	}

	// Lock cart first because we might need to delete the order. Also locking
	// the cart prevents the race condition that there are two orders on the
	// cart which are locked simultaneously. Only one caller can lock the cart.
	err = c.CartRepository.LockCartOfUser(ctx, userID, order.CartID)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, deleteOrder()
	case errors.Is(err, persistence.ErrDeleted):
		return nil, deleteOrder()
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return nil, deleteOrder()
	case errors.Is(err, persistence.ErrLocked):
		return nil, deleteOrder()
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, err
	case err == nil:
		// success
	default:
		panic(err)
	}

	// lock order
	err = c.OrderRepository.LockOrderOfUser(ctx, userID, orderID)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, ErrNotFound
	case errors.Is(err, persistence.ErrDeleted):
		return nil, ErrDeleted
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return nil, ErrForbidden
	case errors.Is(err, persistence.ErrLocked):
		return nil, ErrLocked
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return nil, err
	case err == nil:
		// success
	default:
		panic(err)
	}

	return order, nil
}

// Delete deletes the order with the given id. ErrNotFound is returned if the
// order does not exist. ErrDeleted is returned if the order is already deleted.
// ErrForbidden is returned if the order exists, but is not owned by the current
// user. ErrLocked is returned if the order is already placed and therefore
// cannot be deleted.
func (c *Order) Delete(ctx context.Context, orderID string) error {
	userID := authentication.AuthenticatedUser(ctx).ID
	err := c.OrderRepository.DeleteOrderOfUser(ctx, userID, orderID)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return ErrNotFound
	case errors.Is(err, persistence.ErrDeleted):
		return ErrDeleted
	case errors.Is(err, persistence.ErrNotOwnedByUser):
		return ErrForbidden
	case errors.Is(err, persistence.ErrLocked):
		return ErrLocked
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return err
	case err == nil:
		return nil
	default:
		panic(err)
	}
}

func generateOrderPositions(positions []model.Position, coupons []*model.Coupon) []model.Position {
	positions = consolidatePositions(positions)
	positions = calculatePositionPrices(positions)

	// drop all non-products
	for i := 0; i < len(positions); i++ {
		if positions[i].ProductID != "" {
			continue
		}
		positions = append(positions[:i], positions[i+1:]...)
		i--
	}

	// get best coupons
	productCoupons := make(map[string]*model.Coupon, len(coupons))
	for _, coupon := range coupons {
		existingCoupon, ok := productCoupons[coupon.ProductID]
		if !ok || coupon.Discount > existingCoupon.Discount {
			productCoupons[coupon.ProductID] = coupon
		}
	}
	for i := 0; i < len(positions); i++ {
		coupon, ok := productCoupons[positions[i].ProductID]
		if !ok {
			continue
		}
		// add position for coupon
		price := -coupon.Discount * positions[i].Price / 100
		couponPosition := model.Position{
			Quantity:   1,
			Price:      price,
			Coupon:     coupon,
			CouponCode: coupon.Code,
			SavedPrice: -price,
		}
		positions = append(positions[:i+1],
			append([]model.Position{couponPosition}, positions[i+1:]...)...)
		i++
	}

	// non-coupon discounts

	// If >=7 apples, 10% discount on all apples.
	for i, position := range positions {
		if position.ProductID != "a6da78f8-2be6-49ff-b40a-32aa86a6a986" {
			continue
		}
		if position.Quantity < 7 {
			continue
		}
		// insert positions for discount
		price := -10 * position.Price / 100
		discountPosition := model.Position{
			Quantity:   1,
			Price:      price,
			SavedPrice: -price,
			Product: &model.Product{ // no id
				Name:       "10% off apples",
				Price:      price,
				SavedPrice: -price,
			},
		}
		positions = append(positions[:i+1],
			append([]model.Position{discountPosition}, positions[i+1:]...)...)
		break
	}

	// Sets of 4 pears and 2 bananas get a 30% discount.
	pearPosition, bananaPosition, setPosition := -1, -1, -1
	for i, position := range positions {
		switch position.ProductID {
		case "5438bfe8-6bd2-4a88-ac36-ec29716eb6d7": // pear
			pearPosition = i
		case "b16088e1-9603-4676-a8df-130823cf15a5": // banana
			bananaPosition = i
		case "0de17a66-ea59-4032-9383-2603c6c77d25": // set
			setPosition = i
		}
	}
	var numSets int
	if pearPosition >= 0 && bananaPosition >= 0 {
		numSets = positions[pearPosition].Quantity / 4
		if n := positions[bananaPosition].Quantity / 2; n < numSets {
			numSets = n
		}
	}
	if numSets > 0 {
		// reduce pears and bananas
		positions[pearPosition].Quantity -= numSets * 4
		positions[pearPosition].Price = positions[pearPosition].Quantity * positions[pearPosition].Product.Price
		if positions[pearPosition].Quantity == 0 {
			positions = append(positions[:pearPosition], positions[pearPosition+1:]...)
			if pearPosition < bananaPosition {
				bananaPosition--
			}
		}
		positions[bananaPosition].Quantity -= numSets * 2
		positions[bananaPosition].Price = positions[bananaPosition].Quantity * positions[bananaPosition].Product.Price
		if positions[bananaPosition].Quantity == 0 {
			positions = append(positions[:bananaPosition], positions[bananaPosition+1:]...)
		}
		// update set quantity or add position for sets
		set := getSpecialProduct("0de17a66-ea59-4032-9383-2603c6c77d25")
		if setPosition >= 0 {
			positions[setPosition].Quantity += numSets
			positions[setPosition].Price += numSets * set.Price
			positions[setPosition].SavedPrice += numSets * set.SavedPrice
		} else {
			positions = append(positions, model.Position{
				Quantity:   numSets,
				Price:      numSets * set.Price,
				SavedPrice: numSets * set.SavedPrice,
				ProductID:  "0de17a66-ea59-4032-9383-2603c6c77d25",
				Product:    set,
			})
		}
	}

	return positions
}

func calculatePositionSum(positions []model.Position) (sum int) {
	for _, position := range positions {
		sum += position.Price
	}
	return
}

func hashPositions(positions []model.Position) []byte {
	entries := make(sort.StringSlice, len(positions))
	for i, position := range positions {
		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "%d,%d,", position.Quantity, position.Price)
		switch {
		case position.Product != nil:
			fmt.Fprintf(buf, "product:%s", position.Product.ID)
		case position.Coupon != nil:
			fmt.Fprintf(buf, "coupon:%s,%d,%q", position.Coupon.ProductID, position.Coupon.Discount, position.Coupon.Code)
		default:
			panic("position has no product and no coupon")
		}
		entries[i] = buf.String()
	}
	entries.Sort()
	base := strings.Join(entries, "\n")
	// TODO: choose hash algorithm
	return []byte(base)
}
