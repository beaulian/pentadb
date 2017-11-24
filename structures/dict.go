// Contains the implementation of python-dict-alike dict
package structures

type BasicDict interface {
	// set key and value to map
	Set(interface{}, interface{})

	// get value by key from map
	Get(interface{})

	// judge if key in map, return boolType
	In(interface{})

	// return a array of map's keys
	Keys()

	// return a array of map's values
	Values()

	// return length of keys or values
	Len()
}

type Dict struct {
	Map map[interface{}]interface{}
}


func NewDict() *Dict {
	return &Dict{
		Map: make(map[interface{}]interface{}),
	}
}

func (d *Dict) Set(key interface{}, value interface{}) {
	d.Map[key] = value
}

func (d *Dict) Get(key interface{}) interface{} {
	return d.Map[key]
}

func (d *Dict) In(key interface{}) bool {
	if d.Map[key] == nil {
		return false
	}
	return true
}

func (d *Dict) Len() int {
	return len(d.Map)
}

func (d *Dict) Keys() []interface{} {
	result := make([]interface{}, d.Len())
	for k := range d.Map {
		result = append(result, k)
	}
	return result
}

func (d *Dict) Values() []interface{} {
	result := make([]interface{}, d.Len())
	for _, v := range d.Map {
		result = append(result, v)
	}
	return result
}

