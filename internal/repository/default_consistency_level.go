package repository

import "github.com/gocql/gocql"

/**
 * DEFAULT CONSISTENCY LEVELS MUST ONLY BE USED IN TEST ENVIRONMENTS
 */
var DefaultConsistencyLevel = ConsistencyLevels{
	NewUser:            gocql.One,
	UpdateUsername:     gocql.One,
	DeleteUser:         gocql.One,
	DoesUserExists:     gocql.One,
	DoesUsernameExists: gocql.One,
	GetUserByUsername:  gocql.One,
	GetUserByPhone:     gocql.One,
	RecordSecurityCode: gocql.One,
	GetSecurityCode:    gocql.One,
}