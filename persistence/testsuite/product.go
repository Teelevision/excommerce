package testsuite

import (
	"errors"
	"sync"
	"testing"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/stretchr/testify/suite"
)

// ProductRepositoryTestSuite is the suite that tests that a product repository
// behaves as expected. Use RunSuite to run it.
type ProductRepositoryTestSuite struct {
	suite.Suite
	NewRepository func() persistence.ProductRepository
}

// RunSuite runs the test suite.
func (s *ProductRepositoryTestSuite) RunSuite(t *testing.T) {
	suite.Run(t, s)
}

// TestCreateProduct tests the product creation.
func (s *ProductRepositoryTestSuite) TestCreateProduct() {
	s.Run("one", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx,
			"191f6123-6cbc-4424-a40a-c68018aec9a1", // id
			"Orange",                               // name
			1337,                                   // price
		)
		s.NoError(err)
	})
	s.Run("with empty id, name and price", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx, "", "", 0)
		s.NoError(err)
	})
	s.Run("many", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			id, name string
			price    int
		}{
			{"0e4bd7b9-6895-4abf-8eec-370bab3df055", "name 47", 47},
			{"8f8c0956-ae93-46a7-a093-742c7a9b8efa", "name 48", 48},
			{"2f55d68d-5df6-4c36-8f38-980339a68abb", "name 49", 49},
			{"e5e31584-41aa-481a-aa37-b23aee0df963", "name 50", 50},
			{"829685be-91f0-4788-a457-2497ecac5b37", "name 51", 51},
			{"93137b05-c002-4831-8691-1be66796912d", "name 52", 52},
			{"9cbe377b-1a5a-46a0-8f5d-e2ae3c16b3b3", "name 53", 53},
		} {
			err := r.CreateProduct(ctx, c.id, c.name, c.price)
			s.Require().NoError(err)
		}
	})
	s.Run("conflict on same id", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx, "id", "name1", 1337)
		s.Require().NoError(err)
		err = r.CreateProduct(ctx, "id", "name1", 1337)
		s.True(errors.Is(err, persistence.ErrConflict))
	})
	s.Run("no conflict on same name", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx, "id1", "name", 1337)
		s.Require().NoError(err)
		err = r.CreateProduct(ctx, "id2", "name", 1337)
		s.NoError(err)
	})
	s.Run("is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx, "id", "name", 1337)
		s.Require().NoError(err)
		err = r.CreateProduct(ctx, "ID", "NAME", 1337)
		s.NoError(err)
	})
	s.Run("supports more complex strings", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx, "औकखग \u0000\t\"abc", "‽ⓐ◐\n👽 乐乑", -1337)
		s.NoError(err)
	})
	s.Run("does not trim whitespaces", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx, " a", "b\n", 1337)
		s.Require().NoError(err)
		err = r.CreateProduct(ctx, "a", "b", 1337)
		s.NoError(err)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			id, name string
			price    int
		}
		cases := []singleCase{
			{"620f38d1-142b-4c45-8b0c-6e16bdff02bf", "name 101", 101},
			{"4d0c2de0-a5ad-46eb-92df-4927b0a4e8ff", "name 102", 102},
			{"0202df3c-d178-4ca3-a10c-effd43359d05", "name 103", 103},
			{"2712c536-4882-445a-9419-f962fe86f5e6", "name 104", 104},
			{"ce4810db-cf1e-47b8-a9a7-f27c355461d5", "name 105", 105},
			{"522659ad-636c-46c3-a779-c2fed0391bb5", "name 106", 106},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.CreateProduct(ctx, c.id, c.name, c.price)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(cases[:3])
		go do(cases[3:])
		wg.Wait()
	})
}

// TestFindAllProducts tests finding all products.
func (s *ProductRepositoryTestSuite) TestFindAllProducts() {
	s.Run("finds product", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx, "cc10b2f8-2405-441a-9edf-6ee7d10481c6", "北京市", 1337)
		s.Require().NoError(err)
		products, err := r.FindAllProducts(ctx)
		s.NoError(err)
		s.Equal([]*model.Product{
			{
				ID:    "cc10b2f8-2405-441a-9edf-6ee7d10481c6",
				Name:  "北京市",
				Price: 1337,
			},
		}, products)
	})
	s.Run("finds many products", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			id, name string
			price    int
		}{
			{"075a4c76-a06b-412e-a96f-43cf9ca5cd92", "name 143", 143},
			{"f7062ad5-5d2b-4842-987a-1e58e2a27777", "name 144", 144},
			{"e02c1dbd-2d07-41b7-b368-2258e9a396b8", "name 145", 145},
			{"18970305-6cea-4000-a705-3b70bdaf15da", "name 146", 146},
			{"3c71dea5-74ef-42e5-a22a-2a3b372b685b", "name 147", 147},
			{"c983501f-53f3-4087-b161-b3c257a35fef", "name 148", 148},
			{"59f25fca-a105-454f-bd0a-5711c0689b16", "name 149", 149},
		} {
			err := r.CreateProduct(ctx, c.id, c.name, c.price)
			s.Require().NoError(err)
		}
		products, err := r.FindAllProducts(ctx)
		s.NoError(err)
		s.ElementsMatch([]*model.Product{
			{ID: "075a4c76-a06b-412e-a96f-43cf9ca5cd92", Name: "name 143", Price: 143},
			{ID: "f7062ad5-5d2b-4842-987a-1e58e2a27777", Name: "name 144", Price: 144},
			{ID: "e02c1dbd-2d07-41b7-b368-2258e9a396b8", Name: "name 145", Price: 145},
			{ID: "18970305-6cea-4000-a705-3b70bdaf15da", Name: "name 146", Price: 146},
			{ID: "3c71dea5-74ef-42e5-a22a-2a3b372b685b", Name: "name 147", Price: 147},
			{ID: "c983501f-53f3-4087-b161-b3c257a35fef", Name: "name 148", Price: 148},
			{ID: "59f25fca-a105-454f-bd0a-5711c0689b16", Name: "name 149", Price: 149},
		}, products)
	})
	s.Run("no products exist", func() {
		r := s.NewRepository()
		products, err := r.FindAllProducts(ctx)
		s.NoError(err)
		s.Equal([]*model.Product{}, products)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			id, name string
			price    int
		}
		cases := []singleCase{
			{"0a3ec066-b0d8-4695-9dd1-4d8ba859ed64", "name 180", 180},
			{"2f08d5d9-68fd-48fd-942f-6fd99bd4c56a", "name 181", 181},
			{"36a8aec6-23e3-438e-8d56-8790ab5790d3", "name 182", 182},
			{"519c7aea-667a-44a5-818f-cc11806969b3", "name 183", 183},
			{"0a83e8c5-7479-44ca-931a-f83236476ff3", "name 184", 184},
			{"911f5a56-01c2-4a59-8807-433af11e6755", "name 185", 185},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.CreateProduct(ctx, c.id, c.name, c.price)
				s.Require().NoError(err)
			}
			_, err := r.FindAllProducts(ctx)
			s.Require().NoError(err)
		}
		wg.Add(2)
		go do(cases[:3])
		go do(cases[3:])
		wg.Wait()
	})
}

// TestFindProduct tests finding a product.
func (s *ProductRepositoryTestSuite) TestFindProduct() {
	s.Run("finds product", func() {
		r := s.NewRepository()
		err := r.CreateProduct(ctx, "e6c73a05-c169-452e-8a6c-4afd9ffeb8ba", "北京市", 1337)
		s.Require().NoError(err)
		product, err := r.FindProduct(ctx, "e6c73a05-c169-452e-8a6c-4afd9ffeb8ba")
		s.NoError(err)
		s.Equal(&model.Product{
			ID:    "e6c73a05-c169-452e-8a6c-4afd9ffeb8ba",
			Name:  "北京市",
			Price: 1337,
		}, product)
	})
	s.Run("finds one among many products", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			id, name string
			price    int
		}{
			{"37ed5f9b-684d-4077-bae5-5ac283272a09", "name 223", 223},
			{"8303f166-3e2b-4c4a-a27e-922f6e10d6ad", "name 224", 224},
			{"ebf852c3-0a82-48b9-b51a-1d35bea3a263", "name 225", 225},
			{"11b75e75-63b5-40ba-81c3-bb5fb231854c", "name 226", 226},
			{"494b3baf-3b90-46b9-9026-fc92377e127f", "name 227", 227},
			{"e78d6142-0be8-4355-8cb0-ae4b0d8bcabf", "name 228", 228},
			{"dfeacc59-874d-4f22-8447-9661901f2070", "name 229", 229},
		} {
			err := r.CreateProduct(ctx, c.id, c.name, c.price)
			s.Require().NoError(err)
		}
		product, err := r.FindProduct(ctx, "494b3baf-3b90-46b9-9026-fc92377e127f")
		s.NoError(err)
		s.Equal(&model.Product{
			ID:    "494b3baf-3b90-46b9-9026-fc92377e127f",
			Name:  "name 227",
			Price: 227,
		}, product)
	})
	s.Run("does not find other product", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			id, name string
			price    int
		}{
			{"629faf00-6ffd-4c87-8b7b-7805162709bc", "name 248", 248},
			{"47f8ce74-6581-4be9-b696-1d5ca1065e48", "name 249", 249},
			{"3cc44ebe-177a-4651-8de9-bb2aae53699c", "name 250", 250},
		} {
			err := r.CreateProduct(ctx, c.id, c.name, c.price)
			s.Require().NoError(err)
		}
		product, err := r.FindProduct(ctx, "efc07346-98ac-4c6d-81e9-af3a93c03c71")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(product)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			id, name string
			price    int
		}
		cases := []singleCase{
			{"012887a6-0f3a-447e-a524-2ea2f70264c4", "name 267", 267},
			{"75199849-5dde-4d87-811a-504e1893e31c", "name 268", 268},
			{"74c84ced-0430-4076-b65c-be72930ed7e9", "name 269", 269},
			{"7aee97ac-767e-49c5-9df7-b572dd12875c", "name 270", 270},
			{"ec20f20e-3173-4b8f-a6a2-2194936f81d9", "name 271", 271},
			{"b067259d-c1ef-4915-af2a-bec3ada752da", "name 272", 272},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.CreateProduct(ctx, c.id, c.name, c.price)
				s.Require().NoError(err)
			}
			for _, c := range cases {
				_, err := r.FindProduct(ctx, c.id)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(cases[:3])
		go do(cases[3:])
		wg.Wait()
	})
}
