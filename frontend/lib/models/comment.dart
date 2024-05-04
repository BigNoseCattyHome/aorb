class Comment {
  final String advise;
  final String choose;
  final String userid;

  Comment({
    required this.advise,
    required this.choose,
    required this.userid,
  });

  factory Comment.fromJson(Map<String, dynamic> json) {
    return Comment(
      advise: json['advise'],
      choose: json['choose'],
      userid: json['userid'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'advise': advise,
      'choose': choose,
      'userid': userid,
    };
  }
}
