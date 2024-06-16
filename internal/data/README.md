# Queries for Interesting Stats

### Top users by most messages sent
```
SELECT
    Users.discordTag,
    UserStats.userId,
    UserStats.messagesSent
FROM
    UserStats
JOIN Users ON UserStats.userId = Users.userId
ORDER BY
    UserStats.messagesSent DESC
```

### Top users by most reactions received
```
SELECT
    Users.discordTag,
    UserStats.userId,
    UserStats.reactionsReceived
FROM
    UserStats
JOIN Users ON UserStats.userId = Users.userId
ORDER BY
    UserStats.reactionsReceived DESC
```