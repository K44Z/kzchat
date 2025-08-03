-- name: GetChatMessagesByChatId :many
SELECT messages.*
FROM messages
WHERE messages.chat_id = $1
ORDER BY time ASC;

-- name: GetDmChatMessagesByParticipants :many
SELECT m.*
FROM messages m
JOIN chats c ON m.chat_id = c.id
WHERE c.type = 'dm'
  AND m.chat_id IN (
    SELECT chat_id
    FROM chat_members
    WHERE user_id = $1 OR user_id = $2
    GROUP BY chat_id
    HAVING COUNT(DISTINCT user_id) = 2
  )
ORDER BY m.time ASC;


-- name: StoreChatMessage :execresult
INSERT INTO messages(sender_id, content, chat_id, time, type)
VALUES($1, $2, $3, $4, $5);

-- name: CreateChat :one
INSERT INTO chats("type", name)
VALUES($1, $2)
RETURNING *;

-- name: CreateChatMembers :execresult
INSERT INTO chat_members(chat_id, user_id)
VALUES($1, $2);

-- name: FindChatByParticipants :one
SELECT chat_id
FROM chat_members
WHERE user_id = ANY($1::int[])
GROUP BY chat_id
HAVING COUNT(*) = $2
   AND COUNT(*) = (
       SELECT COUNT(*) FROM chat_members cm2
       WHERE cm2.chat_id = chat_members.chat_id
   );

-- name: GetChatById :one
SELECT * FROM chats
WHERE id = $1;
