package testsuite

import (
	"errors"
	"sync"
	"testing"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/stretchr/testify/suite"
)

// CartRepositoryTestSuite is the suite that tests that a cart repository
// behaves as expected. Use RunSuite to run it.
type CartRepositoryTestSuite struct {
	suite.Suite
	NewRepository func() persistence.CartRepository
}

// RunSuite runs the test suite.
func (s *CartRepositoryTestSuite) RunSuite(t *testing.T) {
	suite.Run(t, s)
}

// TestCreateCart tests the cart creation.
func (s *CartRepositoryTestSuite) TestCreateCart() {
	s.Run("one", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx,
			"0cfb6682-4279-4928-af37-30fd3e4c0b15", // user id
			"71388209-5d1d-4ce7-a7ac-38a2e75fd67c", // id
			[]struct {
				ProductID string
				Quantity  int
				Price     int
			}{
				{"ce21148d-f8c8-437a-bd9c-fd72797803dd", 1, 999},
				{"c92a017b-3e75-43d0-bb38-7f05e7d9b3c3", 999, -1},
			}, // products
		)
		s.NoError(err)
	})
	s.Run("with empty id, positions and locked status", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "", "", nil)
		s.NoError(err)
	})
	s.Run("many", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			userID, id string
		}{
			{"9c7cd4e5-2acb-4355-b0c7-38b0663e2143", "54a0b76f-df5f-427a-94ea-0b143a6b133a"},
			{"1cdd6ebd-0be4-4f13-829e-633d1f2f9ccd", "b8e3d979-7e58-4f70-bd27-ee298e995d34"},
			{"fbc00df4-2041-4afb-84b3-468f4af925f6", "31ae3114-db9f-4f84-9541-aa6729c5624b"},
			{"66825e5b-9db8-49b3-916b-6065731a3618", "640ab425-24c8-49a4-b4f7-8fca10313513"},
			{"b95f29ed-ff2c-40e7-bad8-88e3c69103b3", "350081f0-5c68-44cb-a177-24b9a8522b56"},
			{"053ac78e-4332-4e0f-bb1a-544c3b8bfd83", "2e485bac-dade-4a20-9fee-f5a5bcfe8428"},
			{"1e1d6a41-5827-40a6-9a82-404f278cce5c", "8b7d133d-666b-41c9-8d73-1fd391161b15"},
		} {
			err := r.CreateCart(ctx, c.userID, c.id, nil)
			s.Require().NoError(err)
		}
	})
	s.Run("conflict on same id", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user1", "id", nil)
		s.Require().NoError(err)
		err = r.CreateCart(ctx, "user2", "id", nil)
		s.True(errors.Is(err, persistence.ErrConflict))
	})
	s.Run("no conflict on same positions", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", []struct {
			ProductID string
			Quantity  int
			Price     int
		}{
			{"c7344f02-d3f6-407d-9b0d-9009eb16fcf2", 1, 100},
			{"c7344f02-d3f6-407d-9b0d-9009eb16fcf2", 1, 100},
		})
		s.Require().NoError(err)
	})
	s.Run("no conflict on same user", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id1", nil)
		s.Require().NoError(err)
		err = r.CreateCart(ctx, "user", "id2", nil)
		s.NoError(err)
	})
	s.Run("is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", nil)
		s.Require().NoError(err)
		err = r.CreateCart(ctx, "user", "ID", nil)
		s.NoError(err)
	})
	s.Run("supports more complex strings", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "औकखग \u0000\t\"abc", "‽ⓐ◐\n👽 乐乑", nil)
		s.NoError(err)
	})
	s.Run("does not trim whitespaces", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", " a", nil)
		s.Require().NoError(err)
		err = r.CreateCart(ctx, "user", "a", nil)
		s.NoError(err)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"37e25c1e-b132-4670-98c3-8d1e7e8f4e85",
			"9599f202-8e08-4933-8be7-2763b77f5cca",
			"c6352b8a-683d-41e7-b390-038673072d79",
			"c423009d-cbda-4664-acbf-266e897fd0d1",
			"d7b6aa1d-810f-4dd4-ba77-31ef88343a31",
			"2abe37af-cc8f-41e8-913c-0a1710298ab2",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.CreateCart(ctx, "user", id, nil)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(ids[:3])
		go do(ids[3:])
		wg.Wait()
	})
}

// TestUpdateCartOfUser tests updating carts.
func (s *CartRepositoryTestSuite) TestUpdateCartOfUser() {
	s.Run("updates a cart", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx,
			"user",
			"id",
			[]struct {
				ProductID string
				Quantity  int
				Price     int
			}{
				{"bb364bf5-e1fb-445d-be2d-ebad49316e0c", 1, 999},
				{"77181602-1b4c-463c-b9a5-c2188610fd68", 999, -1},
			}, // products
		)
		s.Require().NoError(err)
		err = r.UpdateCartOfUser(ctx,
			"user",
			"id",
			[]struct {
				ProductID string
				Quantity  int
				Price     int
			}{
				{"e11c7885-92a9-4833-8e52-ed020fef5aff", 2, 234},
				{"ff62397c-cbcf-4cd9-b57d-0a9348dd8ef4", 3, 345},
			}, // products
		)
		s.Require().NoError(err)
	})
	s.Run("user is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", nil)
		s.Require().NoError(err)
		err = r.UpdateCartOfUser(ctx, "USER", "id", nil)
		s.True(errors.Is(err, persistence.ErrNotOwnedByUser))
	})
	s.Run("id is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", nil)
		s.Require().NoError(err)
		err = r.UpdateCartOfUser(ctx, "user", "ID", nil)
		s.True(errors.Is(err, persistence.ErrNotFound))
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"805099eb-27ba-4db8-a7fe-57b63e8afa84",
			"0b2cd45b-879d-42b7-a35a-edbbe7a719f6",
			"6f2b4ecb-56f3-49a3-a5ad-0e946ff26647",
			"44b9ffee-8ad0-40ec-b3c1-ee8d04e0074e",
			"9f879b83-6cfc-4a45-b4be-7cabd1ff17dc",
			"e6cf7804-937c-4e77-9952-f15a5e3bdec5",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.CreateCart(ctx, "user", id, nil)
				s.Require().NoError(err)
			}
			for _, id := range ids {
				err := r.UpdateCartOfUser(ctx, "user", id, nil)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(ids[:3])
		go do(ids[3:])
		wg.Wait()
	})
}

// TestFindAllUnlockedCartsOfUser tests finding all unlocked carts of a user.
func (s *CartRepositoryTestSuite) TestFindAllUnlockedCartsOfUser() {
	s.Run("finds cart with positions", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "8a0f04c7-babb-4ae6-a003-03637cb4396a", "4a33699b-afc5-41e7-b22f-3cdfca5952f8", []struct {
			ProductID string
			Quantity  int
			Price     int
		}{
			{"eb8013e1-74ec-4c20-b57f-19d7a47c8bb0", 1, 123},
			{"80d96241-96de-486e-a9bd-5f31dfb59405", 9, 987},
		})
		s.Require().NoError(err)
		carts, err := r.FindAllUnlockedCartsOfUser(ctx, "8a0f04c7-babb-4ae6-a003-03637cb4396a")
		s.NoError(err)
		s.ElementsMatch([]model.Position{
			{
				ProductID: "eb8013e1-74ec-4c20-b57f-19d7a47c8bb0",
				Quantity:  1,
				Price:     123,
			}, {
				ProductID: "80d96241-96de-486e-a9bd-5f31dfb59405",
				Quantity:  9,
				Price:     987,
			},
		}, carts[0].Positions)
		carts[0].Positions = nil
		s.Equal([]*model.Cart{
			{
				ID:     "4a33699b-afc5-41e7-b22f-3cdfca5952f8",
				Locked: false,
			},
		}, carts)
		s.Run("after updating it", func() {
			err := r.UpdateCartOfUser(ctx, "8a0f04c7-babb-4ae6-a003-03637cb4396a", "4a33699b-afc5-41e7-b22f-3cdfca5952f8", []struct {
				ProductID string
				Quantity  int
				Price     int
			}{
				{"58a89337-e6e3-4ed8-b6b8-1999f79d48d5", 5, -100}, // new one
				{"eb8013e1-74ec-4c20-b57f-19d7a47c8bb0", 1, 123},
				// removed 80d96241-96de-486e-a9bd-5f31dfb59405
			})
			s.Require().NoError(err)
			carts, err := r.FindAllUnlockedCartsOfUser(ctx, "8a0f04c7-babb-4ae6-a003-03637cb4396a")
			s.NoError(err)
			s.ElementsMatch([]model.Position{
				{
					ProductID: "58a89337-e6e3-4ed8-b6b8-1999f79d48d5",
					Quantity:  5,
					Price:     -100,
				}, {
					ProductID: "eb8013e1-74ec-4c20-b57f-19d7a47c8bb0",
					Quantity:  1,
					Price:     123,
				},
			}, carts[0].Positions)
			carts[0].Positions = nil
			s.Equal([]*model.Cart{
				{
					ID:     "4a33699b-afc5-41e7-b22f-3cdfca5952f8",
					Locked: false,
				},
			}, carts)
		})
		s.Run("after removing all positions", func() {
			err := r.UpdateCartOfUser(ctx, "8a0f04c7-babb-4ae6-a003-03637cb4396a", "4a33699b-afc5-41e7-b22f-3cdfca5952f8", nil)
			s.Require().NoError(err)
			carts, err := r.FindAllUnlockedCartsOfUser(ctx, "8a0f04c7-babb-4ae6-a003-03637cb4396a")
			s.NoError(err)
			s.Equal([]*model.Cart{
				{
					ID:        "4a33699b-afc5-41e7-b22f-3cdfca5952f8",
					Positions: []model.Position{},
					Locked:    false,
				},
			}, carts)
		})
	})
	s.Run("finds many", func() {
		r := s.NewRepository()
		for _, id := range []string{
			"ec9d12ab-a7e8-4e27-8f58-4ef62f14d82c",
			"d1111e81-6d8d-4531-bd2f-294fa41eab9b",
			"9b312fc0-4867-42f5-948c-731582193a3d",
			"e7e08f45-0cfd-45a5-8ae4-3bfef52bf590",
			"3f616dc1-ad44-4720-9fa0-02985896ee5d",
			"dc5be657-5bbc-48eb-a529-bf60107bd725",
			"d69829b0-ec64-4608-88f6-3c5005fae6e6",
		} {
			err := r.CreateCart(ctx, "821a9932-f585-4d5a-a383-17091b55adcd", id, nil)
			s.Require().NoError(err)
		}
		carts, err := r.FindAllUnlockedCartsOfUser(ctx, "821a9932-f585-4d5a-a383-17091b55adcd")
		s.NoError(err)
		s.ElementsMatch([]*model.Cart{
			{ID: "ec9d12ab-a7e8-4e27-8f58-4ef62f14d82c", Positions: []model.Position{}},
			{ID: "d1111e81-6d8d-4531-bd2f-294fa41eab9b", Positions: []model.Position{}},
			{ID: "9b312fc0-4867-42f5-948c-731582193a3d", Positions: []model.Position{}},
			{ID: "e7e08f45-0cfd-45a5-8ae4-3bfef52bf590", Positions: []model.Position{}},
			{ID: "3f616dc1-ad44-4720-9fa0-02985896ee5d", Positions: []model.Position{}},
			{ID: "dc5be657-5bbc-48eb-a529-bf60107bd725", Positions: []model.Position{}},
			{ID: "d69829b0-ec64-4608-88f6-3c5005fae6e6", Positions: []model.Position{}},
		}, carts)
	})
	s.Run("no carts exist", func() {
		r := s.NewRepository()
		carts, err := r.FindAllUnlockedCartsOfUser(ctx, "22c5bb43-b950-4718-b2c1-d083eed57fe7")
		s.NoError(err)
		s.Equal([]*model.Cart{}, carts)
	})
	s.Run("returns only carts of the user", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			userID, id string
		}{
			{"userA", "158d7eb1-c82f-46fc-9075-bfa3f545d2fd"},
			{"userA", "bc8b368b-b412-4087-995d-b199e5dffb8c"},
			{"userA", "bf9350a2-8a34-4ae3-8672-590cf740b19d"},
			{"userB", "42f8d3d5-a9a6-4f83-a940-7de23e05a93b"},
			{"userB", "8c33de17-8064-4b32-b4dc-ae3a5abc6234"},
			{"userB", "5b33586b-7a90-45a4-94f8-f396dcd45261"},
			{"userB", "dd742448-09ea-4cc1-adb8-f00bba75c520"},
		} {
			err := r.CreateCart(ctx, c.userID, c.id, nil)
			s.Require().NoError(err)
		}
		carts, err := r.FindAllUnlockedCartsOfUser(ctx, "userA")
		s.NoError(err)
		s.ElementsMatch([]*model.Cart{
			{ID: "158d7eb1-c82f-46fc-9075-bfa3f545d2fd", Positions: []model.Position{}},
			{ID: "bc8b368b-b412-4087-995d-b199e5dffb8c", Positions: []model.Position{}},
			{ID: "bf9350a2-8a34-4ae3-8672-590cf740b19d", Positions: []model.Position{}},
		}, carts)
	})
	s.Run("does not mix up positions", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			id string
			// one product:
			productID       string
			quantity, price int
		}{
			{"id1", "product1", 1, 123},
			{"id2", "product2", 2, 234},
			{"id3", "product3", 3, 345},
		} {
			err := r.CreateCart(ctx, "user", c.id, []struct {
				ProductID string
				Quantity  int
				Price     int
			}{
				{ProductID: c.productID, Quantity: c.quantity, Price: c.price},
			})
			s.Require().NoError(err)
		}
		carts, err := r.FindAllUnlockedCartsOfUser(ctx, "user")
		s.NoError(err)
		s.ElementsMatch([]*model.Cart{
			{ID: "id1", Positions: []model.Position{{ProductID: "product1", Quantity: 1, Price: 123}}},
			{ID: "id2", Positions: []model.Position{{ProductID: "product2", Quantity: 2, Price: 234}}},
			{ID: "id3", Positions: []model.Position{{ProductID: "product3", Quantity: 3, Price: 345}}},
		}, carts)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			userID, id string
		}
		cases := []singleCase{
			{"8924b52e-1868-409e-9283-33149e84e768", "ead2b83a-20b6-45c6-b43b-6389fd96f7fb"},
			{"a8421c6d-4927-4cf0-ba08-82352e9f6429", "210d0aba-62d6-450e-847e-11f94962fff5"},
			{"18bf21cb-f1c0-4de4-a2f4-355e23c3dd06", "1a791edd-0adb-4a65-b291-41798e652a49"},
			{"9151c02b-d055-441f-a902-385e2d46eeaf", "4937d0e1-a624-4ce4-987a-5bf2c8634472"},
			{"5fc60b04-3e6f-4ef6-aaa9-798437303991", "6a5d4bf9-318c-4e3f-8c3b-e23e095000a2"},
			{"0b2b65e0-28fb-44b0-aa95-7aff9a54b9a9", "f3566ee3-31bb-4d82-974d-b9420b7801c9"},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.CreateCart(ctx, c.userID, c.id, nil)
				s.Require().NoError(err)
			}
			for _, c := range cases {
				_, err := r.FindAllUnlockedCartsOfUser(ctx, c.userID)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(cases[:3])
		go do(cases[3:])
		wg.Wait()
	})
}

// TestFindCartOfUser tests finding a cart of a user.
func (s *CartRepositoryTestSuite) TestFindCartOfUser() {
	s.Run("finds a cart", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx,
			"user",
			"id",
			[]struct {
				ProductID string
				Quantity  int
				Price     int
			}{
				{"49dd8502-2d5a-4c71-ac50-e0affcba22c2", 1, 999},
				{"f99a9ea1-8c0f-4e86-b778-18b2537f6234", 999, -1},
			}, // products
		)
		s.Require().NoError(err)
		cart, err := r.FindCartOfUser(ctx, "user", "id")
		s.NoError(err)
		s.Equal(&model.Cart{
			ID: "id",
			Positions: []model.Position{
				{ProductID: "49dd8502-2d5a-4c71-ac50-e0affcba22c2", Quantity: 1, Price: 999},
				{ProductID: "f99a9ea1-8c0f-4e86-b778-18b2537f6234", Quantity: 999, Price: -1},
			},
		}, cart)
	})
	s.Run("user is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", nil)
		s.Require().NoError(err)
		cart, err := r.FindCartOfUser(ctx, "USER", "id")
		s.True(errors.Is(err, persistence.ErrNotOwnedByUser))
		s.Nil(cart)
	})
	s.Run("id is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", nil)
		s.Require().NoError(err)
		cart, err := r.FindCartOfUser(ctx, "user", "ID")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(cart)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"5353f5fb-0553-47b8-8a8d-4cc1f5b9f5e8",
			"b0be6663-50cc-49d4-aaa6-ad45479e76ee",
			"53a65e23-4274-4e24-9a2e-757bf4935a2e",
			"f8b68086-0907-4819-a599-cb3d5d498067",
			"10cbe2bf-4cff-48d5-b912-63e457b12bad",
			"f7865585-5e94-42a3-ac6b-44f3287fa349",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.CreateCart(ctx, "user", id, nil)
				s.Require().NoError(err)
			}
			for _, id := range ids {
				_, err := r.FindCartOfUser(ctx, "user", id)
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
		err := r.CreateCart(ctx,
			"user",
			"id",
			[]struct {
				ProductID string
				Quantity  int
				Price     int
			}{
				{"04d2c9a8-068d-40ac-acd7-7bf3f5357953", 2, 500},
			}, // products
		)
		s.Require().NoError(err)
		cart, err := r.FindCartOfUser(ctx, "user", "id")
		s.Require().NoError(err)
		s.Require().Equal(&model.Cart{
			ID: "id",
			Positions: []model.Position{
				{ProductID: "04d2c9a8-068d-40ac-acd7-7bf3f5357953", Quantity: 2, Price: 500},
			},
		}, cart)
		// changing the result ...
		cart.ID = "changed"
		cart.Locked = true
		cart.Positions[0].ProductID = "changed"
		cart.Positions = nil
		// ... does not have any side effects
		cart, err = r.FindCartOfUser(ctx, "user", "id")
		s.NoError(err)
		s.Equal(&model.Cart{
			ID: "id",
			Positions: []model.Position{
				{ProductID: "04d2c9a8-068d-40ac-acd7-7bf3f5357953", Quantity: 2, Price: 500},
			},
		}, cart)
	})
}

// TestDeleteCartOfUser tests deleting a cart of a user.
func (s *CartRepositoryTestSuite) TestDeleteCartOfUser() {
	s.Run("deletes a cart", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", nil)
		s.Require().NoError(err)
		err = r.DeleteCartOfUser(ctx, "user", "id")
		s.Require().NoError(err)
		s.Run("prevents deleting it again", func() {
			err := r.DeleteCartOfUser(ctx, "user", "id")
			s.True(errors.Is(err, persistence.ErrDeleted))
		})
		s.Run("prevents accessing it", func() {
			_, err := r.FindCartOfUser(ctx, "user", "id")
			s.True(errors.Is(err, persistence.ErrDeleted))
			carts, err := r.FindAllUnlockedCartsOfUser(ctx, "user")
			s.Empty(carts)
		})
		s.Run("prevents re-creating it", func() {
			err := r.CreateCart(ctx, "user", "id", nil)
			s.True(errors.Is(err, persistence.ErrConflict))
		})
		s.Run("prevents updating it", func() {
			err := r.UpdateCartOfUser(ctx, "user", "id", nil)
			s.True(errors.Is(err, persistence.ErrDeleted))
		})
	})
	s.Run("user is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", nil)
		s.Require().NoError(err)
		err = r.DeleteCartOfUser(ctx, "USER", "id")
		s.True(errors.Is(err, persistence.ErrNotOwnedByUser))
	})
	s.Run("id is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateCart(ctx, "user", "id", nil)
		s.Require().NoError(err)
		err = r.DeleteCartOfUser(ctx, "user", "ID")
		s.True(errors.Is(err, persistence.ErrNotFound))
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		ids := []string{
			"1405a152-2dc7-4695-91b3-cfe5359f8019",
			"cedcaf96-7506-467d-8e4c-81c2d3be0a37",
			"ca5d5302-33cb-4a23-87e8-3afb5e3bab2e",
			"fa654176-7929-4ae0-8d48-2ef483f97332",
			"d6e47667-ce97-467f-b0f2-edf8c35710be",
			"59e247f7-80cc-47a1-986d-2d32a40a6d78",
		}
		do := func(ids []string) {
			defer wg.Done()
			for _, id := range ids {
				err := r.CreateCart(ctx, "user", id, nil)
				s.Require().NoError(err)
			}
			for _, id := range ids {
				err := r.DeleteCartOfUser(ctx, "user", id)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(ids[:3])
		go do(ids[3:])
		wg.Wait()
	})
}
