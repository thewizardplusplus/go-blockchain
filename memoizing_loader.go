package blockchain

type loadingParameters struct {
	cursor interface{}
	count  int
}

type loadingResult struct {
	blocks     BlockGroup
	nextCursor interface{}
}
