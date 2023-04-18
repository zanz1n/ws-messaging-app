package utils

import nanoid "github.com/matoous/go-nanoid/v2"

const RandomIdLen = 18

func RandomId() string {
	id, err := nanoid.New(RandomIdLen)
	if err != nil {
		panic(err)
	}
	return id
}
