package neo4jRepo

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"os"
	"social-graph/model"
)

type RepositoryNeo4j struct {
	driver neo4j.Driver
}

const (
	database = "neo4j"
	query    = "MATCH (u:User)%s(following)\nWHERE u.username = $username RETURN following.username as username"
)

func NewRepositoryNeo4j() (*RepositoryNeo4j, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	url := fmt.Sprintf("neo4j://%s:%s", db, dbport)

	driver, err := neo4j.NewDriver(url, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		return nil, err
	}

	return &RepositoryNeo4j{
		driver,
	}, err
}

func (repo *RepositoryNeo4j) SaveFollow(follow *model.Follows) (err error) {
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func() {
		err = session.Close()
	}()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("Merge(f:User {username:$from })\nMerge(t:User {username:$to}) \nMerge(f)-[:FOLLOWS]->(t)", map[string]interface{}{"from": follow.From.Username, "to": follow.To.Username})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return nil, nil
	})

	return err
}
func (repo *RepositoryNeo4j) RemoveFollow(follow *model.Follows) (err error) {
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func() {
		err = session.Close()
	}()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH (f {username: $from})-[r:FOLLOWS]->(t {username: $to})\nDELETE r", map[string]interface{}{"from": follow.From.Username, "to": follow.To.Username})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return nil, nil
	})

	return err
}
func (repo *RepositoryNeo4j) GetFollowing(username string) (users []model.User, err error) {
	return repo.GetUsers(username, fmt.Sprintf(query, "-[:FOLLOWS]->"))
}
func (repo *RepositoryNeo4j) GetFollowers(username string) (users []model.User, err error) {
	return repo.GetUsers(username, fmt.Sprintf(query, "<-[:FOLLOWS]-"))
}
func (repo *RepositoryNeo4j) GetUsers(username string, query string) (users []model.User, err error) {
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		err = session.Close()
	}()
	rez, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run(query, map[string]interface{}{"username": username})
		if err != nil {
			log.Println(err)

			return nil, err
		}
		var results []model.User
		for records.Next() {
			record := records.Record()
			u, _ := record.Get("username")

			results = append(results, model.User{Username: u.(string)})
		}
		return results, nil
	})
	return rez.([]model.User), nil
}
