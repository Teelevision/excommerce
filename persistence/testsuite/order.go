package testsuite

import (
	"errors"
	"sync"
	"testing"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/stretchr/testify/suite"
)

// OrderRepositoryTestSuite is the suite that tests that a order repository
// behaves as expected. Use RunSuite to run it.
type OrderRepositoryTestSuite struct {
	suite.Suite
	NewRepository func() persistence.OrderRepository
}

// RunSuite runs the test suite.
func (s *OrderRepositoryTestSuite) RunSuite(t *testing.T) {
	suite.Run(t, s)
}

// TestCreateOrder tests the order creation.
func (s *OrderRepositoryTestSuite) TestCreateOrder() {
	s.Run("one", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx,
			"c1f8e321-4eb4-4f6f-9674-c9b1e452e7d9", // user id
			"ba3e44b1-59ea-4325-a8a8-600f3a081e73", // id
			persistence.OrderAttributes{
				Hash:   []byte("foo\nbar"),
				CartID: "2c3573ab-1d57-46bf-b979-5eaac02d850b",
				Buyer: persistence.OrderAddress{
					Name:       "Bundeskanzleramt, Bundeskanzlerin Angela Merkel",
					Country:    "DE",
					PostalCode: "10557",
					City:       "Berlin",
					Street:     "Willy-Brandt-Straße 1",
				},
				Recipient: persistence.OrderAddress{
					Name:       "Bundeskanzleramt, Bundeskanzlerin Angela Merkel",
					Country:    "DE",
					PostalCode: "10557",
					City:       "Berlin",
					Street:     "Willy-Brandt-Straße 1",
				},
				Coupons: []string{"orange30"},
			},
		)
		s.NoError(err)
	})
	s.Run("with empty everything", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "", "", persistence.OrderAttributes{})
		s.NoError(err)
	})
	s.Run("many", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			userID, id string
		}{
			{"f6c91ca9-174c-416f-9ca4-46c9e31417cc", "cb478a2e-b058-4601-8efd-4e5a98e6b67a"},
			{"577511d3-e27f-4e42-ba8d-277280eb1793", "a5538d2c-6a7b-4ca7-aa20-ecf2efd91e26"},
			{"fc6694e7-2822-4cb0-a254-e064b3aeb608", "2a9976ea-ea4a-40a2-b3eb-221dec1c03b9"},
			{"6300121d-ebdd-4a0f-9bea-298f031e1149", "6917a7c6-f84a-43b9-85ff-6b0511e173f1"},
			{"850ba076-9062-4d3f-9106-5fd43c94eea0", "10272a91-42d5-4ca2-894b-c101790c95b1"},
			{"376a44ca-5424-4b35-842d-5b2ad6aec15d", "1ea1d63e-4c6d-4224-b7d8-43a73046e59f"},
			{"96dd4f7a-782b-47e2-b9f2-74df37b97dc1", "7d31c8b4-ace6-43fa-a9a9-8290862544ed"},
		} {
			err := r.CreateOrder(ctx, c.userID, c.id, persistence.OrderAttributes{})
			s.Require().NoError(err)
		}
	})
	s.Run("conflict on same id", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user1", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.CreateOrder(ctx, "user2", "id", persistence.OrderAttributes{})
		s.True(errors.Is(err, persistence.ErrConflict))
	})
	s.Run("no conflict on same user", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id1", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.CreateOrder(ctx, "user", "id2", persistence.OrderAttributes{})
		s.NoError(err)
	})
	s.Run("is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.CreateOrder(ctx, "user", "ID", persistence.OrderAttributes{})
		s.NoError(err)
	})
	s.Run("supports more complex strings", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "औकखग \u0000\t\"abc", "‽ⓐ◐\n👽 乐乑", persistence.OrderAttributes{})
		s.NoError(err)
	})
	s.Run("does not trim whitespaces", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", " a", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.CreateOrder(ctx, "user", "a", persistence.OrderAttributes{})
		s.NoError(err)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"1ec6000d-5fa8-4e54-a242-7f28148f78b1",
			"8fc589db-e2c0-4ad8-87af-492c79b147ef",
			"1a0a5b6c-8083-4060-a3d9-e6290e27c3b4",
			"3121f5ec-8adc-4d33-b0dc-40f94366c573",
			"37f62f28-6ec6-4881-99df-1035826be65d",
			"9a6509d2-7485-4fdc-adbe-ba956df98188",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.CreateOrder(ctx, "user", id, persistence.OrderAttributes{})
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(ids[:3])
		go do(ids[3:])
		wg.Wait()
	})
	s.Run("changing the input does not have any side effects", func() {
		r := s.NewRepository()
		attributes := persistence.OrderAttributes{
			Hash:    []byte("foobar"),
			Coupons: []string{"orange30"},
		}
		err := r.CreateOrder(ctx, "user", "id", attributes)
		s.Require().NoError(err)
		// changing the input ...
		attributes.Hash[0] = 'c'
		attributes.Hash = append(attributes.Hash, 'a')
		attributes.Coupons[0] = "changed"
		attributes.Coupons = append(attributes.Coupons, "added")
		// ... does not have any side effects
		order, err := r.FindOrderOfUser(ctx, "user", "id")
		s.NoError(err)
		s.Equal(&model.Order{
			ID:      "id",
			Hash:    []byte("foobar"),
			Coupons: []*model.Coupon{{Code: "orange30"}},
		}, order)
	})
}

// TestFindOrderOfUser tests finding an order of a user.
func (s *OrderRepositoryTestSuite) TestFindOrderOfUser() {
	s.Run("finds an order", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx,
			"user",
			"id",
			persistence.OrderAttributes{
				Hash:   []byte("foo\nbar"),
				CartID: "8a74f13f-c1ad-4d36-96b6-39f1604e77df",
				Buyer: persistence.OrderAddress{
					Name:       "Bundeskanzleramt, Bundeskanzlerin Angela Merkel",
					Country:    "DE",
					PostalCode: "10557",
					City:       "Berlin",
					Street:     "Willy-Brandt-Straße 1",
				},
				Recipient: persistence.OrderAddress{
					Name:       "Bundeskanzleramt, Bundeskanzlerin Angela Merkel",
					Country:    "DE",
					PostalCode: "10557",
					City:       "Berlin",
					Street:     "Willy-Brandt-Straße 1",
				},
				Coupons: []string{"orange30"},
			},
		)
		s.Require().NoError(err)
		order, err := r.FindOrderOfUser(ctx, "user", "id")
		s.NoError(err)
		s.Equal(&model.Order{
			ID:     "id",
			Hash:   []byte("foo\nbar"),
			CartID: "8a74f13f-c1ad-4d36-96b6-39f1604e77df",
			Buyer: model.Address{
				Name:       "Bundeskanzleramt, Bundeskanzlerin Angela Merkel",
				Country:    "DE",
				PostalCode: "10557",
				City:       "Berlin",
				Street:     "Willy-Brandt-Straße 1",
			},
			Recipient: model.Address{
				Name:       "Bundeskanzleramt, Bundeskanzlerin Angela Merkel",
				Country:    "DE",
				PostalCode: "10557",
				City:       "Berlin",
				Street:     "Willy-Brandt-Straße 1",
			},
			Coupons: []*model.Coupon{{Code: "orange30"}},
		}, order)
	})
	s.Run("user is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		order, err := r.FindOrderOfUser(ctx, "USER", "id")
		s.True(errors.Is(err, persistence.ErrNotOwnedByUser))
		s.Nil(order)
	})
	s.Run("id is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		order, err := r.FindOrderOfUser(ctx, "user", "ID")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(order)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"b6fc0f36-742d-43f5-bed5-0b5759e18683",
			"3c4ba632-3a66-4e53-a0ca-bffcd8b48e9e",
			"09445c0b-dc0f-420e-bc2a-169ce3075c2b",
			"0beedae7-79f8-470d-a48b-33e30a5db1ab",
			"1cd9610f-e5b9-49b4-b59a-dedd2ded2749",
			"d89ada42-e763-410e-bb3a-cc1b112118b7",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.CreateOrder(ctx, "user", id, persistence.OrderAttributes{})
				s.Require().NoError(err)
			}
			for _, id := range ids {
				_, err := r.FindOrderOfUser(ctx, "user", id)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(ids[:3])
		go do(ids[3:])
		wg.Wait()
	})
	s.Run("changing the result does not have any side effects", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx,
			"user",
			"id",
			persistence.OrderAttributes{
				Hash:    []byte("foobar"),
				Coupons: []string{"orange30"},
			},
		)
		s.Require().NoError(err)
		order, err := r.FindOrderOfUser(ctx, "user", "id")
		s.Require().NoError(err)
		s.Require().Equal(&model.Order{
			ID:      "id",
			Hash:    []byte("foobar"),
			Coupons: []*model.Coupon{{Code: "orange30"}},
		}, order)
		// changing the result ...
		order.Hash[0] = 'c'
		order.Hash = append(order.Hash, 'a')
		order.Coupons[0].Code = "changed"
		order.Coupons = append(order.Coupons, &model.Coupon{Code: "added"})
		// ... does not have any side effects
		order, err = r.FindOrderOfUser(ctx, "user", "id")
		s.NoError(err)
		s.Equal(&model.Order{
			ID:      "id",
			Hash:    []byte("foobar"),
			Coupons: []*model.Coupon{{Code: "orange30"}},
		}, order)
	})
}

// TestDeleteOrderOfUser tests deleting an order of a user.
func (s *OrderRepositoryTestSuite) TestDeleteOrderOfUser() {
	s.Run("deletes an order", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.DeleteOrderOfUser(ctx, "user", "id")
		s.Require().NoError(err)
		s.Run("prevents deleting it again", func() {
			err := r.DeleteOrderOfUser(ctx, "user", "id")
			s.True(errors.Is(err, persistence.ErrDeleted))
		})
		s.Run("prevents accessing it", func() {
			_, err := r.FindOrderOfUser(ctx, "user", "id")
			s.True(errors.Is(err, persistence.ErrDeleted))
		})
		s.Run("prevents re-creating it", func() {
			err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
			s.True(errors.Is(err, persistence.ErrConflict))
		})
		s.Run("prevents locking it", func() {
			err := r.LockOrderOfUser(ctx, "user", "id")
			s.True(errors.Is(err, persistence.ErrDeleted))
		})
	})
	s.Run("user is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.DeleteOrderOfUser(ctx, "USER", "id")
		s.True(errors.Is(err, persistence.ErrNotOwnedByUser))
	})
	s.Run("id is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.DeleteOrderOfUser(ctx, "user", "ID")
		s.True(errors.Is(err, persistence.ErrNotFound))
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"91c2ec5a-0cfc-4737-8cd9-9fb0efaa5214",
			"7569a2f5-61e8-4fbd-b043-28d26357a7d4",
			"8ab95623-891b-4cd7-b32f-ce000b46cf67",
			"8c29ec61-ba63-47d0-95aa-235c6a60e9dc",
			"81be76ab-ca72-4322-88dd-c49a844159bd",
			"224c2fe7-8461-4284-a521-f682876a40b6",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.CreateOrder(ctx, "user", id, persistence.OrderAttributes{})
				s.Require().NoError(err)
			}
			for _, id := range ids {
				err := r.DeleteOrderOfUser(ctx, "user", id)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(ids[:3])
		go do(ids[3:])
		wg.Wait()
	})
}

// TestLockOrderOfUser tests locking an order of a user.
func (s *OrderRepositoryTestSuite) TestLockOrderOfUser() {
	s.Run("locks an order", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.LockOrderOfUser(ctx, "user", "id")
		s.Require().NoError(err)
		s.Run("prevents locking it again", func() {
			err := r.LockOrderOfUser(ctx, "user", "id")
			s.True(errors.Is(err, persistence.ErrLocked))
		})
		s.Run("prevents deleting it", func() {
			err := r.DeleteOrderOfUser(ctx, "user", "id")
			s.True(errors.Is(err, persistence.ErrLocked))
		})
		s.Run("and accessing it returns the locked state", func() {
			order, err := r.FindOrderOfUser(ctx, "user", "id")
			s.NoError(err)
			s.True(order.Locked)
		})
		s.Run("prevents re-creating it", func() {
			err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
			s.True(errors.Is(err, persistence.ErrConflict))
		})
	})
	s.Run("user is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.LockOrderOfUser(ctx, "USER", "id")
		s.True(errors.Is(err, persistence.ErrNotOwnedByUser))
	})
	s.Run("id is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateOrder(ctx, "user", "id", persistence.OrderAttributes{})
		s.Require().NoError(err)
		err = r.LockOrderOfUser(ctx, "user", "ID")
		s.True(errors.Is(err, persistence.ErrNotFound))
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"0b29ec74-54ff-4869-9b0e-79b592158eab",
			"cbfd4090-7698-46e5-ae18-6a4b51fb8975",
			"7cb3fc3c-ca73-404b-ac88-08f2de941217",
			"fc7f07dc-55a3-4d87-9e7c-58b5282053ec",
			"f3db961b-a048-432c-955e-779171e576c3",
			"bc227ca9-2a7f-48d6-82b1-57096d44e6a5",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.CreateOrder(ctx, "user", id, persistence.OrderAttributes{})
				s.Require().NoError(err)
			}
			for _, id := range ids {
				err := r.LockOrderOfUser(ctx, "user", id)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(ids[:3])
		go do(ids[3:])
		wg.Wait()
	})
}
