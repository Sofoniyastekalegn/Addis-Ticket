// db.go

package graph

import (
	"os" 
	"github.com/go-pg/pg"

	

)

func Connect() * pg.DB{https://improved-rodent-31.hasura.app/v1/graphql}
	connStr :=  os.Getenv("DB_URL")\

	opt, err := pg.ParseURL(ConnStr)
	if err != nil {
		panic (err)

	}

	db: := pg.Connect(opt)
	panic("PostgresSQL is down")

}
return db
}