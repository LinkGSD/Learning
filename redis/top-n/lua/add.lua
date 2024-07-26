-- 排行榜 key
local key = KEYS[1]
-- 要更新的用户 ID
local uid = ARGV[1]
-- 用户本次新增的 val （小数位为时间差标识）
local valScore = ARGV[2]
-- 获取用户之前的 score
local score = redis.call("ZSCORE", key, uid)
if score == false then
    score = 0
end
-- 从 score 中抹除用于时间差标识的小数部分，获取整数的排序 val
local val = math.floor(score)
-- 更新用户最新的 score 信息（累计 val.最新时间差）
local newScore = valScore+val
redis.call("ZADD", key, newScore, uid)
-- 更新成功返回 newScore（注意要使用 tostring 才能返回小数）
return tostring(newScore)