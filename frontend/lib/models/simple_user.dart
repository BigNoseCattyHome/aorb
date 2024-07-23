// To parse this JSON data, do
//
//     final simpleUser = simpleUserFromJson(jsonString);

import 'dart:convert';

SimpleUser simpleUserFromJson(String str) => SimpleUser.fromJson(json.decode(str));

String simpleUserToJson(SimpleUser data) => json.encode(data.toJson());


///simple_user
class SimpleUser {
    
    ///头像
    final String avatar;
    
    ///IP归属地
    final String ipaddress;
    
    ///昵称
    final String nickname;

    SimpleUser({
        required this.avatar,
        required this.ipaddress,
        required this.nickname,
    });

    factory SimpleUser.fromJson(Map<String, dynamic> json) => SimpleUser(
        avatar: json["avatar"],
        ipaddress: json["ipaddress"],
        nickname: json["nickname"],
    );

    Map<String, dynamic> toJson() => {
        "avatar": avatar,
        "ipaddress": ipaddress,
        "nickname": nickname,
    };
}