import 'package:aorb/models/comment.dart';
// 单条投票记录，包含在question中
class Vote {
  final String id;
  final String type; // 这里我们使用 String 替代原 Type 枚举
  final String time;
  final String user_id;
  final String title;
  final String description;
  final List<String> options;
  final String channel;
  final List<Comment> comments;
  final int? fee; // 用 int? 表示可选的整数
  final List<String> inviteIDs;
  final List<String> user_ids;

  Vote({
    required this.id,
    required this.channel,
    required this.comments,
    required this.description,
    this.fee,
    required this.user_id,
    required this.time,
    required this.title,
    required this.options,
    required this.inviteIDs,
    required this.type,
    required this.user_ids,
  });

  // 从JSON数据创建Vote对象的工厂方法
  factory Vote.fromJson(Map<String, dynamic> json) {
    return Vote(
      id: json['id'],
      description: json['description'],
      channel: json['channel'],
      comments: List<Comment>.from(json['comments'].map((x) => Comment.fromJson(x))),
      fee: json['fee'],
      user_id: json['user_id'],
      time: json['time'],
      title: json['title'],
      options: List<String>.from(json['options']),
      inviteIDs: List<String>.from(json['invite_ids'] ?? []),
      type: json['type'],
      user_ids: List<String>.from(json['user_ids']),
    );
  }

  // 将Vote对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'channel': channel,
      'comments': List<dynamic>.from(comments.map((x) => x.toJson())),
      'fee': fee,
      'user_id': user_id,
      'time': time,
      'title': title,
      'description': description,
      'options': options,
      'invite_ids': inviteIDs,
      'type': type,
      'user_ids': user_ids,
    };
  }
}