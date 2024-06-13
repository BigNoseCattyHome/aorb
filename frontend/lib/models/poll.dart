import 'package:aorb/models/comment.dart';
import 'package:aorb/models/vote.dart';

class Poll {
  final String id;
  final String type; // 这里我们使用 String 替代原 Type 枚举
  final String title;
  final String description;
  final List<String> options;
  final List<double> options_rate;
  final DateTime time;
  final String ipaddress;
  final String sponsor;
  final List<Vote> votes;
  final List<Comment> comments;
  final int? fee; // 用 int? 表示可选的整数

  // final String channel;
  // final List<String> inviteIDs;
  // final List<String> voters;

  Poll({
    required this.id,
    required this.type,
    required this.title,
    required this.description,
    required this.options,
    required this.options_rate,
    required this.time,
    required this.ipaddress,
    required this.sponsor,
    required this.votes,
    required this.comments,
    this.fee,
    // required this.channel,
    // required this.inviteIDs,
    // required this.voters,
  });

  // 从JSON数据创建Vote对象的工厂方法
  factory Poll.fromJson(Map<String, dynamic> json) {
    return Poll(
      id: json['id'],
      type: json['type'],
      title: json['title'],
      description: json['description'],
      options: List<String>.from(json['options']),
      options_rate: List<double>.from(json['options_rate']),
      time: json['time'],
      ipaddress: json['ipaddress'],
      sponsor: json['sponsor'],
      votes: List<Vote>.from(json['votes'].map((x) => Vote.fromJson(x))),
      comments: List<Comment>.from(json['comments'].map((x) => Comment.fromJson(x))),
      fee: json['fee'],
      // channel: json['channel'],
      // inviteIDs: List<String>.from(json['invite_ids']),
      // voters: List<String>.from(json['voters']),
    );
  }

  // 将Vote对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'type': type,
      'title': title,
      'description': description,
      'options': options,
      'options_rate': options_rate,
      'time': time,
      'ipaddress': ipaddress,
      'sponsor': sponsor,
      'votes': votes.map((x) => x.toJson()).toList(),
      'comments': comments.map((x) => x.toJson()).toList(),
      'fee': fee,
      // 'channel': channel,
      // 'invite_ids': inviteIDs,
      // 'voters': voters,
    };
  }
  
}

