package hlt

type Commander struct {
	*Map
	Planet map[int][]*Ship
}

func (c *Commander) SetMap(Map *Map) {
	c.Map = Map
}
