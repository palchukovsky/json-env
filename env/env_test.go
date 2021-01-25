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
		  		"x": {
		    		"y": {
		      		"z": "val1"
		    		}
		  		},
		  		"f": "val2"
				}
	*/

	_, err := jsonenv.NewEnv("eyJ4Ijp7InkiOnsieiI6InZhbDEifX0sImYiOiJ2YWwyIn0=")
	// Because has padding:
	assert.EqualError(err,
		`failed to decode Base64: "illegal base64 data at input byte 47"`)

	env, err := jsonenv.NewEnv("eyJ4Ijp7InkiOnsieiI6InZhbDEifX0sImYiOiJ2YWwyIn0")
	assert.NoError(err)

	val, err := env.Read("x/y/z")
	assert.NoError(err)
	assert.Equal("val1", val)

	val, err = env.Read("f")
	assert.NoError(err)
	assert.Equal("val2", val)

	_, err = env.Read("y")
	assert.EqualError(err, `path "" doesn't have value key "y"`)

	_, err = env.Read("y/z/x")
	assert.EqualError(err, `path node "y" is not existing in ""`)

	_, err = env.Read("x/y/v")
	assert.EqualError(err, `path "x/y" doesn't have value key "v"`)
}
