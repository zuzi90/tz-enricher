package tests

import (
	"context"
	"fmt"
	"time"
)

func (s *IntegrationTestSuite) TestResolvers() {

	name := "Rivka"
	ctx := context.Background()
	s.Run("get age by name", func() {
		age, err := s.ageResolver.GetAge(ctx, name)
		s.Require().NoError(err)
		s.Require().Equal(67, age)
		fmt.Println(age)
	})

	time.Sleep(1 * time.Second)
	s.Run("get gender by name", func() {
		gender, err := s.genderResolver.GetGender(ctx, name)
		s.Require().NoError(err)
		s.Require().Equal("female", gender)
		fmt.Println(gender)
	})

	time.Sleep(1 * time.Second)
	s.Run("get country by name", func() {
		country, err := s.countryResolver.GetCountry(ctx, name)
		s.Require().NoError(err)
		s.Require().Equal("IL", country)
		fmt.Println(country)
	})
}
