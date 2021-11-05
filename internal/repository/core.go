package repository

import (
	"github.com/gocql/gocql"
	"github.com/zytell3301/cassandra-query-builder"
)

type Repository struct {
	metadata   cassandraQB.TableMetadata
	connection cassandraQB.Connection
}

var usersMetadata = cassandraQB.TableMetadata{
	Keyspace: "tg",
	Pk:       map[string]struct{}{"id": {}},
	Table:    "users",
	Columns: map[string]struct{}{
		"id":            {},
		"name":          {},
		"lastname":      {},
		"bio":           {},
		"username":      {},
		"phone":         {},
		"online_status": {},
		"created_at":    {},
	},
	Ck:         nil,
	DependsOn:  nil,
	Connection: nil,
}

func NewUsersRepository(hosts []string) (Repository, error) {
	connection := cassandraQB.Connection{
		Cluster: gocql.NewCluster(hosts...),
		Session: nil,
	}

	session, err := connection.Cluster.CreateSession()
	switch err != nil {
	case true:
		return Repository{}, err
	}

	connection.Session = session
	usersMetadata.Connection = connection.Session

	return Repository{connection: connection, metadata: usersMetadata}, nil
}
