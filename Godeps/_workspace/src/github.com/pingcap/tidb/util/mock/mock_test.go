// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package mock

import (
	"testing"

	. "github.com/pingcap/check"
)

func TestT(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testMockSuite{})

type testMockSuite struct {
}

type contextKeyType int

func (k contextKeyType) String() string {
	return "mock_key"
}

const contextKey contextKeyType = 0

func (s *testMockSuite) TestContext(c *C) {
	ctx := NewContext()

	ctx.SetValue(contextKey, 1)
	v := ctx.Value(contextKey)
	c.Assert(v, Equals, 1)

	ctx.ClearValue(contextKey)
	v = ctx.Value(contextKey)
	c.Assert(v, IsNil)

	_, err := ctx.GetTxn(false)
	c.Assert(err, IsNil)

	err = ctx.FinishTxn(false)
	c.Assert(err, IsNil)
}
