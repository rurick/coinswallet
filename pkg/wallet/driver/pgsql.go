package driver

// configuration for database connection
type configuration struct {
	DBName string
	DBUser string
	DBHost string
	DBPass string
}

//return configuration for database connection
func getConfiguration() configuration {
	return configuration{}
}

// Init - initialisation of driver
func Init() error {
	//config := getConfiguration()
	return nil
}

// Driver for work with PostgreSQL database
type PgSql struct {
	name string
}

func (pg *PgSql) Get(name string) error {
	return nil
}

func (pg *PgSql) Save() error {
	return nil
}

func (pg *PgSql) Create() error {
	return nil
}
