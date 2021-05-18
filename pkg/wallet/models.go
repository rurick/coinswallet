package wallet

//
// Model implement interfaces for access to database layer
//
type accountModel interface {
	// ID return id of wallet account
	ID() string

	// Name return name of wallet account
	Name() string

	// Balance return balance of wallet
	Balance() float64

	// Find instance of wallet by account name
	Find(name string) error

	// Save all data of object to database
	Save() error

	// Create new object in database
	Create() error
}

func accountFactory() {

}
