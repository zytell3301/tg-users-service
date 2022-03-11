package repository

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/zytell3301/cassandra-query-builder"
	errors2 "github.com/zytell3301/tg-globals/errors"
	"github.com/zytell3301/tg-users-service/internal/domain"
	"github.com/zytell3301/tg-users-service/internal/errorReporter"
	uuid_generator "github.com/zytell3301/uuid-generator"
	"time"
)

type Repository struct {
	usersMetadata           cassandraQB.TableMetadata
	usersPkPhoneMetadata    cassandraQB.TableMetadata
	usersPkUsernameMetadata cassandraQB.TableMetadata
	securityCodesMetaData   cassandraQB.TableMetadata
	connection              cassandraQB.Connection
	idGenerator             *uuid_generator.Generator
	consistencyLevels       ConsistencyLevels
}

type Configs struct {
	Hosts             []string
	Keyspace          string
	Port              int
	ConsistencyLevels ConsistencyLevels
}

type ConsistencyLevels struct {
	NewUser            gocql.Consistency
	UpdateUsername     gocql.Consistency
	DeleteUser         gocql.Consistency
	DoesUserExists     gocql.Consistency
	DoesUsernameExists gocql.Consistency
	GetUserByUsername  gocql.Consistency
	GetUserByPhone     gocql.Consistency
	RecordSecurityCode gocql.Consistency
	GetSecurityCode    gocql.Consistency
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
		"name":     {},
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

var securityCodesMetaData = cassandraQB.TableMetadata{
	Keyspace: "tg",
	Pk:       map[string]struct{}{"phone": {}},
	Table:    "security_codes",
	Columns: map[string]struct{}{
		"phone": {},
		"code":  {},
	},
}

func NewUsersRepository(configs Configs, generator *uuid_generator.Generator) (Repository, error) {
	connection := cassandraQB.Connection{
		Cluster: gocql.NewCluster(configs.Hosts...),
		Session: nil,
	}
	connection.Cluster.Keyspace = configs.Keyspace
	connection.Cluster.Port = configs.Port
	connection.Cluster.Consistency = gocql.All
	session, err := connection.Cluster.CreateSession()
	switch err != nil {
	case true:
		reportError("creating connection to cassandra database", err)
		return Repository{}, err
	}

	connection.Session = session
	usersMetadata.Connection = connection.Session
	usersPkPhoneMetadata.Connection = connection.Session
	usersPkUsernameMetadata.Connection = connection.Session
	securityCodesMetaData.Connection = connection.Session
	return Repository{
		connection:              connection,
		usersMetadata:           usersMetadata,
		usersPkPhoneMetadata:    usersPkPhoneMetadata,
		usersPkUsernameMetadata: usersPkUsernameMetadata,
		securityCodesMetaData:   securityCodesMetaData,
		idGenerator:             generator,
		consistencyLevels:       configs.ConsistencyLevels,
	}, nil
}

func (r Repository) NewUser(user domain.User) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	id, err := r.idGenerator.GenerateV4()
	switch err != nil {
	case true:
		reportError("generating uuid", err)
		return errors2.InternalError{}
	}
	data := map[string]interface{}{
		"id":       id.String(),
		"name":     user.Name,
		"lastname": user.Lastname,
		"phone":    user.Phone,
	}
	err = r.usersMetadata.NewRecord(data, batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}

	err = r.usersPkPhoneMetadata.NewRecord(map[string]interface{}{
		"id":       id.String(),
		"name":     data["name"],
		"phone":    data["phone"],
		"lastname": data["lastname"],
		"bio":      data["bio"],
		"username": data["username"],
	}, batch)

	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}
	batch.SetConsistency(r.consistencyLevels.NewUser)
	err = r.connection.Session.ExecuteBatch(batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}
	return
}

func (r Repository) UpdateUsername(phone string, username string) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	user, err := r.getUserByPhone(phone)

	/**
	 * In this case gocql.ErrNotFound is not needed to be checked because the username
	 * existence is checked before in core
	 */
	switch err != nil {
	case true:
		return errors2.InternalError{}
	}

	err = r.usersPkPhoneMetadata.UpdateRecord(map[string]interface{}{"phone": phone}, map[string]interface{}{"username": username}, batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}

	err = r.usersMetadata.UpdateRecord(map[string]interface{}{"id": user.Id}, map[string]interface{}{"username": username}, batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}

	switch user.Username != "" {
	case true:
		err = r.usersPkUsernameMetadata.DeleteRecord(map[string]interface{}{"username": user.Username}, batch)
		switch err != nil {
		case true:
			reportQueryError(err)
			return errors2.InternalError{}
		}
	}
	batch.SetConsistency(r.consistencyLevels.UpdateUsername)
	err = r.connection.Session.ExecuteBatch(batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}

	batch = r.connection.Session.NewBatch(gocql.UnloggedBatch)
	err = r.usersPkUsernameMetadata.NewRecord(map[string]interface{}{"username": username, "id": user.Id}, batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}

	batch.SetConsistency(r.consistencyLevels.UpdateUsername)
	err = r.connection.Session.ExecuteBatch(batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}
	return
}

func (r Repository) DeleteUser(phone string) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	user, err := r.getUserByPhone(phone)
	switch err != nil {
	case true:
		return errors2.InternalError{}
	}
	err = r.usersPkPhoneMetadata.DeleteRecord(map[string]interface{}{"phone": phone}, batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}

	err = r.usersMetadata.DeleteRecord(map[string]interface{}{"id": user.Id}, batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}
	batch.SetConsistency(r.consistencyLevels.DeleteUser)
	err = r.connection.Session.ExecuteBatch(batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}
	return
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
		return false, errors2.InternalError{}
	}
	return true, nil
}

func (r Repository) DoesUsernameExists(username string) (bool, error) {
	_, err := r.getIdByUsername(username, r.consistencyLevels.DoesUsernameExists)
	switch err != nil {
	case true:
		switch errors.Is(err, gocql.ErrNotFound) {
		case true:
			return false, nil
		default:
			return false, errors2.InternalError{}
		}
	default:
		return true, nil
	}
}

func (r Repository) getIdByUsername(username string, consistencyLevel gocql.Consistency) (string, error) {
	statement, err := r.usersPkUsernameMetadata.GetSelectStatement(map[string]interface{}{"username": username}, []string{"id"})
	switch err != nil {
	case true:
		reportQueryError(err)
		return "", err
	}
	statement.SetConsistency(consistencyLevel)
	user, err := r.usersPkUsernameMetadata.FetchFromSelectStatement(statement)
	switch err != nil {
	case true:
		switch errors.Is(err, gocql.ErrNotFound) {
		case false:
			reportQueryError(err)
			return "", errors2.InternalError{}
		}
		return "", err
	}
	return user["id"].(gocql.UUID).String(), nil
}

func (r Repository) GetUserByUsername(username string) (domain.User, error) {
	return r.getUserByUsername(username)
}

func (r Repository) getUserByUsername(username string) (domain.User, error) {
	id, err := r.getIdByUsername(username, r.consistencyLevels.GetUserByUsername)
	switch err != nil {
	case true:
		return domain.User{}, errors2.InternalError{}
	}
	user, err := r.usersMetadata.GetRecord(map[string]interface{}{"id": id}, []string{"*"})
	switch err != nil {
	case true:
		switch errors.Is(err, gocql.ErrNotFound) {
		case true:
			return domain.User{}, errors2.EntityNotFound{}
		}
		reportQueryError(err)
		return domain.User{}, errors2.InternalError{}
	}
	return domain.User{
		Id:            user["id"].(gocql.UUID).String(),
		Name:          user["name"].(string),
		Lastname:      user["lastname"].(string),
		Bio:           user["bio"].(string),
		Username:      username,
		Phone:         user["phone"].(string),
		Online_status: user["online_status"].(bool),
		Created_at:    user["created_at"].(time.Time),
	}, nil
}

func (r Repository) GetUserByPhone(phone string) (domain.User, error) {
	user, err := r.getUserByPhone(phone)
	switch err != nil {
	case true:
		switch errors.Is(err, gocql.ErrNotFound) {
		case true:
			return domain.User{}, errors2.EntityNotFound{}
		}
		return domain.User{}, errors2.InternalError{}
	}
	return user, nil
}

func (r Repository) getUserByPhone(phone string) (domain.User, error) {
	statement, err := r.usersPkPhoneMetadata.GetSelectStatement(map[string]interface{}{"phone": phone}, []string{"*"})
	switch err != nil {
	case true:
		reportQueryError(err)
		return domain.User{}, errors2.InternalError{}
	}
	statement.SetConsistency(r.consistencyLevels.GetUserByPhone)
	user, err := r.usersPkPhoneMetadata.FetchFromSelectStatement(statement)
	switch err != nil {
	case true:
		switch errors.Is(err, gocql.ErrNotFound) {
		case true:
			return domain.User{}, err
		}
		reportQueryError(err)
		return domain.User{}, errors2.InternalError{}
	}
	return domain.User{
		Id:       user["id"].(gocql.UUID).String(),
		Name:     user["name"].(string),
		Lastname: user["lastname"].(string),
		Bio:      user["lastname"].(string),
		Username: user["username"].(string),
		Phone:    phone,
	}, nil
}

func (r Repository) RecordSecurityCode(securityCode domain.SecurityCode) (err error) {
	batch := r.connection.Session.NewBatch(gocql.UnloggedBatch)
	err = r.securityCodesMetaData.NewRecord(map[string]interface{}{
		"phone":  securityCode.Phone,
		"code":   securityCode.SecurityCode,
		"action": securityCode.Action,
	}, batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return
	}

	batch.SetConsistency(r.consistencyLevels.RecordSecurityCode)
	err = r.connection.Session.ExecuteBatch(batch)
	switch err != nil {
	case true:
		reportQueryError(err)
		return errors2.InternalError{}
	}
	return nil
}

func (r Repository) GetSecurityCode(phone string) (domain.SecurityCode, error) {
	statement, err := r.securityCodesMetaData.GetSelectStatement(map[string]interface{}{"phone": phone}, []string{"phone", "code", "writetime(code) as created_at", "action"})
	switch err != nil {
	case true:
		reportQueryError(err)
		return domain.SecurityCode{}, errors2.InternalError{}
	}
	statement.SetConsistency(r.consistencyLevels.GetSecurityCode)
	securityCode, err := r.securityCodesMetaData.FetchFromSelectStatement(statement)
	switch err != nil {
	case true:
		switch errors.Is(err, gocql.ErrNotFound) {
		case true:
			return domain.SecurityCode{}, errors2.EntityNotFound{}
		}
		reportQueryError(err)
		return domain.SecurityCode{}, errors2.InternalError{}
	}
	return domain.SecurityCode{
		Phone:        securityCode["phone"].(string),
		SecurityCode: securityCode["code"].(string),
		Action:       securityCode["action"].(string),
		CreatedAt:    parseMicroSeconds(securityCode["created_at"].(int64)),
	}, nil
}

/**
 * Reports errors to central error recorder
 */
func reportError(subject string, err error) {
	errorReporter.ReportError("An error occurred while %s. Error message: %s", subject, err.Error())
}

func reportQueryError(err error) {
	reportError("executing a query", err)
}

func parseMicroSeconds(seconds int64) time.Time {
	return time.Unix(0, 100*seconds)
}
