package domain

type Contact struct {
	ID          int
	FullName    string
	PhoneNumber string
}

func NewContact(id int, firstName, lastName, middleName, phoneNumber string) *Contact {
	fullName := lastName + " " + firstName + " " + middleName
	return &Contact{
		ID:          id,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}
