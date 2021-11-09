package repository

import (
	"github.com/gocql/gocql"
	"github.com/zytel3301/tg-users-service/internal/domain"
	"github.com/zytell3301/cassandra-query-builder"
	uuid_generator "github.com/zytell3301/uuid-generator"
)

type Repository struct {
	metadata    cassandraQB.TableMetadata
	connection  cassandraQB.Connection
	idGenerator *uuid_generator.Generator
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

func NewUsersRepository(hosts []string, keyspace string, generator *uuid_generator.Generator) (Repository, error) {
	connection := cassandraQB.Connection{
		Cluster: gocql.NewCluster(hosts...),
		Session: nil,
	}
	connection.Cluster.Keyspace = keyspace
	connection.Cluster.Consistency = gocql.All
	session, err := connection.Cluster.CreateSession()
	switch err != nil {
	case true:
		return Repository{}, err
	}

	connection.Session = session
	usersMetadata.Connection = connection.Session

	return Repository{connection: connection, metadata: usersMetadata, idGenerator: generator}, nil
}

func (r Repository) NewUser(user domain.User) (err error) {
	batch := r.connection.Session.NewBatch(gocql.LoggedBatch)
	data := map[string]interface{}{
		"id":            user.Id,
		"username":      user.Username,
		"name":          user.Name,
		"lastname":      user.Lastname,
		"phone":         user.Phone,
		"bio":           user.Bio,
		"online_status": user.Online_status,
		"created_at":    user.Created_at,
	}
	err = cassandraQB.AddId(&data, nil, r.idGenerator)
	switch err != nil {
	case true:
		return
	}
	err = r.metadata.NewRecord(data, batch)
	switch err != nil {
	case true:
		return
	}

	err = r.connection.Session.ExecuteBatch(batch)
	return
}

func (r Repository) UpdateUser(user domain.User) (err error) {
	data := map[string]interface{}{
		"id":            user.Id,
		"username":      user.Username,
		"name":          user.Name,
		"lastname":      user.Lastname,
		"phone":         user.Phone,
		"bio":           user.Bio,
		"online_status": user.Online_status,
		"created_at":    user.Created_at,
	}
	batch := r.connection.Session.NewBatch(gocql.LoggedBatch)

	err = r.metadata.UpdateRecord(map[string]interface{}{"id": user.Id}, data, batch)
	switch err != nil {
	case true:
		return
	}

	err = r.connection.Session.ExecuteBatch(batch)
	return
}

func (r Repository) DeleteUser(id string) error {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)

	return r.metadata.DeleteRecord(map[string]interface{}{"id": id}, batch)
}
