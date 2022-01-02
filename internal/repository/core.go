package repository

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/zytell3301/cassandra-query-builder"
	"github.com/zytell3301/tg-users-service/internal/domain"
	uuid_generator "github.com/zytell3301/uuid-generator"
	"time"
)

type Repository struct {
	usersMetadata           cassandraQB.TableMetadata
	usersPkPhoneMetadata    cassandraQB.TableMetadata
	usersPkUsernameMetadata cassandraQB.TableMetadata
	connection              cassandraQB.Connection
	idGenerator             *uuid_generator.Generator
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

var usersPkUsernameMetadata = cassandraQB.TableMetadata{
	Keyspace: "tg",
	Pk:       map[string]struct{}{"username": {}},
	Table:    "users_pk_username",
	Columns: map[string]struct{}{
		"id":       {},
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

	return Repository{
		connection:              connection,
		usersMetadata:           usersMetadata,
		usersPkPhoneMetadata:    usersPkPhoneMetadata,
		usersPkUsernameMetadata: usersPkUsernameMetadata,
		idGenerator:             generator,
	}, nil
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

	r.usersPkPhoneMetadata.NewRecord(map[string]interface{}{
		"id":       id.String(),
		"phone":    data["phone"],
		"lastname": data["lastname"],
		"bio":      data["bio"],
		"username": data["username"],
	}, batch)

	switch err != nil {
	case true:
		return
	}
	err = r.connection.Session.ExecuteBatch(batch)
	return
}

func (r Repository) UpdateUsername(phone string, username string) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	user, err := r.getUserByPhone(phone)

	switch err != nil {
	case true:
		return
	}

	err = r.usersPkPhoneMetadata.UpdateRecord(map[string]interface{}{"phone": phone}, map[string]interface{}{"username": username}, batch)
	switch err != nil {
	case true:
		return
	}

	err = r.usersMetadata.UpdateRecord(map[string]interface{}{"id": user.Id}, map[string]interface{}{"username": username}, batch)
	switch err != nil {
	case true:
		return
	}

	err = r.connection.Session.ExecuteBatch(batch)
	return
}

func (r Repository) DeleteUser(phone string) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	user, err := r.getUserByPhone(phone)
	switch err != nil {
	case true:
		return
	}
	err = r.usersPkPhoneMetadata.DeleteRecord(map[string]interface{}{"phone": phone}, batch)
	switch err != nil {
	case true:
		return
	}

	err = r.usersMetadata.DeleteRecord(map[string]interface{}{"id": user.Id}, batch)

	switch err != nil {
	case true:
		return
	}

	return r.connection.Session.ExecuteBatch(batch)
}

/**
gocql package will default return a not found error and it is not needed to
check for the returned data
*/
func (r Repository) DoesUserExists(phone string) (bool, error) {
	_, err := r.getUserByPhone(phone)

	switch errors.Is(err, gocql.ErrNotFound) {
	case true:
		return false, nil
	}
	switch err != nil {
	case true:
		return false, err
	}
	return true, nil
}

func (r Repository) getUserByUsername(username string) (domain.User, error) {
	user, err := r.usersPkUsernameMetadata.GetRecord(map[string]interface{}{"username": username}, []string{"id"})
	switch err != nil {
	case true:
		return domain.User{}, err
	}
	id := user["id"].(gocql.UUID)
	u, err := r.usersMetadata.GetRecord(map[string]interface{}{"id": id.String()}, []string{"*"})
	switch err != nil {
	case true:
		return domain.User{}, err
	}
	return domain.User{
		Id:            u["id"].(string),
		Name:          u["name"].(string),
		Lastname:      u["lastname"].(string),
		Bio:           u["bio"].(string),
		Username:      username,
		Phone:         u["phone"].(string),
		Online_status: u["online_status"].(bool),
		Created_at:    u["created_at"].(time.Time),
	}, nil
}

func (r Repository) getUserByPhone(phone string) (domain.User, error) {
	user, err := r.usersPkPhoneMetadata.GetRecord(map[string]interface{}{"phone": phone}, []string{"*"})
	return domain.User{
		Id:       user["id"].(string),
		Name:     user["name"].(string),
		Lastname: user["lastname"].(string),
		Bio:      user["lastname"].(string),
		Username: user["username"].(string),
		Phone:    phone,
	}, err
}
