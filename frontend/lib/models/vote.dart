import 'package:flutter/material.dart';
import './comment.dart';

class Vote {
  final String id;
  final String type; // 这里我们使用 String 替代原 Type 枚举
  final String time;
  final String sponsor;
  final String title;
  final String description;
  final List<String> options;
  final String channel;
  final List<Comment> comments;
  final int? fee; // 用 int? 表示可选的整数
  final List<String> inviteIDs;
  final List<String> voters;

  Vote({
    required this.id,
    required this.channel,
    required this.comments,
    required this.description,
    this.fee,
    required this.sponsor,
    required this.time,
    required this.title,
    required this.options,
    required this.inviteIDs,
    required this.type,
    required this.voters,
  });

  // 从JSON数据创建Vote对象的工厂方法
  factory Vote.fromJson(Map<String, dynamic> json) {
    return Vote(
      id: json['id'],
      description: json['description'],
      channel: json['channel'],
      comments: List<Comment>.from(json['comments'].map((x) => Comment.fromJson(x))),
      fee: json['fee'],
      sponsor: json['sponsor'],
      time: json['time'],
      title: json['title'],
      options: List<String>.from(json['options']),
      inviteIDs: List<String>.from(json['invite_ids'] ?? []),
      type: json['type'],
      voters: List<String>.from(json['voters']),
    );
  }

  // 将Vote对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'channel': channel,
      'comments': List<dynamic>.from(comments.map((x) => x.toJson())),
      'fee': fee,
      'sponsor': sponsor,
      'time': time,
      'title': title,
      'description': description,
      'options': options,
      'invite_ids': inviteIDs,
      'type': type,
      'voters': voters,
    };
  }
}

