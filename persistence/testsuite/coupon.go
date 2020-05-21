package testsuite

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/stretchr/testify/suite"
)

// CouponRepositoryTestSuite is the suite that tests that a coupon repository
// behaves as expected. Use RunSuite to run it.
type CouponRepositoryTestSuite struct {
	suite.Suite
	NewRepository func() persistence.CouponRepository
}

// RunSuite runs the test suite.
func (s *CouponRepositoryTestSuite) RunSuite(t *testing.T) {
	suite.Run(t, s)
}

// TestStoreCoupon tests the coupon creation/updates.
func (s *CouponRepositoryTestSuite) TestStoreCoupon() {
	s.Run("one", func() {
		r := s.NewRepository()
		err := r.StoreCoupon(ctx,
			"ORANGE30",                             // code
			"30% off oranges",                      // name
			"0061f256-d4b8-4dd3-85e3-aaaa88a050d2", // product id
			30,                                     // discount
			time.Now().Add(time.Hour),              // expires at
		)
		s.NoError(err)
	})
	s.Run("with empty everything", func() {
		r := s.NewRepository()
		err := r.StoreCoupon(ctx, "", "", "", 0, time.Time{})
		s.NoError(err)
	})
	s.Run("many", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			code, name, productID string
			discount              int
		}{
			{"0E856E61", "name 50", "9032ca95-d07d-4e04-8376-41c02713ed8f", 50},
			{"B6849519", "name 51", "33d8e17d-7af8-4f09-bb1e-216011f1ba89", 51},
			{"05DE29F9", "name 52", "f76d53e4-898e-4613-aa1a-caed73a09990", 52},
			{"60301EB5", "name 53", "c50e328e-a2ae-462d-9408-a65a9a04b78b", 53},
			{"261958DC", "name 54", "6552047c-10e5-4d51-9147-21a6f1f8dd30", 54},
			{"0CBAF5D3", "name 55", "39f70198-359d-4e51-beb4-d023ec1e8809", 55},
			{"9D338E6C", "name 56", "2a9f0d57-65f1-4cbf-8381-51f63fbc7ef5", 56},
		} {
			err := r.StoreCoupon(ctx, c.code, c.name, c.productID, c.discount, time.Now().Add(time.Minute))
			s.Require().NoError(err)
		}
	})
	s.Run("no conflict", func() {
		r := s.NewRepository()
		t := time.Now()
		err := r.StoreCoupon(ctx, "code", "name", "product", 42, t)
		s.Require().NoError(err)
		err = r.StoreCoupon(ctx, "code", "name", "product", 42, t)
		s.NoError(err)
	})
	s.Run("supports more complex strings", func() {
		r := s.NewRepository()
		err := r.StoreCoupon(ctx, "औकखग", "\u0000\t\"abc", "‽ⓐ◐\n👽 乐乑", -1337, time.Now())
		s.NoError(err)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			code, name, productID string
			discount              int
		}
		cases := []singleCase{
			{"C991C4C0", "name 83", "114dd241-06f6-41e3-94ee-4660e9c6d147", 83},
			{"E7A21876", "name 84", "bdbf9cf5-ff1a-4644-884b-9a5995682b1e", 84},
			{"88C56FE3", "name 85", "51f9492e-9ad3-4bf8-9b30-a92bef97ddd9", 85},
			{"6407EB6B", "name 86", "a7abd254-02cb-4012-a5b6-8949a6b8a2fe", 86},
			{"E70DB433", "name 87", "df18412e-10d7-4070-bc30-9d6cbe7ddb63", 87},
			{"9EEAB3D9", "name 88", "26c9e026-32de-4c41-9a53-d9f93db05d05", 88},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.StoreCoupon(ctx, c.code, c.name, c.productID, c.discount, time.Now().Add(time.Hour))
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(cases[:3])
		go do(cases[3:])
		wg.Wait()
	})
}

// TestFindValidCoupon tests finding a non-expired coupon.
func (s *CouponRepositoryTestSuite) TestFindValidCoupon() {
	s.Run("finds coupon", func() {
		r := s.NewRepository()
		t := time.Now().Add(time.Minute)
		err := r.StoreCoupon(ctx, "11D10425", "北京市", "d933396f-20a3-46c6-a4f6-bf6b2ad4fec9", 42, t)
		s.Require().NoError(err)
		coupon, err := r.FindValidCoupon(ctx, "11D10425")
		s.NoError(err)
		s.WithinDuration(t, coupon.ExpiresAt, time.Second)
		coupon.ExpiresAt = time.Time{}
		s.Equal(&model.Coupon{
			Code:      "11D10425",
			Name:      "北京市",
			ProductID: "d933396f-20a3-46c6-a4f6-bf6b2ad4fec9",
			Discount:  42,
		}, coupon)
	})
	s.Run("does not find expired coupon", func() {
		r := s.NewRepository()
		err := r.StoreCoupon(ctx, "9e539b7f", "name", "6a1dc32a-c626-4739-8d85-30c5144b3315", 1, time.Now())
		s.Require().NoError(err)
		coupon, err := r.FindValidCoupon(ctx, "9e539b7f")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(coupon)
	})
	s.Run("finds one among many coupons", func() {
		r := s.NewRepository()
		t := time.Now().Add(time.Minute)
		for _, c := range []struct {
			code, name, productID string
			discount              int
		}{
			{"66350184", "name 129", "9f747ba7-d095-4059-83e8-f058bb12aca4", 129},
			{"888c5d93", "name 130", "eff98b07-4608-4fe5-a7b7-0f34cd536baa", 130},
			{"f5fad7ad", "name 131", "eb2e21e8-877f-44c8-b719-26da1a371850", 131},
			{"2cf878cf", "name 132", "44df107f-1b54-4e92-9382-537f078ca197", 132},
			{"8d4ce1b5", "name 133", "f06ce025-f153-449a-ae71-df8e264b33de", 133},
			{"d4a89757", "name 134", "81c22a54-9434-426e-8158-f8c1bba3bfeb", 134},
			{"3ba909d6", "name 135", "d83025eb-7ab2-467d-826e-69c175c47414", 135},
		} {
			err := r.StoreCoupon(ctx, c.code, c.name, c.productID, c.discount, t)
			s.Require().NoError(err)
		}
		coupon, err := r.FindValidCoupon(ctx, "2cf878cf")
		s.NoError(err)
		s.WithinDuration(t, coupon.ExpiresAt, time.Second)
		coupon.ExpiresAt = time.Time{}
		s.Equal(&model.Coupon{
			Code:      "2cf878cf",
			Name:      "name 132",
			ProductID: "44df107f-1b54-4e92-9382-537f078ca197",
			Discount:  132,
		}, coupon)
	})
	s.Run("does not find other coupon", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			code, name, productID string
			discount              int
			expiresAt             time.Time
		}{
			{"3E8B32FB", "name 158", "c85a9077-f834-4719-b206-aa06550728fe", 158, time.Now().Add(158 * time.Second)},
			{"631768F8", "name 159", "d58ac65d-ba44-4d3a-848a-1a8012b84b35", 159, time.Now().Add(159 * time.Second)},
			{"6E28465C", "name 160", "deab8615-d163-4093-8980-1e89ccd7dde5", 160, time.Now().Add(160 * time.Second)},
		} {
			err := r.StoreCoupon(ctx, c.code, c.name, c.productID, c.discount, c.expiresAt)
			s.Require().NoError(err)
		}
		coupon, err := r.FindValidCoupon(ctx, "49252F6B")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(coupon)
	})
	s.Run("does find the latest version of a coupon", func() {
		r := s.NewRepository()
		t := time.Now().Add(time.Hour)
		for _, c := range []struct {
			name, productID string
			discount        int
		}{
			{"name 176", "a130cb92-e08b-47e1-b32b-7c654946f02a", 176},
			{"name 177", "bbdad944-ee78-4833-bacd-dbdac9574b6f", 177},
			{"name 178", "e80bb51e-175b-4560-81a8-7f4a404baf09", 178},
		} {
			err := r.StoreCoupon(ctx, "ED2BBF84", c.name, c.productID, c.discount, t)
			s.Require().NoError(err)
		}
		coupon, err := r.FindValidCoupon(ctx, "ED2BBF84")
		s.NoError(err)
		s.WithinDuration(t, coupon.ExpiresAt, time.Second)
		coupon.ExpiresAt = time.Time{}
		s.Equal(&model.Coupon{
			Code:      "ED2BBF84",
			Name:      "name 178",
			ProductID: "e80bb51e-175b-4560-81a8-7f4a404baf09",
			Discount:  178,
		}, coupon)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			id, name, productID string
			discount            int
			expiresAt           time.Time
		}
		cases := []singleCase{
			{"CC0B6A8D", "name 178", "f5519ed9-58b2-4090-8eae-58df71b308b6", 178, time.Now().Add(178 * time.Second)},
			{"BBCCED70", "name 179", "c3de413b-8341-4f63-ba75-6715fb08eef0", 179, time.Now().Add(179 * time.Second)},
			{"AC2E824E", "name 180", "8b68c0f3-ac3a-46fc-9a62-beb0b10a7f56", 180, time.Now().Add(180 * time.Second)},
			{"F17EEE8C", "name 181", "0fb3e95d-730e-48a0-8514-8da447af5e2f", 181, time.Now().Add(181 * time.Second)},
			{"189732CF", "name 182", "51cd8d65-053d-404d-adf6-4a2f31fb1a08", 182, time.Now().Add(182 * time.Second)},
			{"F8222AE8", "name 183", "69f29332-9135-47d8-a4ff-bb9126fb8078", 183, time.Now().Add(183 * time.Second)},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.StoreCoupon(ctx, c.id, c.name, c.productID, c.discount, c.expiresAt)
				s.Require().NoError(err)
			}
			for _, c := range cases {
				_, err := r.FindValidCoupon(ctx, c.id)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(cases[:3])
		go do(cases[3:])
		wg.Wait()
	})
	s.Run("changing the result does not have any side effects", func() {
		r := s.NewRepository()
		t := time.Now().Add(time.Hour)
		err := r.StoreCoupon(ctx, "51D11123", "北京市", "7839c105-add5-443b-9018-f6bc73136195", 1337, t)
		s.Require().NoError(err)
		coupon, err := r.FindValidCoupon(ctx, "51D11123")
		s.Require().NoError(err)
		s.Require().Equal(&model.Coupon{
			Code:      "51D11123",
			Name:      "北京市",
			ProductID: "7839c105-add5-443b-9018-f6bc73136195",
			Discount:  1337,
			ExpiresAt: t,
		}, coupon)
		// changing the result ...
		coupon.Code = "changed"
		coupon.Name = "changed"
		coupon.ProductID = "changed"
		coupon.Discount--
		coupon.ExpiresAt = time.Now()
		// ... does not have any side effects
		coupon, err = r.FindValidCoupon(ctx, "51D11123")
		s.NoError(err)
		s.Equal(&model.Coupon{
			Code:      "51D11123",
			Name:      "北京市",
			ProductID: "7839c105-add5-443b-9018-f6bc73136195",
			Discount:  1337,
			ExpiresAt: t,
		}, coupon)
	})
}
