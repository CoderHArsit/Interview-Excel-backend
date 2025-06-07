package utils

import (
	"interviewexcel-backend-go/config"
	"time"
)

//AddTokenToBlacklist adds a token to the Redis blacklist with an expiration

func AddTokenToBlacklist(token string,expiration time.Duration)error{
	return config.RedisClient.Set(config.Ctx,token,true,expiration).Err()
}


//IsTokenBlacklisted checks if a token is blacklisted

func IsTokenBlacklisted(token string)(bool,error){
	exists,err:=config.RedisClient.Exists(config.Ctx,token).Result()
	if err != nil{
			return false,err
		}
	
	return exists>0, nil
}