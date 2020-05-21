package testsuite_test

import (
	"testing"

	"github.com/Teelevision/excommerce/persistence"
	"github.com/Teelevision/excommerce/persistence/inmemory"
	"github.com/Teelevision/excommerce/persistence/testsuite"
)

// Tests the suite against the reference in-memory implementation.
func TestReferenceImplementation(t *testing.T) {
	{ // user
		suite := &testsuite.UserRepositoryTestSuite{
			NewRepository: func() persistence.UserRepository {
				return inmemory.NewAdapter(inmemory.FastLessSecureHashingForTesting())
			},
		}
		suite.RunSuite(t)
	}
	{ // product
		suite := &testsuite.ProductRepositoryTestSuite{
			NewRepository: func() persistence.ProductRepository {
				return inmemory.NewAdapter()
			},
		}
		suite.RunSuite(t)
	}
	{ // cart
		suite := &testsuite.CartRepositoryTestSuite{
			NewRepository: func() persistence.CartRepository {
				return inmemory.NewAdapter()
			},
		}
		suite.RunSuite(t)
	}
	{ // coupon
		suite := &testsuite.CouponRepositoryTestSuite{
			NewRepository: func() persistence.CouponRepository {
				return inmemory.NewAdapter()
			},
		}
		suite.RunSuite(t)
	}
	{ // order
		suite := &testsuite.OrderRepositoryTestSuite{
			NewRepository: func() persistence.OrderRepository {
				return inmemory.NewAdapter()
			},
		}
		suite.RunSuite(t)
	}
}
