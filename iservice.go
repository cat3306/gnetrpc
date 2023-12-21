package gnetrpc

type IService interface {
	Init(v ...interface{}) IService
	Alias() string
}

type BaseService struct {
}

func (b *BaseService) Alias() string {
	return ""
}
func (b *BaseService) Init(v ...interface{}) IService {
	return b
}
