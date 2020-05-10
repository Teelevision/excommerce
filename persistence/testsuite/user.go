package testsuite

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
	"github.com/stretchr/testify/suite"
)

// UserRepositoryTestSuite is the suite that tests that a user repository
// behaves as expected. Use RunSuite to run it.
type UserRepositoryTestSuite struct {
	suite.Suite
	NewRepository func() persistence.UserRepository
}

// RunSuite runs the test suite.
func (s *UserRepositoryTestSuite) RunSuite(t *testing.T) {
	suite.Run(t, s)
}

// TestCreateUser tests the user creation.
func (s *UserRepositoryTestSuite) TestCreateUser() {
	s.Run("one", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx,
			"8d8084c3-72a8-4a5b-978d-ebe8775c5f0b", // id
			"strobbery",                            // name
			"correct horse battery staple",         // plain password
		)
		s.NoError(err)
	})
	s.Run("with empty id, name and password", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "", "", "")
		s.NoError(err)
	})
	s.Run("many", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			id, name, pass string
		}{
			{"3cd904eb-9ded-422b-b997-3f46cb595344", "name 47", "some pass"},
			{"cf291bf5-9be1-499a-a50e-0c06a3d5c3fd", "name 48", "some pass"},
			{"1a0eb77e-7321-4a03-84ac-cc559b5d9335", "name 49", "some pass"},
			{"49e4ca2e-d786-4483-9d5b-624fad105a89", "name 50", "some pass"},
			{"3791f12f-776d-4b76-a280-1db8906e5570", "name 51", "some pass"},
			{"aeb658a7-9c31-4820-b00b-7b453d192ab2", "name 52", "some pass"},
			{"413be127-4f8c-4bb4-8315-d56088dba10a", "name 53", "some pass"},
		} {
			err := r.CreateUser(ctx, c.id, c.name, c.pass)
			s.Require().NoError(err)
		}
	})
	s.Run("conflict on same id", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "id", "name1", "password")
		s.Require().NoError(err)
		err = r.CreateUser(ctx, "id", "name1", "password")
		s.True(errors.Is(err, persistence.ErrConflict))
	})
	s.Run("conflict on same name", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "id1", "name", "password")
		s.Require().NoError(err)
		err = r.CreateUser(ctx, "id2", "name", "password")
		s.True(errors.Is(err, persistence.ErrConflict))
	})
	s.Run("is case-sensitive", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "id", "name", "password")
		s.Require().NoError(err)
		err = r.CreateUser(ctx, "ID", "NAME", "password")
		s.NoError(err)
	})
	s.Run("supports more complex strings", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "औकखग", "‽ⓐ◐\n👽 乐乑", "\u0000\t\"abc")
		s.NoError(err)
	})
	s.Run("does not trim whitespaces", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, " a", "b\n", "password")
		s.Require().NoError(err)
		err = r.CreateUser(ctx, "a", "b", "password")
		s.NoError(err)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			id, name, pass string
		}
		cases := []singleCase{
			{"ef947eff-ff71-42fa-b4a7-cda9d47e2c48", "name 99", "pass 99"},
			{"4ad7b56b-c028-4a3d-b83e-c56654a9ac9e", "name 100", "pass 100"},
			{"14f2fab3-9c9c-4cf4-b6d5-849d44bdee8c", "name 101", "pass 101"},
			{"497c3956-6d2d-4c37-9980-36385a10a858", "name 102", "pass 102"},
			{"3c6976d4-c852-4b85-b8a3-60f59a29f2a1", "name 103", "pass 103"},
			{"2354d927-8da8-4c02-8998-0e393ec07445", "name 104", "pass 104"},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.CreateUser(ctx, c.id, c.name, c.pass)
				s.Require().NoError(err)
			}
		}
		wg.Add(2)
		go do(cases[:3])
		go do(cases[3:])
		wg.Wait()
	})
}

// TestFindUserByNameAndPassword tests finding a user by name and password.
func (s *UserRepositoryTestSuite) TestFindUserByNameAndPassword() {
	s.Run("finds user", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "85987e96-9bf2-426a-92b3-961bb652eea9", "北京市", "广州市")
		s.Require().NoError(err)
		user, err := r.FindUserByNameAndPassword(ctx, "北京市", "广州市")
		s.NoError(err)
		s.Equal(&model.User{
			ID:   "85987e96-9bf2-426a-92b3-961bb652eea9",
			Name: "北京市",
		}, user)
	})
	s.Run("finds user among many", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			id, name, pass string
		}{
			{"244e6147-e74a-487b-a74b-bd78a1381c4f", "name 138", "pass 138"},
			{"de4daaee-bc0d-4ecf-bfee-54a88905ca37", "name 139", "pass 139"},
			{"5015ca2a-edea-4fb8-9ebb-267455b68696", "name 140", "pass 140"},
			{"9e87f3e9-6652-4476-9986-471321cb1b81", "name 141", "pass 141"},
			{"b9e30bd3-11be-44af-9e42-79846824005a", "name 142", "pass 142"},
			{"c017cb7d-cbd9-46bb-bf95-3e7bab845330", "name 143", "pass 143"},
			{"91e6ac1f-63a1-4617-878e-179142951d72", "name 144", "pass 144"},
		} {
			err := r.CreateUser(ctx, c.id, c.name, c.pass)
			s.Require().NoError(err)
		}
		user, err := r.FindUserByNameAndPassword(ctx, "name 141", "pass 141")
		s.NoError(err)
		s.Equal(&model.User{
			ID:   "9e87f3e9-6652-4476-9986-471321cb1b81",
			Name: "name 141",
		}, user)
	})
	s.Run("wrong password", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "5539cbaa-2c2c-4485-a44e-5698f30b51e9", "marius", "ExCommerce")
		s.Require().NoError(err)
		user, err := r.FindUserByNameAndPassword(ctx, "marius", "excommerce")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(user)
	})
	s.Run("user does not exist", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "f0ca324d-7255-4d18-82ef-b3e59fcd675c", "marius", "ExCommerce")
		s.Require().NoError(err)
		user, err := r.FindUserByNameAndPassword(ctx, "marius2", "ExCommerce")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(user)
	})
	s.Run("works with empty values", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "", "", "")
		s.Require().NoError(err)
		user, err := r.FindUserByNameAndPassword(ctx, "", "")
		s.NoError(err)
		s.Equal(&model.User{}, user)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			id, name, pass string
		}
		cases := []singleCase{
			{"41662045-4f0d-4224-a559-32e68b5a37e2", "name 187", "pass 187"},
			{"d986d0d1-556c-4bba-af13-2c419ad29a12", "name 188", "pass 188"},
			{"71e0b126-7446-4c56-ae3a-41d731fb77e3", "name 189", "pass 189"},
			{"98123b14-5ba0-4bdb-b64b-47f3753d599b", "name 190", "pass 190"},
			{"c1d5f926-20ab-4e3c-b781-c1c8dd4083e1", "name 191", "pass 191"},
			{"f21479d0-8f01-4a23-8224-7616084cc38e", "name 192", "pass 192"},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.CreateUser(ctx, c.id, c.name, c.pass)
				s.Require().NoError(err)
			}
			for _, c := range cases {
				_, err := r.FindUserByNameAndPassword(ctx, c.name, c.pass)
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
		err := r.CreateUser(ctx, "9b19601c-701b-47c8-8f55-68326905f6c5", "北京市", "广州市")
		s.Require().NoError(err)
		user, err := r.FindUserByNameAndPassword(ctx, "北京市", "广州市")
		s.Require().NoError(err)
		s.Require().Equal(&model.User{
			ID:   "9b19601c-701b-47c8-8f55-68326905f6c5",
			Name: "北京市",
		}, user)
		// changing the result ...
		user.ID = "changed"
		user.Name = "changed"
		// ... does not have any side effects
		user, err = r.FindUserByNameAndPassword(ctx, "北京市", "广州市")
		s.NoError(err)
		s.Equal(&model.User{
			ID:   "9b19601c-701b-47c8-8f55-68326905f6c5",
			Name: "北京市",
		}, user)
	})
}

// TestFindUserByIDAndPassword tests finding a user by name and password.
func (s *UserRepositoryTestSuite) TestFindUserByIDAndPassword() {
	s.Run("finds user", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "85987e96-9bf2-426a-92b3-961bb652eea9", "北京市", "广州市")
		s.Require().NoError(err)
		user, err := r.FindUserByIDAndPassword(ctx, "85987e96-9bf2-426a-92b3-961bb652eea9", "广州市")
		s.NoError(err)
		s.Equal(&model.User{
			ID:   "85987e96-9bf2-426a-92b3-961bb652eea9",
			Name: "北京市",
		}, user)
	})
	s.Run("finds user among many", func() {
		r := s.NewRepository()
		for _, c := range []struct {
			id, name, pass string
		}{
			{"ed9b4944-e7d6-4500-9c41-d249fe843f86", "name 230", "pass 230"},
			{"5cbfa727-639b-47a2-831e-886597aad8d3", "name 231", "pass 231"},
			{"f413b753-76ab-457a-9900-ca886ec78477", "name 232", "pass 232"},
			{"e1a43f5a-d236-4d74-9cbd-4b1109bd0569", "name 233", "pass 233"},
			{"105b8bb4-f273-4572-b1cf-0eac1a04f38f", "name 234", "pass 234"},
			{"b5ab835b-488d-483b-a6f7-e825c5baff58", "name 235", "pass 235"},
			{"f61584d3-7442-4575-9a2f-730cc430d9ec", "name 236", "pass 236"},
		} {
			err := r.CreateUser(ctx, c.id, c.name, c.pass)
			s.Require().NoError(err)
		}
		user, err := r.FindUserByIDAndPassword(ctx, "e1a43f5a-d236-4d74-9cbd-4b1109bd0569", "pass 233")
		s.NoError(err)
		s.Equal(&model.User{
			ID:   "e1a43f5a-d236-4d74-9cbd-4b1109bd0569",
			Name: "name 233",
		}, user)
	})
	s.Run("wrong password", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "5539cbaa-2c2c-4485-a44e-5698f30b51e9", "marius", "ExCommerce")
		s.Require().NoError(err)
		user, err := r.FindUserByIDAndPassword(ctx, "5539cbaa-2c2c-4485-a44e-5698f30b51e9", "excommerce")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(user)
	})
	s.Run("user does not exist", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "f0ca324d-7255-4d18-82ef-b3e59fcd675c", "marius", "ExCommerce")
		s.Require().NoError(err)
		user, err := r.FindUserByIDAndPassword(ctx, "04b5aab7-6679-40d7-b261-a53e84b9dd9e", "ExCommerce")
		s.True(errors.Is(err, persistence.ErrNotFound))
		s.Nil(user)
	})
	s.Run("works with empty values", func() {
		r := s.NewRepository()
		err := r.CreateUser(ctx, "", "", "")
		s.Require().NoError(err)
		user, err := r.FindUserByIDAndPassword(ctx, "", "")
		s.NoError(err)
		s.Equal(&model.User{}, user)
	})
	s.Run("works concurrently", func() {
		r := s.NewRepository()
		var wg sync.WaitGroup
		type singleCase struct {
			id, name, pass string
		}
		cases := []singleCase{
			{"58105cb7-b9d6-40a4-8983-0c1cadf680ff", "name 279", "pass 279"},
			{"a1e59ba2-7e52-4c74-8842-172e21d3e81d", "name 280", "pass 280"},
			{"adddf189-04a1-4a26-b9a5-793cb1b86342", "name 281", "pass 281"},
			{"0e1d84c0-c7bb-497f-bb9b-a80d3a577e61", "name 282", "pass 282"},
			{"8ba8c96f-014c-4c6c-8696-9bc67b6d38f4", "name 283", "pass 283"},
			{"5fd5294f-c41e-4aff-9adc-d0759309f557", "name 284", "pass 284"},
		}
		do := func(cases []singleCase) {
			defer wg.Done()
			for _, c := range cases {
				err := r.CreateUser(ctx, c.id, c.name, c.pass)
				s.Require().NoError(err)
			}
			for _, c := range cases {
				_, err := r.FindUserByIDAndPassword(ctx, c.id, c.pass)
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
		err := r.CreateUser(ctx, "9cd5698a-cdf1-42c6-8d47-0f26db05740b", "北京市", "广州市")
		s.Require().NoError(err)
		user, err := r.FindUserByIDAndPassword(ctx, "9cd5698a-cdf1-42c6-8d47-0f26db05740b", "广州市")
		s.Require().NoError(err)
		s.Require().Equal(&model.User{
			ID:   "9cd5698a-cdf1-42c6-8d47-0f26db05740b",
			Name: "北京市",
		}, user)
		// changing the result ...
		user.ID = "changed"
		user.Name = "changed"
		// ... does not have any side effects
		user, err = r.FindUserByIDAndPassword(ctx, "9cd5698a-cdf1-42c6-8d47-0f26db05740b", "广州市")
		s.NoError(err)
		s.Equal(&model.User{
			ID:   "9cd5698a-cdf1-42c6-8d47-0f26db05740b",
			Name: "北京市",
		}, user)
	})
}

var ctx = context.Background()
