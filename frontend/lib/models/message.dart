class Message {
  final String id;
  final String fromUserId;
  final String toUserId;
  final String content;
  final DateTime timestamp;
  bool isRead;  // 消息是否已读

  Message({
    required this.id,
    required this.fromUserId,
    required this.toUserId,
    required this.content,
    required this.timestamp,
    this.isRead = false,  // 默认未读
  });

  // 从JSON数据创建Message对象的工厂方法
  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      id: json['id'],
      fromUserId: json['fromUserId'],
      toUserId: json['toUserId'],
      content: json['content'],
      timestamp: DateTime.parse(json['timestamp']),
      isRead: json['isRead'] ?? false,  // 安全处理，如果无此字段则默认为false
    );
  }

  // 将Message对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'fromUserId': fromUserId,
      'toUserId': toUserId,
      'content': content,
      'timestamp': timestamp.toIso8601String(),
      'isRead': isRead,
    };
  }

  // 标记消息为已读
  void markAsRead() {
    isRead = true;
  }
}
