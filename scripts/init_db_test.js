// 连接到 MongoDB 服务
conn = new Mongo();

// 切换到数据库 aorb
db = conn.getDB("aorb");

// 检查是否存在数据库 aorb，如果存在则清空所有集合
if (db.getCollectionInfos().length > 0) {
    db.getCollectionNames().forEach(function (collectionName) {
        db[collectionName].drop();
    });
    print("Database 'aorb_test' has been cleared.");
}



// insert users
db.createCollection("users")
for (let i = 1; i <= 6; i++) {
    db.users.insertOne({
        id: ObjectId(),
        username: "user" + i,
        avatar: "https://s2.loli.net/2024/08/01/RiYJwlW5kNzDnPq.jpg",
        nickname: "nickname" + i,
        createat: {
            seconds: Math.floor(Date.now() / 1000) + i,
            nanos: 293536000 + i
        },
        updateat: {
            seconds: Math.floor(Date.now() / 1000) + i,
            nanos: 293713000 + i
        },
        deleteat: {
            seconds: -62135596800,
            nanos: 0
        },
        password: "password" + i,
        coins: i * 10,
        coinsrecord: {
            records: []
        },
        followed: {
            usernames: [
                "user" + ((i + 1) % 20 + 1),
                "user" + ((i + 2) % 20 + 1),
            ]
        },
        follower: {
            usernames: [
                "user" + ((i + 3) % 20 + 1),
                "user" + ((i + 4) % 20 + 1),
            ]
        },
        blacklist: {
            usernames: [
                "user" + ((i + 5) % 20 + 1),
                "user" + ((i + 6) % 20 + 1),
            ]
        },
        ipaddress: "Shanghai",
        pollask: {
            pollids: []
        },
        pollans: {
            pollids: []
        },
        pollcollect: {
            pollids: []
        }
    });
}
for (let i = 7; i <= 14; i++) {
    db.users.insertOne({
        id: ObjectId(),
        username: "user" + i,
        avatar: "https://s2.loli.net/2024/08/01/UNWt8EDQGMzY2id.jpg",
        nickname: "nickname" + i,
        createat: {
            seconds: Math.floor(Date.now() / 1000) + i,
            nanos: 293536000 + i
        },
        updateat: {
            seconds: Math.floor(Date.now() / 1000) + i,
            nanos: 293713000 + i
        },
        deleteat: {
            seconds: -62135596800,
            nanos: 0
        },
        password: "password" + i,
        coins: i * 10,
        coinsrecord: {
            records: []
        },
        followed: {
            usernames: [
                "user" + ((i + 1) % 20 + 1),
                "user" + ((i + 2) % 20 + 1),
            ]
        },
        follower: {
            usernames: [
                "user" + ((i + 3) % 20 + 1),
                "user" + ((i + 4) % 20 + 1),
            ]
        },
        blacklist: {
            usernames: [
                "user" + ((i + 5) % 20 + 1),
                "user" + ((i + 6) % 20 + 1),
            ]
        },
        ipaddress: "Shanghai",
        pollask: {
            pollids: []
        },
        pollans: {
            pollids: []
        },
        pollcollect: {
            pollids: []
        }
    });
}
for (let i = 15; i <= 20; i++) {
    db.users.insertOne({
        id: ObjectId(),
        username: "user" + i,
        avatar: "https://s2.loli.net/2024/08/01/p8YwiSyQGeBMXqZ.jpg",
        nickname: "nickname" + i,
        createat: {
            seconds: Math.floor(Date.now() / 1000) + i,
            nanos: 293536000 + i
        },
        updateat: {
            seconds: Math.floor(Date.now() / 1000) + i,
            nanos: 293713000 + i
        },
        deleteat: {
            seconds: -62135596800,
            nanos: 0
        },
        password: "password" + i,
        coins: i * 10,
        coinsrecord: {
            records: []
        },
        followed: {
            usernames: [
                "user" + ((i + 1) % 20 + 1),
                "user" + ((i + 2) % 20 + 1),
            ]
        },
        follower: {
            usernames: [
                "user" + ((i + 3) % 20 + 1),
                "user" + ((i + 4) % 20 + 1),
            ]
        },
        blacklist: {
            usernames: [
                "user" + ((i + 5) % 20 + 1),
                "user" + ((i + 6) % 20 + 1),
            ]
        },
        ipaddress: "Shanghai",
        pollask: {
            pollids: []
        },
        pollans: {
            pollids: []
        },
        pollcollect: {
            pollids: []
        }
    });
}

const polls = [];

for (let i = 1; i <= 10; i++) {
    const poll = {
        commentList: [],
        createAt: new ISODate(`2023-07-15T10:30:00.000Z`),
        options: ["选项A", "选项B"],
        optionsCount: [i * 2, i * 3],
        pollType: "public",
        pollUuid: `dfbeb25d-f79d-4003-aafb-995ea6fe345${i}`,
        title: `投票标题${i}`,
        userName: `user${i}`,
        voteList: []
    };

    for (let j = 1; j <= 10; j++) {
        poll.commentList.push({
            commentUuid: `06ec93af-7b22-4fde-8061-e91b87c5f79${j}`,
            commentUserName: `user${j}`,
            content: `评论内容${j}`,
            createAt: new ISODate(`2023-07-15T10:30:00.000Z`)
        });

        poll.voteList.push({
            voteUuid: `79421094-de8e-44e6-bb03-adc381a0b5a${j}`,
            voteUserName: `user${j}`,
            choice: j % 2 === 0 ? "选项A" : "选项B",
            createAt: new ISODate(`2023-07-15T10:30:00.000Z`)
        });
    }

    polls.push(poll);
}

db.polls.insertMany(polls);