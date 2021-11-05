package repository

import (
	"github.com/gocql/gocql"
	"github.com/zytell3301/cassandra-query-builder"
)

type Repository struct {
	metadata   cassandraQB.TableMetadata
	connection cassandraQB.Connection
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
	return Repository{connection: connection}, nil
}
