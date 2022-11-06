package neo4jRepo

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"social-graph/model"
)

type RepositoryNeo4j struct {
	Driver   neo4j.Driver
	Database string
}

func (repo *RepositoryNeo4j) SaveFollow(follow *model.Follows) (err error) {
	session := repo.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: repo.Database})

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
	session := repo.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: repo.Database})

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
func (repo *RepositoryNeo4j) Get(username string, query string) (users []model.User, err error) {
	session := repo.Driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: repo.Database})
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
