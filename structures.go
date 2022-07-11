package main

// Appends to the slice. if len reached cap removed the first item
// This allow a slice to be limited by its capacity
func limitedAppend[T any](c *[]T, i T) {
	if len(*c) >= cap(*c) {
		copy((*c)[:len(*c)-1],(*c)[1:])
		*c = (*c)[:len(*c)-1]
	}

	*c = append(*c, i)
}
