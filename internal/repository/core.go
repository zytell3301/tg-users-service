package repository

import (
	"github.com/zytell3301/cassandra-query-builder"
)

type Repository struct {
	metadata cassandraQB.TableMetadata
	connection cassandraQB.Connection
}

func NewUsersRepository() Repository {
	return Repository{}
}