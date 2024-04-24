class User {
  final String id;
  final String nickname;
  final String avatarUrl;
  final List<String> following;
  final List<String> followers;
  final List<String> blocked;
  final int coins;
  final List<String> coinTransactions;

  User({
    required this.id,
    required this.nickname,
    required this.avatarUrl,
    required this.following,
    required this.followers,
    required this.blocked,
    required this.coins,
    required this.coinTransactions,
  });

  // 从JSON数据创建User对象的工厂方法
  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      nickname: json['nickname'],
      avatarUrl: json['avatarUrl'],
      following: List<String>.from(json['following']),
      followers: List<String>.from(json['followers']),
      blocked: List<String>.from(json['blocked']),
      coins: json['coins'],
      coinTransactions: List<String>.from(json['coinTransactions']),
    );
  }

  // 将User对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'nickname': nickname,
      'avatarUrl': avatarUrl,
      'following': following,
      'followers': followers,
      'blocked': blocked,
      'coins': coins,
      'coinTransactions': coinTransactions,
    };
  }

  // 方法：关注另一个用户
  void followUser(String userId) {
    if (!following.contains(userId)) {
      following.add(userId);
    }
  }

  // 方法：取消关注另一个用户
  void unfollowUser(String userId) {
    following.remove(userId);
  }

  // 方法：屏蔽用户
  void blockUser(String userId) {
    if (!blocked.contains(userId)) {
      blocked.add(userId);
    }
  }

  // 方法：解除屏蔽用户
  void unblockUser(String userId) {
    blocked.remove(userId);
  }
}
