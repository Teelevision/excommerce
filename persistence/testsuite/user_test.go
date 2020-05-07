package testsuite_test

import (
	"testing"

	"github.com/Teelevision/excommerce/persistence"
	"github.com/Teelevision/excommerce/persistence/inmemory"
	"github.com/Teelevision/excommerce/persistence/testsuite"
)

// Tests the suite against the reference in-memory implementation.
func TestReferenceImplementation(t *testing.T) {
	suite := &testsuite.UserRepositoryTestSuite{
		NewRepository: func() persistence.UserRepository {
			return inmemory.NewAdapter(inmemory.FastLessSecureHashingForTesting())
		},
	}
	suite.RunSuite(t)
}
