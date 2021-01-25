package jsonenv_test

import (
	"testing"

	jsonenv "github.com/palchukovsky/json-env/env"
	"github.com/stretchr/testify/assert"
)

func Test_Env(test *testing.T) {
	assert := assert.New(test)

	/*
		{
		  "f": "val2",
		  "x": {
		    "y": {
		      "z": "val1"
		    }
		  }
		}
	*/

	_, err := jsonenv.NewEnv("eyJmIjoidmFsMiIsIngiOnsieSI6eyJ6IjoidmFsMSJ9fX0=")
	// Because has padding:
	assert.EqualError(err,
		`failed to decode Base64: "illegal base64 data at input byte 47"`)

	env, err := jsonenv.NewEnv("eyJmIjoidmFsMiIsIngiOnsieSI6eyJ6IjoidmFsMSJ9fX0")
	assert.NoError(err)

	val, err := env.Read("x/y/z")
	assert.NoError(err)
	assert.NotNil(val)
	assert.Equal("val1", *val)

	val, err = env.Read("f")
	assert.NoError(err)
	assert.NotNil(val)
	assert.Equal("val2", *val)

	val, err = env.Read("y")
	assert.NoError(err)
	assert.Nil(val)

	val, err = env.Read("y/z/x")
	assert.NoError(err)
	assert.Nil(val)

	val, err = env.Read("x/y/v")
	assert.NoError(err)
	assert.Nil(val)

	//////////////////////////////////////////////////////////////////////////////

	dump, err := env.Dump()
	assert.NoError(err)
	assert.Equal(`eyJmIjoidmFsMiIsIngiOnsieSI6eyJ6IjoidmFsMSJ9fX0`, dump)

	//////////////////////////////////////////////////////////////////////////////

	assert.NoError(env.Set("x/y/z2", "val1.2"))
	assert.NoError(env.Set("x/y2/z", "val1.3"))
	assert.NoError(env.Set("f", "val2.2"))
	assert.NoError(env.Set("f2", "val2.3"))

	val, err = env.Read("x/y/z2")
	assert.NoError(err)
	assert.NotNil(val)
	assert.Equal("val1.2", *val)

	val, err = env.Read("x/y2/z")
	assert.NoError(err)
	assert.NotNil(val)
	assert.Equal("val1.3", *val)

	val, err = env.Read("f")
	assert.NoError(err)
	assert.NotNil(val)
	assert.Equal("val2.2", *val)

	val, err = env.Read("f2")
	assert.NoError(err)
	assert.NotNil(val)
	assert.Equal("val2.3", *val)

	val, err = env.Read("x/y/v")
	assert.NoError(err)
	assert.Nil(val)

	//////////////////////////////////////////////////////////////////////////////

	/*
		{
		  "f": "val2.2",
		  "f2": "val2.3",
		  "x": {
		    "y": {
		      "z": "val1",
		      "z2": "val1.2"
		    },
		    "y2": {
		      "z": "val1.3"
		    }
		  }
		}
	*/
	dump, err = env.Dump()
	assert.NoError(err)
	assert.Equal(`eyJmIjoidmFsMi4yIiwiZjIiOiJ2YWwyLjMiLCJ4Ijp7InkiOnsieiI6InZhbDEiLCJ6MiI6InZhbDEuMiJ9LCJ5MiI6eyJ6IjoidmFsMS4zIn19fQ`,
		dump)

	dump, err = env.Export()
	assert.NoError(err)
	assert.Equal(`{"f":"val2.2","f2":"val2.3","x":{"y":{"z":"val1","z2":"val1.2"},"y2":{"z":"val1.3"}}}`,
		dump)
}
