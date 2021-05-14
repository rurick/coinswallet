package wallet

//
// Model implement interfaces for access to database layer
//
type Model interface {
	// Name return name of wallet account
	Name() string

	// Find instance of wallet by account name
	Find(name string) error

	// Save all data of object to database
	Save() error

	// Create new object in database
	Create() error
}
