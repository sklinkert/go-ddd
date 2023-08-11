package entities

type ValidatedProduct struct {
	Product
	isValidated bool
}

func (vp *ValidatedProduct) IsValid() bool {
	return vp.isValidated
}

func NewValidatedProduct(product *Product) (*ValidatedProduct, error) {
	if err := product.validate(); err != nil {
		return nil, err
	}

	return &ValidatedProduct{
		Product:     *product,
		isValidated: true,
	}, nil
}
