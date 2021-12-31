package repository

import (
	"github.com/gocql/gocql"
	"github.com/zytell3301/cassandra-query-builder"
	"github.com/zytell3301/tg-users-service/internal/domain"
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
	id, err := r.idGenerator.GenerateV4()
	switch err != nil {
	case true:
		return err
	}
	data := map[string]interface{}{
		"id":            id.String(),
		"username":      user.Username,
		"name":          user.Name,
		"lastname":      user.Lastname,
		"phone":         user.Phone,
		"bio":           user.Bio,
		"online_status": user.Online_status,
		"created_at":    user.Created_at,
	}
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

func (r Repository) DeleteUser(id string) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	err = r.metadata.DeleteRecord(map[string]interface{}{"id": id}, batch)
	switch err != nil {
	case true:
		return
	}

	return r.connection.Session.ExecuteBatch(batch)
}
