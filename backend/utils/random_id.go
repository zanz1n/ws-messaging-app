package utils

import nanoid "github.com/matoous/go-nanoid/v2"

func RandomId() string {
	id, err := nanoid.New(18)
	if err != nil {
		panic(err)
	}
	return id
}
