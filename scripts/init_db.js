// 连接到 MongoDB 服务
conn = new Mongo();

// 切换到数据库 aorb
db = conn.getDB("aorb");

// 创建 users 集合并插入数据
db.createCollection("users");
db.users.insertMany([
  {
    _id: ObjectId(),
    nickname: "爱吃饭的小袁同学",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/05/27/2MgJcvLtOVKmAdn.jpg",
    coins: 100,
    coins_record: [],
    followed: ["2", "3"],
    follower: ["2"],
    blacklist: [],
    questions_ask: ["1"],
    questions_asw: [],
    questions_collect: ["2"],
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    nickname: "花枝鼠gogo来帮忙",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png",
    coins: 50,
    coins_record: [],
    followed: ["1"],
    follower: ["1", "3"],
    blacklist: [],
    questions_ask: ["2"],
    questions_asw: [],
    questions_collect: ["1"],
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    nickname: "风见澈Siri",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/05/27/QzKM41C3Vs5FeHW.jpg",
    coins: 70,
    coins_record: [],
    followed: ["1", "2"],
    follower: [],
    blacklist: [],
    questions_ask: ["3"],
    questions_asw: [],
    questions_collect: ["2"],
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    nickname: "Anti Cris",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/05/27/alt3BKPYhzmV4E7.jpg",
    coins: 90,
    coins_record: [],
    followed: [],
    follower: [],
    blacklist: [],
    questions_ask: ["4"],
    questions_asw: [],
    questions_collect: [],
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    nickname: "Sirius Wild",
    password: "e10adc3949ba59abbe56e057f20f883e",
    avatar: "https://s2.loli.net/2024/06/07/newUserA.jpg",
    coins: 120,
    coins_record: [],
    followed: [],
    follower: [],
    blacklist: [],
    questions_ask: [],
    questions_asw: [],
    questions_collect: [],
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  }
]);

// 创建 polls 集合并插入数据
db.createCollection("polls");
db.polls.insertMany([
  {
    _id: ObjectId(),
    type: "public",
    title: "午饭吃什么呀?",
    options: ["麻辣烫", "炸鸡汉堡"],
    options_rate: [0.4, 0.6],
    sponsor: "爱吃饭的小袁同学",
    votes: ["vote1", "vote2"],
    collections: ["1"],
    comments: ["comment1", "comment2"],
    fee: 10,
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    type: "public",
    title: "下午去哪里玩?",
    options: ["顾村公园", "外滩"],
    options_rate: [0.16, 0.84],
    sponsor: "花枝鼠gogo来帮忙",
    votes: ["vote3"],
    collections: ["1", "3"],
    comments: ["comment3"],
    fee: 20,
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    type: "public",
    title: "要不要去小美家玩啊？",
    options: ["麻辣烫", "炸鸡汉堡"],
    options_rate: [0.4, 0.6],
    sponsor: "风见澈Siri",
    votes: ["vote4", "vote5"],
    collections: ["2"],
    comments: ["comment4"],
    fee: 15,
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    type: "public",
    title: "Exploring the Enigmatic World of Quantum Mechanics",
    options: ["Yes, I'd love to.", "I'm familiar with quantum mechanics."],
    options_rate: [0.4, 0.6],
    sponsor: "Anti Cris",
    votes: ["vote6"],
    collections: [],
    comments: [],
    fee: 30,
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  }
]);

// 创建 comments 集合并插入数据
db.createCollection("comments");
db.comments.insertMany([
  {
    _id: ObjectId(),
    user_id: "2",
    poll_id: "1",
    reply_id: "",
    vote: "vote1",
    content: "我觉得麻辣烫不错",
    creare_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    user_id: "3",
    poll_id: "1",
    reply_id: "",
    vote: "vote2",
    content: "我还是喜欢炸鸡汉堡",
    creare_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    user_id: "1",
    poll_id: "2",
    reply_id: "",
    vote: "vote3",
    content: "外滩风景更好",
    creare_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    user_id: "2",
    poll_id: "3",
    reply_id: "",
    vote: "vote4",
    content: "麻辣烫好吃",
    creare_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    user_id: "3",
    poll_id: "3",
    reply_id: "",
    vote: "vote5",
    content: "炸鸡汉堡更香",
    creare_at: new Date(),
    delete_at: null,
  }
]);

// 创建 messages 集合并插入数据
db.createCollection("messages");
db.messages.insertMany([
  {
    _id: ObjectId(),
    from_user_id: "1",
    to_user_id: "2",
    type: "comment",
    content: "comment1",
    time: new Date().toISOString(),
    status: "unread"
  },
  {
    _id: ObjectId(),
    from_user_id: "2",
    to_user_id: "3",
    type: "comment",
    content: "comment2",
    time: new Date().toISOString(),
    status: "unread"
  },
  {
    _id: ObjectId(),
    from_user_id: "3",
    to_user_id: "1",
    type: "comment",
    content: "comment3",
    time: new Date().toISOString(),
    status: "read"
  },
  {
    _id: ObjectId(),
    from_user_id: "2",
    to_user_id: "1",
    type: "comment",
    content: "comment4",
    time: new Date().toISOString(),
    status: "unread"
  }
]);

// 创建 votes 集合并插入数据
db.createCollection("votes");
db.votes.insertMany([
  {
    _id: ObjectId(),
    poll_id: "1",
    voter: "2",
    choice: "麻辣烫",
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    poll_id: "1",
    voter: "3",
    choice: "炸鸡汉堡",
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    poll_id: "2",
    voter: "1",
    choice: "外滩",
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    poll_id: "3",
    voter: "2",
    choice: "麻辣烫",
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    poll_id: "3",
    voter: "3",
    choice: "炸鸡汉堡",
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  },
  {
    _id: ObjectId(),
    poll_id: "4",
    voter: "1",
    choice: "Yes, I'd love to.",
    creare_at: new Date(),
    update_at: new Date(),
    delete_at: null,
  }
]);
