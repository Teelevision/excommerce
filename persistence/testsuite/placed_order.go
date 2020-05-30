package testsuite

import (
	"sync"
	"testing"

	"github.com/Teelevision/excommerce/persistence"
	"github.com/stretchr/testify/suite"
)

// PlacedOrderRepositoryTestSuite is the suite that tests that a placed order
// repository behaves as expected. Use RunSuite to run it.
type PlacedOrderRepositoryTestSuite struct {
	suite.Suite
	NewRepository func() persistence.PlacedOrderRepository
}

// RunSuite runs the test suite.
func (s *PlacedOrderRepositoryTestSuite) RunSuite(t *testing.T) {
	suite.Run(t, s)
}

// TestPlaceOrder tests placing orders.
func (s *PlacedOrderRepositoryTestSuite) TestPlaceOrder() {
	s.Run("one", func() {
		r := s.NewRepository()
		err := r.PlaceOrder(ctx, persistence.PlacedOrder{
			UserID: "8e668ea2-ba30-421b-a773-6e289b5b68fd",
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
			Coupons: map[string]persistence.OrderCoupon{
				"orange30": {
					ProductID: "5b31a473-4b5e-48ad-8033-bcccdfb373f9",
					Name:      "30% off oranges",
					Discount:  30,
				},
			},
			Products: map[string]persistence.OrderProduct{
				"5b31a473-4b5e-48ad-8033-bcccdfb373f9": {
					Name:  "Orange",
					Price: 79,
				},
				"a67d84d3-3417-478f-b93f-fb5990ce0052": {
					Name:  "Apple",
					Price: 49,
				},
			},
			Price: 160,
			Positions: []persistence.OrderPosition{
				{
					ProductID:  "5b31a473-4b5e-48ad-8033-bcccdfb373f9",
					CouponCode: "",
					Quantity:   2,
					Price:      158,
				}, {
					ProductID:  "",
					CouponCode: "orange30",
					Quantity:   1,
					Price:      -47,
				}, {
					ProductID:  "a67d84d3-3417-478f-b93f-fb5990ce0052",
					CouponCode: "",
					Quantity:   1,
					Price:      49,
				},
			},
		})
		s.NoError(err)
	})
	s.Run("with empty everything", func() {
		r := s.NewRepository()
		err := r.PlaceOrder(ctx, persistence.PlacedOrder{})
		s.NoError(err)
	})
	s.Run("many", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			userID string
		}{
			{"50d710f1-e8e6-4449-a690-57585e970699"},
			{"b00a387d-dcac-4ecf-9ad0-88857d9bc6ee"},
			{"249ed004-4c5e-47f4-bf55-9ceea0c6a385"},
			{"d062616b-52eb-48c3-9372-7876f7eaff53"},
			{"2e0bd181-b76e-46a2-a538-bb8f0da7b6c0"},
			{"ca77448b-3eca-445a-b3c4-28c66d132879"},
			{"f1af4738-4b7a-492f-ac09-9efe8da20731"},
		} {
			err := r.PlaceOrder(ctx, persistence.PlacedOrder{
				UserID: c.userID,
			})
			s.Require().NoError(err)
		}
	})
	s.Run("supports more complex strings", func() {
		r := s.NewRepository()
		err := r.PlaceOrder(ctx, persistence.PlacedOrder{
			UserID: "‽ⓐ◐\n👽",
			Buyer: persistence.OrderAddress{
				Name:       "औकखग",
				Country:    "ÖÄÜß",
				PostalCode: "乐乑 ",
				City:       "‽ⓐ◐\n👽",
				Street:     " \u0000\t\"abc",
			},
			Recipient: persistence.OrderAddress{
				Name:       "乐乑",
				Country:    "DE",
				PostalCode: "\u0000\t\"abc ",
				City:       "‽ⓐ◐\n👽",
				Street:     " औकखग",
			},
		})
		s.NoError(err)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"5bd52b3a-ed59-460f-b6fb-d0c0d49f81e0",
			"196bcfaa-154b-4c7e-bf35-907d96ce1f15",
			"721e6fdc-605a-4528-9d96-2418747cdf0e",
			"1386a8bd-2663-4ed6-b36e-47553c89e036",
			"6a0152cf-46fb-4606-a02b-b4bbef8421f3",
			"9449ff31-949f-43f0-88a1-d2051c5479e4",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.PlaceOrder(ctx, persistence.PlacedOrder{
					UserID: id,
				})
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(ids[:3])
		go do(ids[3:])
		wg.Wait()
	})
}
