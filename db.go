package main

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/howbazaar/loggo"
	"net/http"
	"sync"
	"time"
)

type Db struct {
	logger               *loggo.Logger
	mutex                sync.Mutex
	primaryRedisHost     string
	secondaryRedisHosts  []string
	primaryConnection    *redis.Conn
	secondaryConnections map[string]*redis.Conn
}

type RedisCmd struct {
	Cmd    string
	Params []interface{}
}

func NewDb(l *loggo.Logger) *Db {
	return &Db{logger: l}
}

func (db *Db) StartPollingSentinel(sentinelUrl string) {
	for {
		func() {

			var data struct {
				Nodes []struct {
					Node struct {
						Host string
						Role string
					}
				}
			}

			// Get server list sentinel
			resp, err := http.Get(sentinelUrl)
			if err != nil {
				db.logger.Errorf("unable to connect to sentinel: %v", err)
				return
			}
			defer resp.Body.Close()
			dec := json.NewDecoder(resp.Body)
			dec.Decode(&data)

			var primaryNode string
			var secondaryNodes []string
			for _, node := range data.Nodes {
				switch node.Node.Role {
				case "primary":
					primaryNode = node.Node.Host + ":6379"
				case "secondary":
					secondaryNodes = append(secondaryNodes, node.Node.Host+":6379")
				}
			}
			db.logger.Debugf("primaryNode:%s, secondaryNodes:%s", primaryNode, secondaryNodes)
			db.setServers(primaryNode, secondaryNodes)
		}()
		time.Sleep(5 * time.Second)
	}
}

func (db *Db) setServers(primaryHost string, secondaryHosts []string) {
	// Create connections before locking the mutex to avoid blocking active threads that
	// need to access redis

	var primaryConn *redis.Conn

	// Connect to the primary server if needed
	if primaryHost != db.primaryRedisHost || db.isConnBroken(db.primaryConnection) {
		c, err := redis.Dial("tcp", primaryHost)
		if err != nil {
			db.logger.Errorf("unable to connect to primary redis host %s: %v", primaryHost, err)
			primaryHost = ""
		} else {
			primaryConn = &c
		}
	} else {
		primaryConn = db.primaryConnection
	}

	// Connect to all secondary servers if needed
	var newSecondaryConns map[string]*redis.Conn = make(map[string]*redis.Conn)
	for _, host := range secondaryHosts {
		//var conn redis.Conn
		//var present bool
		conn, present := db.secondaryConnections[host]

		// If connection already exists, keep it.
		// Otherwise make a new one.
		if present && !db.isConnBroken(conn) {
			newSecondaryConns[host] = conn
		} else {
			conn, err := redis.Dial("tcp", host)
			if err != nil {
				db.logger.Errorf("unable to connect to a secondary redis host %s: %v", host, err)
			} else {
				newSecondaryConns[host] = &conn
			}
		}
	}

	// Now lock the mutex and assign new hosts
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.primaryRedisHost = primaryHost
	db.primaryConnection = primaryConn
	db.secondaryConnections = newSecondaryConns
	//db.logger.Debugf("db:%s", db)
}

func (db *Db) isConnBroken(conn *redis.Conn) bool {
	if conn == nil {
		return true
	}

	return (*conn).Err() != nil
}

func (db *Db) Del(key string) error {
	return db.doDel("DEL", key)
}

func (db *Db) Set(key string, value string) error {
	return db.doSet("SET", key, value)
}

func (db *Db) Srem(set string, key string) error {
	return db.doDel("SREM", set, key)
}

func (db *Db) Get(key string) (string, error) {
	v, err := db.doGet("GET", key)
	if err != nil {
		return "", err
	}

	if v == nil {
		return "", nil
	}
	return string(v.([]byte)), nil
}

func (db *Db) Sadd(key string, value string) error {
	return db.doSet("SADD", key, value)
}

func (db *Db) Sismember(set string, key string) (bool, error) {
	v, err := db.doGet("SISMEMBER", set, key)
	if err != nil {
		return false, err
	}

	var exists int64
	exists = 1
	return exists == v, nil
}

func (db *Db) Hset(hash string, key string, value string) error {
	return db.doSet("HSET", hash, key, value)
}

func (db *Db) Hget(hash string, key string) (string, error) {
	v, err := db.doGet("HGET", hash, key)
	if err != nil {
		return "", err
	}

	return string(v.([]byte)), nil
}

func (db *Db) Hdel(hash string, key string) error {
	return db.doSet("HDEL", hash, key)
}

func (db *Db) Smembers(key string) ([]string, error) {
	reply, err := db.doSend("SMEMBERS", key)
	if err != nil {
		return nil, err
	}

	var values []string
	for _, x := range reply {
		var v, _ = x.([]byte)
		values = append(values, string(v))
	}

	return values, nil
}

func (db *Db) Hgetall(key string, value interface{}) error {
	// TODO: Can't we do this with HGET?
	values, err := redis.Values(db.doGet("HGETALL", key))
	if err != nil {
		db.logger.Errorf("error in hgetall: %v", err)
		return err
	}

	if err = redis.ScanStruct(values, value); err != nil {
		db.logger.Errorf("error in scanstruct: %v", err)
		return err
	}

	return nil
}

func (db *Db) HgetallRaw(key string) (map[string]string, error) {
	values, err := redis.Values(db.doGet("HGETALL", key))
	if err != nil {
		db.logger.Errorf("error in hgetall: %v", err)
		return nil, err
	}

	raw := make(map[string]string)

	for i := 0; i < len(values); i += 2 {
		key := string(values[i].([]byte))
		val := string(values[i+1].([]byte))
		raw[key] = val
	}

	return raw, nil
}

func (db *Db) Exists(key string) (bool, error) {
	exists, err := db.doGet("EXISTS", key)
	if err != nil {
		return false, err
	}

	return exists.(int64) == 1, err
}

func (db *Db) Multi(cmds []*RedisCmd) ([]interface{}, error) {
	runPipeline := func(conn *redis.Conn, cmds []*RedisCmd) ([]interface{}, error) {
		var err error
		if conn != nil {
		} else {
			err = errors.New("connection is missing")
			return nil, err
		}

		results, err := db.doMulti(conn, cmds)
		return results, err
	}

	var results []interface{}
	var err error
	if db.primaryConnection != nil {
		results, err = runPipeline(db.primaryConnection, cmds)
	} else {
		err = errors.New("primary connection is missing")
	}

	for _, conn := range db.secondaryConnections {
		if conn != nil {
			_, err = runPipeline(conn, cmds)
			if err != nil {
				return nil, err
			}
		}
	}

	return results, err

}

func (db *Db) Zcount(key string, min int64, max int64) (int64, error) {
	value, err := db.doGet("ZCOUNT", key, min, max)
	if err != nil {
		return 0, err
	}

	return value.(int64), nil
}

func (db *Db) Expire(key string, timeout int64) error {
	return db.doSet("EXPIRE", key, timeout)
}

func (db *Db) ZrangeByScore(key string, min string, max string) ([]interface{}, error) {
	result, err := db.doGet("ZRANGEBYSCORE", key, min, max)
	if err != nil {
		return nil, err
	}

	return result.([]interface{}), err
}

func (db *Db) ZremRangeByScore(key string, min string, max string) error {
	return db.doDel("ZREMRANGEBYSCORE", key, min, max)
}

func (db *Db) doMulti(conn *redis.Conn, cmds []*RedisCmd) ([]interface{}, error) {
	_, err := (*conn).Do("MULTI")
	if err != nil {
		return nil, err
	}

	for _, cmd := range cmds {
		_, err = (*conn).Do(cmd.Cmd, cmd.Params...)
		if err != nil {
			return nil, err
		}
	}

	results, err := (*conn).Do("EXEC")

	return results.([]interface{}), err
}

func (db *Db) doZremRangeByScore(key string, min string, max string) error {
	return nil
}

func (db *Db) doSet(cmd string, args ...interface{}) error {
	var err error
	if db.primaryConnection != nil {
		_, err = (*db.primaryConnection).Do(cmd, args...)
	} else {
		err = errors.New("primary connection is missing")
	}

	for host, conn := range db.secondaryConnections {
		if conn != nil {
			db.logger.Debugf("updating: %s", host)
			go (*conn).Do(cmd, args...)
		}
	}

	if err != nil {
		db.logger.Errorf("erro writing %s to primary host: %v", cmd, err)
	}

	return err
}

func (db *Db) doGet(cmd string, args ...interface{}) (interface{}, error) {
	var err error
	if db.primaryConnection == nil {
		return nil, errors.New("primary connection is missing")
	}

	result, err := (*db.primaryConnection).Do(cmd, args...)
	if err != nil {
		db.logger.Errorf("erro reading %s from primary host: %v", cmd, err)
	}

	return result, err
}

func (db *Db) doDel(cmd string, args ...interface{}) error {
	return db.doSet(cmd, args...)
}

func (db *Db) doSend(cmd string, args ...interface{}) ([]interface{}, error) {
	if db.primaryConnection == nil {
		return nil, errors.New("primary connection is missing")
	}

	err := (*db.primaryConnection).Send(cmd, args...)
	if err != nil {
		db.logger.Errorf("error in %s: %v", cmd, err)
		return nil, err
	}
	(*db.primaryConnection).Flush()

	reply, err := redis.Values((*db.primaryConnection).Receive())
	if err != nil {
		db.logger.Errorf("error receiving reply in %s: %v", cmd, err)
		return nil, err
	}

	return reply, nil
}
