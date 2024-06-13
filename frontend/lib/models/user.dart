class User {
  final String id;
  final String nickname;
  final String avatar;
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
    required this.avatar,
    required this.followed,
    required this.followers,
    required this.blacklist,
    required this.coins,
    required this.coinsRecord,
    required this.questionsAsk,
    required this.questionsAnswer,
    required this.channels,
  });

  // 从JSON数据创建User对象的工厂方法，提供默认值，可以防止转换的时候出现字段缺失导致的错误
  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      nickname: json['nickname'] ?? '',
      avatar: json['avatar'] ?? '',
      followed: json['followed'] != null ? List<String>.from(json['followed']) : [],
      followers: json['follower'] != null ? List<String>.from(json['follower']) : [],
      blacklist: json['blacklist'] != null ? List<String>.from(json['blacklist']) : [],
      coins: json['coins'] != null ? json['coins'].toInt() : 0,
      coinsRecord: json['coins_record'] != null ? List<String>.from(json['coins_record']) : [],
      questionsAsk: json['questions_ask'] != null ? List<String>.from(json['questions_ask']) : [],
      questionsAnswer: json['questions_asw'] != null ? List<String>.from(json['questions_asw']) : [],
      channels: json['channels'] != null ? List<String>.from(json['channels']) : [],
    );
  }

  // 将User对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'nickname': nickname,
      'avatar': avatar,
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
