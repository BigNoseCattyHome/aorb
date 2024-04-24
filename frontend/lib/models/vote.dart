class Vote {
  final String id;
  final String userId;
  final String title;
  final String description;
  final List<String> options;
  Map<String, int> results; // key是选项，value是该选项的投票数
  final bool isAnonymous;
  final DateTime startTime;
  final DateTime endTime;

  Vote({
    required this.id,
    required this.userId,
    required this.title,
    required this.description,
    required this.options,
    required this.results,
    required this.isAnonymous,
    required this.startTime,
    required this.endTime,
  });

  // 从JSON数据创建Vote对象的工厂方法
  factory Vote.fromJson(Map<String, dynamic> json) {
    return Vote(
      id: json['id'],
      userId: json['userId'],
      title: json['title'],
      description: json['description'],
      options: List<String>.from(json['options']),
      results: Map<String, int>.from(json['results']),
      isAnonymous: json['isAnonymous'],
      startTime: DateTime.parse(json['startTime']),
      endTime: DateTime.parse(json['endTime']),
    );
  }

  // 将Vote对象转换成JSON的方法
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'userId': userId,
      'title': title,
      'description': description,
      'options': options,
      'results': results,
      'isAnonymous': isAnonymous,
      'startTime': startTime.toIso8601String(),
      'endTime': endTime.toIso8601String(),
    };
  }

  // 检查投票是否在有效期内
  bool isVoteActive() {
    final now = DateTime.now();
    return now.isAfter(startTime) && now.isBefore(endTime);
  }
}
