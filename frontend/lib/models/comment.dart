class Comment {
  final String id;
  final String userId;
  final String voteId;
  final String content;
  final DateTime timestamp;
  final String userAvatarUrl; // 用户头像URL
  final String userName; // 用户昵称

  Comment({
    required this.id,
    required this.userId,
    required this.voteId,
    required this.content,
    required this.timestamp,
    required this.userAvatarUrl,
    required this.userName,
  });

  // 从JSON创建Comment对象的工厂方法
  factory Comment.fromJson(Map<String, dynamic> json) {
    return Comment(
      id: json['id'],
      userId: json['userId'],
      voteId: json['voteId'],
      content: json['content'],
      timestamp: DateTime.parse(json['timestamp']),
      userAvatarUrl: json['userAvatarUrl'],
      userName: json['userName'],
    );
  }

  // 将Comment对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'userId': userId,
      'voteId': voteId,
      'content': content,
      'timestamp': timestamp.toIso8601String(),
      'userAvatarUrl': userAvatarUrl,
      'userName': userName,
    };
  }
}
