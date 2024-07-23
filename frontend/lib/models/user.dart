// To parse this JSON data, do
//
//     final user = userFromJson(jsonString);

import 'dart:convert';

User userFromJson(String str) => User.fromJson(json.decode(str));

String userToJson(User data) => json.encode(data.toJson());


///user
class User {
    
    ///用户头像
    final String avatar;
    
    ///屏蔽好友
    final List<String> blacklist;
    
    ///用户的金币数
    final double coins;
    
    ///用户金币流水记录
    final List<CoinRecord>? coinsRecord;
    
    ///关注者
    final List<String> followed;
    
    ///被关注者
    final List<String> follower;
    
    ///用户ID
    final String id;
    
    ///IP归属地
    final String ipaddress;
    
    ///用户昵称
    final String nickname;
    
    ///发起过的问题
    final List<String> questionsAsk;
    
    ///回答过的问题
    final List<String> questionsAsw;
    
    ///收藏的问题
    final List<String> questionsCollect;

    User({
        required this.avatar,
        required this.blacklist,
        required this.coins,
        required this.coinsRecord,
        required this.followed,
        required this.follower,
        required this.id,
        required this.ipaddress,
        required this.nickname,
        required this.questionsAsk,
        required this.questionsAsw,
        required this.questionsCollect,
    });

    factory User.fromJson(Map<String, dynamic> json) => User(
        avatar: json["avatar"],
        blacklist: List<String>.from(json["blacklist"].map((x) => x)),
        coins: json["coins"]?.toDouble(),
        coinsRecord: json["coins_record"] == null ? [] : List<CoinRecord>.from(json["coins_record"]!.map((x) => CoinRecord.fromJson(x))),
        followed: List<String>.from(json["followed"].map((x) => x)),
        follower: List<String>.from(json["follower"].map((x) => x)),
        id: json["id"],
        ipaddress: json["ipaddress"],
        nickname: json["nickname"],
        questionsAsk: List<String>.from(json["questions_ask"].map((x) => x)),
        questionsAsw: List<String>.from(json["questions_asw"].map((x) => x)),
        questionsCollect: List<String>.from(json["questions_collect"].map((x) => x)),
    );

    Map<String, dynamic> toJson() => {
        "avatar": avatar,
        "blacklist": List<dynamic>.from(blacklist.map((x) => x)),
        "coins": coins,
        "coins_record": coinsRecord == null ? [] : List<dynamic>.from(coinsRecord!.map((x) => x.toJson())),
        "followed": List<dynamic>.from(followed.map((x) => x)),
        "follower": List<dynamic>.from(follower.map((x) => x)),
        "id": id,
        "ipaddress": ipaddress,
        "nickname": nickname,
        "questions_ask": List<dynamic>.from(questionsAsk.map((x) => x)),
        "questions_asw": List<dynamic>.from(questionsAsw.map((x) => x)),
        "questions_collect": List<dynamic>.from(questionsCollect.map((x) => x)),
    };
}


///一条金币流水记录
///
///coin_record
class CoinRecord {
    
    ///消耗的金币数
    final int consume;
    
    ///为其投币的问题ID
    final String questionId;
    
    ///使用者的ID
    final String userId;

    CoinRecord({
        required this.consume,
        required this.questionId,
        required this.userId,
    });

    factory CoinRecord.fromJson(Map<String, dynamic> json) => CoinRecord(
        consume: json["consume"],
        questionId: json["question_id"],
        userId: json["user_id"],
    );

    Map<String, dynamic> toJson() => {
        "consume": consume,
        "question_id": questionId,
        "user_id": userId,
    };
}