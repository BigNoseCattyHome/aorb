// 连接到 MongoDB 服务
conn = new Mongo();

// 切换到数据库 aorb
db = conn.getDB("aorb");

// 检查是否存在数据库 aorb，如果存在则清空所有集合
if (db.getCollectionInfos().length > 0) {
  db.getCollectionNames().forEach(function (collectionName) {
    db[collectionName].drop();
  });
  print("Database 'aorb' has been cleared.");
}

// 创建 users 集合并插入数据
db.createCollection("users");
db.createCollection("userCounter");
db.userCounter.insertOne(
    {
      _id: "userId",
      sequence_value: 0
    }
);

db.users.insertMany([
  {
    _id: getNextUserSequence("userId"),
    username: "aichifan",
    nickname: "爱吃饭的小袁同学",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/05/27/2MgJcvLtOVKmAdn.jpg",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextUserSequence("userId"),
    username: "gopher",
    nickname: "花枝鼠gogo来帮忙",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextUserSequence("userId"),
    username: "siri",
    nickname: "风见澈Siri",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/05/27/QzKM41C3Vs5FeHW.jpg",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextUserSequence("userId"),
    username: "anti_cris",
    nickname: "Anti Cris",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/05/27/alt3BKPYhzmV4E7.jpg",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextUserSequence("userId"),
    username: "sirius",
    nickname: "Sirius Wild",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/06/07/newUserA.jpg",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  }
]);

// 创建 polls 集合并插入数据
db.createCollection("polls");
db.createCollection("pollCounter");
db.pollCounter.insertOne(
    {
      _id: "pollId",
      sequence_value: 0
    }
);

db.polls.insertMany([
  {
    _id: getNextPollSequence("pollId"),
    poll_type: "public",
    title: "午饭吃什么呀?",
    options: ["麻辣烫", "炸鸡汉堡"],
    options_count: [0.4, 0.6],
    user_id: 1,
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextPollSequence("pollId"),
    poll_type: "public",
    title: "下午去哪里玩?",
    options: ["顾村公园", "外滩"],
    options_count: [0.16, 0.84],
    user_id: 2,
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextPollSequence("pollId"),
    poll_type: "public",
    title: "要不要去小美家玩啊？",
    options: ["麻辣烫", "炸鸡汉堡"],
    options_count: [0.4, 0.6],
    user_id: 3,
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextPollSequence("pollId"),
    poll_type: "public",
    title: "Exploring the Enigmatic World of Quantum Mechanics",
    options: ["Yes, I'd love to.", "I'm familiar with quantum mechanics."],
    options_count: [0.4, 0.6],
    user_id: 4,
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  }
]);

// 创建 comments 集合并插入数据
db.createCollection("comments");
db.createCollection("commentCounter");
db.commentCounter.insertOne(
    {
      _id: "commentId",
      sequence_value: 0
    }
);

db.comments.insertMany([
  {
    _id: getNextCommentSequence("commentId"),
    user_id: 2,
    poll_id: 1,
    content: "我觉得麻辣烫不错",
    create_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextCommentSequence("commentId"),
    user_id: 3,
    poll_id: 1,
    content: "我还是喜欢炸鸡汉堡",
    create_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextCommentSequence("commentId"),
    user_id: 1,
    poll_id: 2,
    content: "外滩风景更好",
    create_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextCommentSequence("commentId"),
    user_id: 2,
    poll_id: 3,
    content: "麻辣烫好吃",
    create_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextCommentSequence("commentId"),
    user_id: 3,
    poll_id: 3,
    content: "炸鸡汉堡更香",
    create_at: new Date(),
    delete_at: null,
  }
]);

// 创建 messages 集合并插入数据
// db.createCollection("messages");
// db.createCollection("messageCounter");
// db.messages.insertMany([
//   {
//     _id: ObjectId(),
//     from_user_id: "1",
//     to_user_id: "2",
//     type: "comment",
//     content: "comment1",
//     time: new Date().toISOString(),
//     status: "unread"
//   },
//   {
//     _id: ObjectId(),
//     from_user_id: "2",
//     to_user_id: "3",
//     type: "comment",
//     content: "comment2",
//     time: new Date().toISOString(),
//     status: "unread"
//   },
//   {
//     _id: ObjectId(),
//     from_user_id: "3",
//     to_user_id: "1",
//     type: "comment",
//     content: "comment3",
//     time: new Date().toISOString(),
//     status: "read"
//   },
//   {
//     _id: ObjectId(),
//     from_user_id: "2",
//     to_user_id: "1",
//     type: "comment",
//     content: "comment4",
//     time: new Date().toISOString(),
//     status: "unread"
//   }
// ]);

// 创建 votes 集合并插入数据
db.createCollection("votes");
db.createCollection("voteCounter");
db.voteCounter.insertOne(
    {
      _id: "voteId",
      sequence_value: 0
    }
);
db.votes.insertMany([
  {
    _id: getNextVoteSequence("voteId"),
    poll_id: 1,
    user_id: 2,
    choice: "麻辣烫",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextVoteSequence("voteId"),
    poll_id: 1,
    user_id: 3,
    choice: "炸鸡汉堡",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextVoteSequence("voteId"),
    poll_id: 2,
    user_id: 1,
    choice: "外滩",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextVoteSequence("voteId"),
    poll_id: 3,
    user_id: 2,
    choice: "麻辣烫",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextVoteSequence("voteId"),
    poll_id: 3,
    user_id: 3,
    choice: "炸鸡汉堡",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: getNextVoteSequence("voteId"),
    poll_id: 4,
    user_id: 1,
    choice: "Yes, I'd love to.",
    create_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  }
]);

// db.createCollection("refresh_tokens");
