class User {
  final String id;
  final String nickname;
  final String avatarUrl;
  final List<String> followed;
  final List<String> followers;
  final List<String> blacklist;
  final int coins;
  final List<String> coinsRecord;
  final List<String> questionsAsk;
  final List<String> questionsAnswer;
  final List<String> channels;

  User({
    required this.id,
    required this.nickname,
    required this.avatarUrl,
    required this.followed,
    required this.followers,
    required this.blacklist,
    required this.coins,
    required this.coinsRecord,
    required this.questionsAsk,
    required this.questionsAnswer,
    required this.channels,
  });

  // 从JSON数据创建User对象的工厂方法
  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      nickname: json['nickname'],
      avatarUrl: json['avatar'],
      followed: List<String>.from(json['followed']),
      followers: List<String>.from(json['follower']),
      blacklist: List<String>.from(json['blacklist']),
      coins: json['coins'].toInt(),
      coinsRecord: List<String>.from(json['coins_record'] ?? []),
      questionsAsk: List<String>.from(json['questions_ask']),
      questionsAnswer: List<String>.from(json['questions_asw']),
      channels: List<String>.from(json['channels']),
    );
  }

  // 将User对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'nickname': nickname,
      'avatar': avatarUrl,
      'followed': followed,
      'follower': followers,
      'blacklist': blacklist,
      'coins': coins,
      'coins_record': coinsRecord,
      'questions_ask': questionsAsk,
      'questions_asw': questionsAnswer,
      'channels': channels,
    };
  }
}
