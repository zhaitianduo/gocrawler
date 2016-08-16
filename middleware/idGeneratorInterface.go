package middleware

type IdGenerator interface {
	GetUint32() uint32
}
