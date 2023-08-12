package entities

type ValidatedSeller struct {
	Seller
	isValidated bool
}

func (vp *ValidatedSeller) IsValid() bool {
	return vp.isValidated
}

func NewValidatedSeller(seller *Seller) (*ValidatedSeller, error) {
	if err := seller.validate(); err != nil {
		return nil, err
	}

	return &ValidatedSeller{
		Seller:      *seller,
		isValidated: true,
	}, nil
}
