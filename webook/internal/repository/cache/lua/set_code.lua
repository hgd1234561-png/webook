--你的验证码在Redis上的key
-- phone_code:login:152xxxxxxx
local key = KEYS[1]
-- 验证次数，我们一个验证码，最多重复三次，这个记录了几次
-- phone_code:login:152xxxxxxx:cnt
local cntKey = key .. ":cnt"
-- 你的验证码 123456
local val = ARGV[1]

--过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    --    key 存在，但是没有过期时间
    return -2
    -- 540 = 600-60 九分钟
elseif ttl == -2 or ttl < 540 then
    --    可以发验证码
    redis.call("set", key, val)
    -- 600 秒
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end