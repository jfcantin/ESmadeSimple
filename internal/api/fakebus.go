package api

type FakeBus struct{}

func NewFakeBus() *FakeBus {
	return &FakeBus{}
}

func (b *FakeBus) Send() {

}
