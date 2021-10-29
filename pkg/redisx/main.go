package redisx

import (
	"errors"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type Config struct {
	prefix string
	client *redis.ClusterClient
}

var config *Config

func Init(address []string, prefix string) error {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    address,
		PoolSize: 200,
	})

	if err := rdb.ReloadState(); err != nil {
		return err
	}

	config = &Config{prefix: prefix, client: rdb}

	return nil
}

func Redis() *Config {
	if config.client == nil {
		log.Println("get redis client nil")
	}

	return config
}

/**
 * @Description: redis Incr
 * @date 8/29/2019 10:52
 */
func (slf *Config) Incr(_key string) (int64, error) {
	if slf.client == nil {
		return 0, errors.New("client is nil")
	}

	code, err := slf.client.Incr(_key).Result()
	if err != nil && err != redis.Nil {
		log.Printf("redis incr err: %v\n", err)
		return 0, err
	}

	return code, nil
}

/**
 * @Description: 为key设置过期时间
 * @date 8/29/2019 10:52
 */
func (slf *Config) SetExpire(_key string, _timer string) error {
	if slf.client == nil {
		return errors.New("client is nil")
	}

	timer, _ := time.ParseDuration(_timer)
	_, err := slf.client.Expire(_key, timer).Result()
	if err != nil && err != redis.Nil {
		log.Printf("redis SetExpire is err: %v\n", err)
		return err
	}

	return nil
}

func (slf *Config) SetString(key string, value string) error {
	err := slf.client.Set(key, value, 0).Err()
	if err != nil {
		log.Println("redis set failed:", err)
		return errors.New("存值失败")
	}
	return nil
}

func (slf *Config) GetString(key string) (string, error) {
	value, err := slf.client.Get(key).Result()
	if err != nil {
		log.Println("redis get failed:", err) // redigo: nil returned
		return "", errors.New("取值失败")
	}
	return value, nil
}

// 锁定某个变量
func (slf *Config) TryLock(_key string, value uint64) (bool, error) {
	setSuccess, err := slf.client.SetNX(_key, value, 0).Result()
	if err != nil {
		return false, err
	}
	if !setSuccess {
		return false, err
	}

	return true, nil
}

// 锁定某个变量
func (slf *Config) UnLock(_key string) (bool, error) {
	delCount, err := slf.client.Del(_key).Result()
	if err != nil {
		log.Print("ERRRRRR  Del " + _key)
		return false, err
	}
	return delCount > 0, nil
}

func (slf *Config) ZRevRangeByScoreWithScores(_key string, _start int64, _stop int64) ([]redis.Z, error) {
	conn := slf.client
	values, err := conn.ZRevRangeWithScores(_key, _start, _stop).Result()
	if err != nil {
		log.Print("ZRevRangeByScoreWithScores err", _key, err)
		return nil, err
	}
	return values, nil
}

func (slf *Config) SAdd(_key string, _val ...string) (int64, error) {
	conn := slf.client
	return conn.SAdd(_key, _val).Result()
}

// LockRedis redis锁定
func (slf *Config) LockRedis(_key string, _ex int64) int64 {
	conn := config.client
	res, err := conn.SetNX(_key, _ex, 30*1e9).Result() // key 不存在返回1，存在返回0
	if err != nil && err != redis.Nil {
		log.Printf("redis lock err:%v\n", err)
	}
	if res {
		return 1
	}
	return 0
}

// UnLockRedis redis解锁
func (slf *Config) UnLockRedis(_key string) int64 {
	conn := config.client
	value, err := conn.Del(_key).Result()
	if err != nil && err != redis.Nil {
		log.Printf("redis del err:%v\n", err)
	}
	return value
}

// SelectAssets 查询资产
func (slf *Config) SelectAssets(_key, _name string) ([]byte, error) {
	conn := config.client
	v, err := conn.HGet(_key, _name).Result()
	return []byte(v), err
}

// UpdateAssets 修改资产
func (slf *Config) UpdateAssets(_key, _name string, _buf []byte) bool {
	conn := config.client
	err := conn.HSet(_key, _name, _buf).Err()
	if err != nil && err != redis.Nil {
		log.Printf("redis del err:%v\n", err)
		return false
	}
	return true
}

// IncrAssets 增加或减少资产
func (slf *Config) IncrAssets(_key, _name string, _buf float64) error {
	conn := config.client
	err := conn.HIncrByFloat(_key, _name, _buf).Err()
	if err != nil && err != redis.Nil {
		log.Printf("redis HIncrByFloat err:%v\n", err)
		return err
	}
	return nil
}

// RedisIncr redis Incr
func (slf *Config) IncrByFloat(_key string, _buf float64) float64 {
	conn := config.client
	code, err := conn.IncrByFloat(_key, _buf).Result()
	if err != nil && err != redis.Nil {
		log.Printf("redis RedisIncr err: %v\n", err)
	}
	return code
}

type PipelineParam struct {
	Amount decimal.Decimal `json:"amount"` // 锁定或解锁金额
	SubKey string          `json:"subKey"` // 目标 Hash key
}

// AssetPipeline opt user asset through Pipeline
func (slf *Config) AssetTxPipeline(_key string, _pipelineParam ...PipelineParam) ([]decimal.Decimal, error) {
	conn := config.client
	result := make([]decimal.Decimal, 0)
	cmdList := make([]*redis.IntCmd, 0)
	pipe := conn.TxPipeline()
	for _, param := range _pipelineParam {
		//str := pipe.HGet(_key, param.SubKey)
		//log.Printf(str.String())
		incrCmd := pipe.HIncrBy(_key, param.SubKey, param.Amount.Mul(decimal.New(1e8, 0)).IntPart())
		cmdList = append(cmdList, incrCmd)
	}
	_, err := pipe.Exec()
	if err != nil {
		/***手动进行 pipeline rollback--start***/
		pipe := conn.TxPipeline()
		cmdLen := len(cmdList)                       // pipeline 命令执行数组 (包含成功或失败)
		revertedPipeline := make([]PipelineParam, 0) // 手动构建 pipeline rollback 数组初始化
		for i := 0; i < cmdLen; i++ {
			// 针对 pipeline 中执行成功的命令构建 rollback 数组
			cmd := cmdList[i]
			if cmd.Err() == nil && len(cmd.Args()) > 3 {
				intValue, ok := cmd.Args()[3].(int64)
				if ok {
					_pipelineParam[i].Amount = decimal.New(intValue, 0).Mul(decimal.New(-1, 0))
					revertedPipeline = append(revertedPipeline, _pipelineParam[i])
				}
			}
		}
		if len(revertedPipeline) > 0 {
			// 组合 pipeline rollback 命令
			for _, param := range revertedPipeline {
				amountInt64 := param.Amount.IntPart()
				_ = pipe.HIncrBy(_key, param.SubKey, amountInt64)
			}
			// 执行
			_, _ = pipe.Exec()
		}
		/***手动进行 pipeline rollback--end***/

		// 返回用户操作 pipeline 中的错误(注意:不要和自己构建的 pipeline rollback 混淆)
		log.Printf("[redis error] %v\n", err)
		return nil, err
	}
	for _, cmd := range cmdList {
		//log.Printf("%d", cmd.Val())
		log.Printf("%f", float64(cmd.Val()*1.0)/1e8)
		log.Printf("%v", decimal.NewFromFloat(float64(cmd.Val()/(1e8*1.0))))
		result = append(result, decimal.NewFromFloat(float64(cmd.Val())/1e8))
	}

	return result, err
}

// HMGet get user c2c asset
func (slf *Config) HMGet(_key string, _subKeys ...string) (map[string]decimal.Decimal, error) {
	conn := config.client
	result := make(map[string]decimal.Decimal)
	sliceRes, err := conn.HMGet(_key, _subKeys...).Result()
	if err != nil {
		log.Printf("redis TxPipeline exec err: %v\n", err)
		return nil, err
	}
	for index, value := range sliceRes {
		valueStr, ok := value.(string)
		if ok {
			valDecimal, err := decimal.NewFromString(valueStr)
			log.Printf("valDecimal: %v\n", valDecimal)
			if err != nil {
				continue
			}
			result[_subKeys[index]] = valDecimal.Div(decimal.New(1e8, 0))
		}
	}

	return result, nil
}

// HGet get user c2c asset
func (slf *Config) HGetAll(_key string) (map[string]string, error) {
	conn := config.client
	sliceRes, err := conn.HGetAll(_key).Result()
	if err != nil {
		log.Printf("redis TxPipeline exec err: %v\n", err)
		return nil, err
	}

	return sliceRes, nil
}

// HGet get user c2c asset
func (slf *Config) Keys(_key string) ([]string, error) {
	conn := config.client
	sliceRes, err := conn.Keys(_key).Result()
	if err != nil {
		log.Printf("redis TxPipeline exec err: %v\n", err)
		return nil, err
	}

	return sliceRes, nil
}

// HGet get user c2c asset
func (slf *Config) HExist(_key string) error {
	conn := config.client
	_, err := conn.Exists(_key).Result()
	if err != nil {
		log.Printf("redis TxPipeline exec err: %v\n", err)
		return err
	}

	return nil
}

// HGet get asset
func (slf *Config) HSet(hash, key, val string) error {
	conn := config.client
	err := conn.HSet(hash, key, val).Err()
	return err
}

// HGet get asset
func (slf *Config) HGet(hash, _key string) (string, error) {
	conn := config.client
	res, err := conn.HGet(hash, _key).Result()
	return res, err
}

// HGet get user c2c asset
func (slf *Config) DelKey(_key string) (int64, error) {
	conn := config.client
	sliceRes, err := conn.Del(_key).Result()
	if err != nil {
		log.Printf("redis TxPipeline exec err: %v\n", err)
		return -1, err
	}

	return sliceRes, nil
}

// HGet get user c2c asset
func (slf *Config) KeysByNode(_key string) ([]string, error) {
	conn := config.client
	result := make(map[string]uint64, 0)

	_ = conn.ForEachNode(func(client *redis.Client) error {
		var cursor uint64
		iter := client.Scan(cursor, _key, 5000).Iterator()

		for iter.Next() {
			result[iter.Val()] = cursor
		}
		if err := iter.Err(); err != nil {
			log.Printf("iter key err")
			return err
		}

		return nil
	})

	keys := make([]string, len(result))

	i := 0
	for k := range result {
		keys[i] = k
		i++
	}
	return keys, nil
}

// TxPipelineAsset opt user asset through TxPipeline
func (slf *Config) TxPipelineAsset(_key string, _pipelineParam ...PipelineParam) ([]decimal.Decimal, error) {
	conn := config.client
	result := make([]decimal.Decimal, 0)
	cmdList := make([]*redis.IntCmd, 0)
	pipe := conn.TxPipeline()
	for _, param := range _pipelineParam {
		amountInt64 := param.Amount.Mul(decimal.New(1e8, 0)).IntPart()
		incrCmd := pipe.HIncrBy(_key, param.SubKey, amountInt64)
		cmdList = append(cmdList, incrCmd)
	}
	_, err := pipe.Exec()
	if err != nil {
		log.Printf("redis TxPipeline exec err: %v\n", err)
		return nil, err
	}
	for _, cmd := range cmdList {
		result = append(result, decimal.NewFromFloat(float64(cmd.Val())/1e8))
	}

	return result, err
}
