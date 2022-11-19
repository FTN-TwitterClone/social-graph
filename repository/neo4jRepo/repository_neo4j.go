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
	query = "MATCH (u:User)%s(following)\nWHERE u.username = $username RETURN following.username as username"
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

func (repo *RepositoryNeo4j) SaveFollow(follow *model.Follows, isPrivate bool) (err error) {
	exists, _ := repo.CheckIfFollowExists(follow.From.Username, follow.To.Username, true)
	if exists {
		return nil
	}
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func() {
		err = session.Close()
	}()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("Merge(f:User {username:$from })\nMerge(t:User {username:$to}) \nMerge(f)-[:FOLLOWS {approved:$approved}]->(t)", map[string]interface{}{"from": follow.From.Username, "to": follow.To.Username, "approved": !isPrivate})
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
	return repo.GetUsers(username, fmt.Sprintf(query, "-[:FOLLOWS{approved :true}]->"))
}
func (repo *RepositoryNeo4j) GetFollowers(username string) (users []model.User, err error) {
	return repo.GetUsers(username, fmt.Sprintf(query, "<-[:FOLLOWS{approved :true}]-"))
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
	if rez == nil || rez.([]model.User) == nil {
		return []model.User{}, nil
	}
	return rez.([]model.User), nil
}
func (repo *RepositoryNeo4j) CheckIfFollowExists(usernameFrom string, usernameTo string, includeAll bool) (ex bool, err error) {
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		err = session.Close()
	}()
	cond := "{approved :true}"
	if includeAll {
		cond = ""
	}
	rez, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (f:User {username: $from }), (t:User {username: $to}) RETURN EXISTS( (f)-[:FOLLOWS"+cond+"]->(t)) as rez", map[string]interface{}{"from": usernameFrom, "to": usernameTo})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		result.Next()
		r := result.Record()
		if r == nil {
			return false, nil
		}
		res, _ := r.Get("rez")
		return res, nil
	})

	return rez.(bool), nil
}

func (repo *RepositoryNeo4j) AcceptRejectFollowRequest(from string, to string, approved bool) (err error) {
	exists, _ := repo.CheckIfFollowExists(from, to, true)
	if !exists {
		return nil
	}
	if !approved {
		err := repo.RemoveFollow(&model.Follows{From: model.User{Username: from}, To: model.User{Username: to}})
		if err != nil {
			return err
		}
		return nil
	}

	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		err = session.Close()
	}()
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH ((:User {username: $from})-[f:FOLLOWS]->(:User {username: $to}))\nSET f.approved = true ", map[string]interface{}{"from": from, "to": to})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return nil, nil
	})

	return nil
}
