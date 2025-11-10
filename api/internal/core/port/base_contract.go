package port

import "microservice/internal/adapter/orm"

type IRepository interface {
	// Tx receives the transactional instance(db.Begin) of DB
	//to handle multiple repositories by the Service layer
	Tx(db orm.ISql)
}
