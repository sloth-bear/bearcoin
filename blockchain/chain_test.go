package blockchain

type fakeDB struct {
	fakeLocadChain func() []byte
	fakeFindBlock  func(hash string) []byte
}

func (f fakeDB) FindBlock(hash string) []byte {
	return f.fakeFindBlock(hash)
}

func (fakeDB) SaveBlock(hash string, block []byte) {}

func (fakeDB) DeleteAllBlocks() {}

func (fakeDB) SaveChain(blockchain []byte) {}

func (f fakeDB) LoadChain() []byte {
	return f.fakeLocadChain()
}
