package gnetrpc

type IService interface {
	Init(v ...interface{}) IService
}
type BaseService struct {
}

func (b *BaseService) Init(v ...interface{}) IService {
	return b
}
