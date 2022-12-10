package neo4jRepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"social-graph/model"
)

type RepositoryNeo4j struct {
	driver neo4j.Driver
	tracer trace.Tracer
}

const (
	query       = "MATCH (u:User)%s(following)\nWHERE u.username = $username RETURN following.username as username, following.private as private"
	followQuery = "Match(f:User {username:$from })\nMatch(t:User {username:$to}) \nMerge(f)-[:%s]->(t)"
	removeQuery = "MATCH (f {username: $from})-[r:%s]->(t {username: $to})DELETE r"
)

func NewRepositoryNeo4j(tracer trace.Tracer) (*RepositoryNeo4j, error) {
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
		tracer,
	}, err
}

func (repo *RepositoryNeo4j) GetUser(ctx context.Context, username string) (*model.User, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.GetUser")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	rez, er := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("match (u:User {username: $username}) return u.username as username, u.private as private", map[string]interface{}{"username": username})
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			log.Println(err)
			return nil, err
		}
		result.Next()
		r := result.Record()
		if r == nil {
			return nil, nil
		}
		u, _ := r.Get("username")
		p, _ := r.Get("private")
		return model.User{Username: u.(string), IsPrivate: p.(bool)}, nil
	})
	if er != nil {
		span.SetStatus(codes.Error, er.Error())
		return &model.User{}, er
	}

	return rez.(*model.User), nil
}

func (repo *RepositoryNeo4j) CreateNewUser(ctx context.Context, user model.User) error {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.CreateNewUser")
	defer span.End()

	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer session.Close()
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {

		_, err := tx.Run("Merge(u:User {username: $username, town: $town, gender: $gender, age: $age, private: $private})", map[string]interface{}{"username": user.Username, "town": user.Town, "gender": user.Gender, "age": user.YearOfBirth, "private": user.IsPrivate})
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			log.Println(err)
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil

}

func (repo *RepositoryNeo4j) SaveApprovedFollow(ctx context.Context, fromUsername string, toUsername string) error {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.SaveApprovedFollow")
	defer span.End()
	return repo.SaveFollow(ctx, fromUsername, toUsername, fmt.Sprintf(followQuery, "FOLLOWS"))
}
func (repo *RepositoryNeo4j) SaveFollowRequest(ctx context.Context, fromUsername string, toUsername string) error {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.SaveFollowRequest")
	defer span.End()
	return repo.SaveFollow(ctx, fromUsername, toUsername, fmt.Sprintf(followQuery, "FOLLOWS_REQUEST"))
}
func (repo *RepositoryNeo4j) SaveFollow(ctx context.Context, fromUsername string, toUsername string, query string) error {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.SaveFollow")
	defer span.End()

	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer session.Close()
	_, er := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(query, map[string]interface{}{"from": fromUsername, "to": toUsername})
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			log.Println(err)
			return nil, err
		}
		return nil, nil
	})
	if er != nil {
		span.SetStatus(codes.Error, er.Error())
		return er
	}
	return nil
}
func (repo *RepositoryNeo4j) RemoveApprovedFollow(ctx context.Context, fromUsername string, toUsername string) error {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.RemoveApprovedFollow")
	defer span.End()
	return repo.RemoveFollow(ctx, fromUsername, toUsername, fmt.Sprintf(removeQuery, "FOLLOWS"))
}
func (repo *RepositoryNeo4j) RemoveFollowRequest(ctx context.Context, fromUsername string, toUsername string) error {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.RemoveFollowRequest")
	defer span.End()
	return repo.RemoveFollow(ctx, fromUsername, toUsername, fmt.Sprintf(removeQuery, "FOLLOWS_REQUEST"))
}
func (repo *RepositoryNeo4j) RemoveFollow(ctx context.Context, fromUsername string, toUsername string, query string) error {

	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.RemoveFollow")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer session.Close()
	_, er := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(query, map[string]interface{}{"from": fromUsername, "to": toUsername})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return nil, nil
	})
	if er != nil {
		span.SetStatus(codes.Error, er.Error())
		return er
	}
	return nil
}
func (repo *RepositoryNeo4j) GetFollowing(ctx context.Context, username string) ([]model.User, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.GetFollowing")
	defer span.End()
	return repo.GetAllFollow(ctx, username, fmt.Sprintf(query, "-[:FOLLOWS]->"))
}
func (repo *RepositoryNeo4j) GetFollowers(ctx context.Context, username string) ([]model.User, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.GetFollowers")
	defer span.End()
	return repo.GetAllFollow(ctx, username, fmt.Sprintf(query, "<-[:FOLLOWS]-"))
}
func (repo *RepositoryNeo4j) GetAllFollow(ctx context.Context, username string, query string) ([]model.User, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.GetAllFollow")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	rez, _ := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run(query, map[string]interface{}{"username": username})
		if err != nil {
			log.Println(err)

			return nil, err
		}
		var results []model.User
		for records.Next() {
			record := records.Record()
			u, _ := record.Get("username")
			p, _ := record.Get("private")

			results = append(results, model.User{Username: u.(string), IsPrivate: p.(bool)})
		}
		return results, nil
	})

	if rez == nil || rez.([]model.User) == nil {
		return []model.User{}, nil
	}
	return rez.([]model.User), nil
}

func (repo *RepositoryNeo4j) GetAllFollowRequests(ctx context.Context, username string) ([]model.User, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.GetAllFollowRequests")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	rez, _ := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (u:User)<-[:FOLLOWS_REQUEST]-(request) WHERE u.username = $username RETURN request.username as username, request.private as private", map[string]interface{}{"username": username})
		if err != nil {
			log.Println(err)

			return nil, err
		}
		var results []model.User
		for records.Next() {
			record := records.Record()
			u, _ := record.Get("username")
			p, _ := record.Get("private")

			results = append(results, model.User{Username: u.(string), IsPrivate: p.(bool)})
		}
		return results, nil
	})
	if rez == nil || rez.([]model.User) == nil {
		return []model.User{}, nil
	}
	return rez.([]model.User), nil
}

func (repo *RepositoryNeo4j) CheckIfFollowRequestExists(ctx context.Context, usernameFrom string, usernameTo string) (bool, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.CheckIfFollowExists")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	rez, _ := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (f:User {username: $from }), (t:User {username: $to}) RETURN EXISTS( (f)-[:FOLLOWS_REQUEST]->(t)) as rez", map[string]interface{}{"from": usernameFrom, "to": usernameTo})
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

func (repo *RepositoryNeo4j) CheckIfFollowExists(ctx context.Context, usernameFrom string, usernameTo string) (bool, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.CheckIfFollowExists")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	rez, _ := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (f:User {username: $from }), (t:User {username: $to}) RETURN EXISTS( (f)-[:FOLLOWS]->(t)) as rez", map[string]interface{}{"from": usernameFrom, "to": usernameTo})
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

func (repo *RepositoryNeo4j) CanAccessTweetOfAnotherUser(ctx context.Context, usernameFromToken string, usernameForAccess string) (bool, error) {
	if usernameFromToken == usernameForAccess {
		return true, nil
	}
	userForAccess, _ := repo.GetUser(ctx, usernameForAccess)
	if userForAccess == nil {
		return false, errors.New("user don't exists")
	}
	if !userForAccess.IsPrivate {
		return true, nil
	}

	return repo.CheckIfFollowExists(ctx, usernameFromToken, usernameForAccess)

}
func (repo *RepositoryNeo4j) AcceptRejectFollowRequest(ctx context.Context, from string, to string, approved bool) error {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.AcceptRejectFollowRequest")
	defer span.End()
	exists, _ := repo.CheckIfFollowRequestExists(ctx, from, to)
	if !exists {
		return nil
	}
	err := repo.RemoveFollowRequest(ctx, from, to)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if approved {
		err := repo.SaveApprovedFollow(ctx, from, to)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
func (repo *RepositoryNeo4j) UpdateUser(ctx context.Context, isPrivate bool, authUsername string) error {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.UpdateUser")
	defer span.End()

	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer session.Close()
	_, er := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("match (u:User {username:$username}) set u.private= $isPrivate", map[string]interface{}{"username": authUsername, "isPrivate": isPrivate})
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			log.Println(err)
			return nil, err
		}
		return nil, nil
	})

	if er != nil {
		span.SetStatus(codes.Error, er.Error())
		return er
	}
	return nil
}
func (repo *RepositoryNeo4j) GetRecommendationsProfile(ctx context.Context, username string) ([]model.User, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.GetRecommendationsProfile")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	rez, _ := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (u:User {username:$username })-[:FOLLOWS*2]-> (r:User) where not (u)-[:FOLLOWS]->(r) and not r.username =~ u.username RETURN r.username as username, r.private as private", map[string]interface{}{"username": username})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var results []model.User
		for records.Next() {
			record := records.Record()
			u, _ := record.Get("username")
			p, _ := record.Get("private")
			results = append(results, model.User{Username: u.(string), IsPrivate: p.(bool)})
		}
		return results, nil
	})

	if rez == nil || rez.([]model.User) == nil {
		return []model.User{}, nil
	}
	return rez.([]model.User), nil
}
func (repo *RepositoryNeo4j) GetAllUsersNotFollowedByUser(ctx context.Context, username string) ([]model.User, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.GetAllUsersNotFollowedByUser")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	rez, _ := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (u:User {username:$username}), (p:User) WHERE NOT (u)-[:FOLLOWS]->(p) and NOT (u)-[:FOLLOWS_REQUEST]->(p) AND p.username <> $myUsername RETURN p.username as username, p.private as private", map[string]interface{}{"username": username, "myUsername": username})
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var results []model.User
		for records.Next() {
			record := records.Record()
			u, _ := record.Get("username")
			p, _ := record.Get("private")

			results = append(results, model.User{Username: u.(string), IsPrivate: p.(bool)})
		}
		return results, nil
	})
	if rez == nil || rez.([]model.User) == nil {
		return []model.User{}, nil
	}
	return rez.([]model.User), nil
}
func (repo *RepositoryNeo4j) GetTargetGroupUser(ctx context.Context, username string, targetUserGroup model.TargetUserGroup) ([]model.User, error) {
	_, span := repo.tracer.Start(ctx, "RepositoryNeo4j.GetTargetGroupUser")
	defer span.End()
	session := repo.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	rez, _ := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("Match (u:User {town: $town, gender: $gender}) match (p:User {username:$username}) WHERE NOT (u)-[:FOLLOWS]->(p) AND u.username <> $usernameBissnis and  20 <= u.age <= 30 return u.username as username",
			map[string]interface{}{"username": username, "town": targetUserGroup.Town, "gender": targetUserGroup.Gender, "minAge": targetUserGroup.MinAge, "maxAge": targetUserGroup.MaxAge, "usernameBissnis": username})
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
