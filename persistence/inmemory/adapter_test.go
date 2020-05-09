package inmemory_test

import (
	"testing"

	"github.com/Teelevision/excommerce/persistence"
	"github.com/Teelevision/excommerce/persistence/inmemory"
	"github.com/Teelevision/excommerce/persistence/testsuite"
)

func TestAdapterImplementsUserRepository(t *testing.T) {
	suite := &testsuite.UserRepositoryTestSuite{
		NewRepository: func() persistence.UserRepository {
			return inmemory.NewAdapter(inmemory.FastLessSecureHashingForTesting())
		},
	}
	suite.RunSuite(t)
}

func TestAdapterImplementsProductRepository(t *testing.T) {
	suite := &testsuite.ProductRepositoryTestSuite{
		NewRepository: func() persistence.ProductRepository {
			return inmemory.NewAdapter()
		},
	}
	suite.RunSuite(t)
}

func TestAdapterImplementsCartRepository(t *testing.T) {
	suite := &testsuite.CartRepositoryTestSuite{
		NewRepository: func() persistence.CartRepository {
			return inmemory.NewAdapter()
		},
	}
	suite.RunSuite(t)
}
