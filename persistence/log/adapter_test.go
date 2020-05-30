package log_test

import (
	"testing"

	"github.com/Teelevision/excommerce/persistence"
	"github.com/Teelevision/excommerce/persistence/log"
	"github.com/Teelevision/excommerce/persistence/testsuite"
)

func TestAdapterImplementsPlacedOrderRepository(t *testing.T) {
	suite := &testsuite.PlacedOrderRepositoryTestSuite{
		NewRepository: func() persistence.PlacedOrderRepository {
			return log.NewAdapter()
		},
	}
	suite.RunSuite(t)
}
