package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) Test__UnitExample() {
	input := 2
	out := 2

	assert.Equal(suite.T(), input, out)
}

func Test_UnitExample(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
