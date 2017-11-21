package structures

type BasicList interface {
	// return length of the list.
	Len()

	// append value to the list.
	Append(interface{})

	// get value from the list by index.
	Get(int)

	// find index according to value
	// if exists repetition return the first value
	Index(interface{})
}
