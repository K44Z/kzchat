-- name: GetChatMessagesByChatId :many
SELECT messages.*
FROM messages
WHERE messages.chat_id = $1
ORDER BY time ASC;

-- name: GetDmChatMessagesByParticipants :many
SELECT
    m.id,
    m.content,
    m.chat_id,
    m.time,
    m.type,
    u_sender.id AS sender_id,
    u_sender.username AS sender_username,
    u_receiver.id AS receiver_id,
    u_receiver.username AS receiver_username
FROM messages m
JOIN chats c
    ON m.chat_id = c.id
JOIN users u_sender
    ON m.sender_id = u_sender.id
JOIN chat_members cm_receiver
    ON cm_receiver.chat_id = m.chat_id
    AND cm_receiver.user_id != m.sender_id
JOIN users u_receiver
    ON cm_receiver.user_id = u_receiver.id
WHERE c.type = 'dm'
  AND m.chat_id = (
      SELECT cm1.chat_id
      FROM chat_members cm1
      WHERE cm1.user_id IN ($1, $2)
      GROUP BY cm1.chat_id
      HAVING COUNT(DISTINCT cm1.user_id) = 2
         AND COUNT(*) = 2
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
SELECT cm.chat_id
FROM chat_members cm
WHERE cm.user_id = ANY($1::int[])
GROUP BY cm.chat_id
HAVING COUNT(*) = 2
   AND COUNT(*) = (
       SELECT COUNT(*)
       FROM chat_members cm2
       WHERE cm2.chat_id = cm.chat_id
   );

-- name: GetChatById :one
SELECT * FROM chats
WHERE id = $1;
