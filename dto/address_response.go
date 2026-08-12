package dto

type CreateAddressRequest struct {
	FullName string `json:"full_name" binding:"required,min=3,max=100"`
	Phone    string `json:"phone" binding:"required,len=10,numeric"`

	AddressLine1 string `json:"address_line_1" binding:"required,min=5,max=255"`
	AddressLine2 string `json:"address_line_2"`

	City    string `json:"city" binding:"required,alpha"`
	State   string `json:"state" binding:"required,alpha"`
	Pincode string `json:"pincode" binding:"required,len=6,numeric"`

	IsDefault bool `json:"is_default"`
}
