package models

//
//  NOTE:  this file is NOT autogenerated!!!
//

// AddressFormats returns the set of available address formats
func AddressFormats() (out []Externaldatav1AddressFormat) {
	for _, f := range externaldatav1AddressFormatEnum {
		out = append(out, f.(Externaldatav1AddressFormat))
	}

	return
}

// Curves returns the set of cryptographic curves supported by the API.
func Curves() (out []Externaldatav1Curve) {
	for _, c := range externaldatav1CurveEnum {
		out = append(out, c.(Externaldatav1Curve))
	}

	return
}

// TransactionTypes returns the set of transaction types supported by the API.
func TransactionTypes() (out []Immutableactivityv1TransactionType) {
	for _, t := range immutableactivityv1TransactionTypeEnum {
		out = append(out, t.(Immutableactivityv1TransactionType))
	}

	return
}