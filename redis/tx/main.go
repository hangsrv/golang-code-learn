package main

import (
	"context"
	"fmt"
	"golang-code-learn/redis/common"

	"github.com/go-redis/redis/v8"
)

const str = `	
local key1 = KEYS[1]
local key2 = KEYS[2]
local new_value1 = ARGV[1]
local new_value2 = ARGV[2]

local current_value1 = redis.call('GET', key1)
local current_value2 = redis.call('GET', key2)

if current_value1 == false then
	current_value1 = ""
end
if current_value2 == false then
	current_value2 = ""
end

if current_value1 == ARGV[3] and current_value2 == ARGV[4] then
	local mset_result = redis.call('MSET', key1, new_value1, key2, new_value2)

	if mset_result == false then  
		return redis.error_reply('Failed to set all keys')
	else  
		return "All keys are set successfully"
	end
else
    return redis.error_reply('Concurrent modification detected')  
end
`

func main() {
	c := common.NewClient()
	s := redis.NewScript(str)

	res, err := s.Eval(context.Background(), c, []string{"k1", "k2"}, "v1", "v2", "v1", "v2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

}
