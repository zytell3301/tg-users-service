package repository

import (
	"github.com/gocql/gocql"
	"github.com/zytell3301/cassandra-query-builder"
	"github.com/zytell3301/tg-users-service/internal/domain"
	uuid_generator "github.com/zytell3301/uuid-generator"
)

type Repository struct {
	usersMetadata        cassandraQB.TableMetadata
	usersPkPhoneMetadata cassandraQB.TableMetadata
	connection           cassandraQB.Connection
	idGenerator          *uuid_generator.Generator
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

var usersPkPhoneMetadata = cassandraQB.TableMetadata{
	Keyspace: "tg",
	Pk:       map[string]struct{}{"phone": {}},
	Table:    "users_pk_phone",
	Columns: map[string]struct{}{
		"id":       {},
		"phone":    {},
		"lastname": {},
		"bio":      {},
		"username": {},
	},
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

	return Repository{connection: connection, usersMetadata: usersMetadata, usersPkPhoneMetadata: usersPkPhoneMetadata, idGenerator: generator}, nil
}

func (r Repository) NewUser(user domain.User) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
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
	err = r.usersMetadata.NewRecord(data, batch)
	switch err != nil {
	case true:
		return
	}

	err = r.connection.Session.ExecuteBatch(batch)
	return
}

func (r Repository) UpdateUsername(id string, username string) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	data := map[string]interface{}{"username": username}
	err = r.usersMetadata.UpdateRecord(map[string]interface{}{"id": id}, data, batch)
	switch err != nil {
	case true:
		return
	}

	err = r.connection.Session.ExecuteBatch(batch)
	return
}

func (r Repository) DeleteUser(phone string) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	err = r.usersMetadata.DeleteRecord(map[string]interface{}{"id": id}, batch)
	switch err != nil {
	case true:
		return
	}

	return r.connection.Session.ExecuteBatch(batch)
}
