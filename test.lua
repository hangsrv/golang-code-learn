-- =============================================
-- 检查列表长度
local currentLength = redis.call('LLEN', KEYS[1])

-- 判断长度是否在 0 和 100 之间
if currentLength < 100 then
    -- 计算实际要 PUSH 的元素数量，最多 PUSH 20 个
    local pushCount = math.min(20, 100 - currentLength)
    
    -- 准备要 PUSH 的元素
    local elements = {}
    for i = 1, pushCount do
        table.insert(elements, ARGV[i])
    end
    
    -- 执行 LPUSH 操作
    redis.call('LPUSH', KEYS[1], unpack(elements))
    
    return pushCount  -- 返回实际 PUSH 的元素数量
else
    return 0  -- 如果列表已满，返回 0
end

-- =============================================
-- 检查 key1 的长度是否大于或等于 100
local currentLength = redis.call('LLEN', KEYS[1])

if currentLength >= 100 then
    -- 弹出 100 个元素并推送到 key2 中
    local elements = redis.call('LRANGE', KEYS[1], 0, 99)
    redis.call('RPUSH', KEYS[2], unpack(elements))
    redis.call('LTRIM', KEYS[1], 100, -1)
    
    return 1 -- 返回 1 表示操作成功
else
    return 0 -- 返回 0 表示 key1 长度不足
end
