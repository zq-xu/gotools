package response

const (
	CustomerErrorSectionCode ErrorSectionCode = iota
	SupplierErrorSectionCode
	SalesmanErrorSectionCode
	DelivererErrorSectionCode
	ProductErrorSectionCode
	OrderErrorSectionCode
)

func RegisterServiceErrorInfo(errType ErrorSectionCode, list ...ServiceErrorInfo) {
	for _, sei := range list {
		sei.ErrorSectionCode = TransServiceCode(errType, sei.ErrorSectionCode)
		RegisterErrorInfo(ServiceErrorSectionCode, sei)
	}
}

func TransServiceCode(errType, code ErrorSectionCode) ErrorSectionCode {
	return errType*100 + code
}
