package food

type Service interface {
	SaveQuantity(int32, string) (Record, error)
}
