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
