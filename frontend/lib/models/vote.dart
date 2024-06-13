// 单条投票记录，包含在question中
class Vote {
    
    // 投票者的选择
    final String choice;
    
    // 投票者ID
    final String userId;

    // 构造函数
    Vote({
        required this.choice,
        required this.userId,
    });

    // 从JSON数据创建Vote对象的工厂方法
    factory Vote.fromJson(Map<String, dynamic> json) {
        return Vote(
            choice: json['choice'],
            userId: json['user_id'],
        );
    }

    // 将Vote对象转换成JSON的方法
    Map<String, dynamic> toJson() {
        return {
            'choice': choice,
            'user_id': userId,
        };
    }
}