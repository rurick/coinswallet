package wallet

//
// Model implement interfaces for access to database layer
//
type Model interface {
	// Find instance of wallet by account name

	Get(name string) error

	// Save all data of object to database
	Save() error

	// Create new object in database
	Create() error
}
