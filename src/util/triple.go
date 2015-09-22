package util

type Triple struct {
	X interface{}
	Y interface{}
	Z interface{}
}

func Triple_new(x, y, z interface{}) *Triple {
	return &Triple{x, y, z}
}
