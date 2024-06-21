package com.example.user.controller;

import com.example.user.model.UserId;
import com.example.user.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import com.example.user.exception.UserNotFoundException;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/user")
public class UserController {

    @Autowired
    private UserService userService;

    @GetMapping("/{id}")
    public UserId getUserById(@PathVariable("id") String id, 
                              @RequestParam(value = "field", required = false) List<String> fields) {
        Optional<UserId> userOpt = userService.findById(id);
        if (userOpt.isPresent()) {
            UserId user = userOpt.get();
            if (fields != null && !fields.isEmpty()) {
                Set<String> fieldSet = fields.stream()
            .flatMap(f -> Arrays.stream(f.split(",")))
            .collect(Collectors.toSet());
                return filterFields(user, fieldSet);
            }
            return user;
        } else {
            throw new UserNotFoundException("User with ID not found");
        }
    }

    //判断用户是否关注一个人
    @GetMapping("/{user_id}/isFollowed")
    public ResponseEntity<Map<String, String>> isFollowed(
        @PathVariable("user_id") String userId,
        @RequestParam("someone_user_id") String someoneUserId
    ) {
        System.out.println("Received request to check if user " + userId + " follows " + someoneUserId);
        boolean isFollowed = userService.isUserFollowed(userId, someoneUserId);
        Map<String, String> response = new HashMap<>();
        response.put("isFollowed", isFollowed ? "YES" : "NO");
        System.out.println("Response: " + response);
        return ResponseEntity.ok(response);
    }
    //返回用户关注列表用户的头像昵称
    @GetMapping("/{user_id}/followed")
    public ResponseEntity<List<Map<String, String>>> getFollowedUsersInfo(
        @PathVariable("user_id") String userId
    ) {
        List<Map<String, String>> followedUsersInfo = userService.getFollowedUsersInfo(userId, null);
        return ResponseEntity.ok(followedUsersInfo);
    }
    //返回用户粉丝的头像昵称
    @GetMapping("/{user_id}/fans")
    public ResponseEntity<List<Map<String, String>>> getUserFans(
        @PathVariable("user_id") String userId
    ) {
        List<Map<String, String>> fansInfo = userService.getUserFansInfo(userId);
        return ResponseEntity.ok(fansInfo);
    }
    //关注/取关用户
    @PostMapping("/{user_id}/follower")
    public ResponseEntity<Void> toggleFollow(
        @PathVariable("user_id") String user_id,
        @RequestBody Map<String, String> request
    ) {
        String follow_id = request.get("follow_id");
        if (follow_id == null || follow_id.isEmpty()) {
            return ResponseEntity.badRequest().build();
        }
        userService.toggleFollow(user_id, follow_id);
        return ResponseEntity.ok().build();
    }
    //拉黑用户
    @PostMapping("/{user_id}/blacklist")
    public ResponseEntity<Void> addToBlacklist(
        @PathVariable("user_id") String userId,
        @RequestBody Map<String, String> requestBody
    ) {
        String blackId = requestBody.get("black_id");
        if (blackId == null || blackId.isEmpty()) {
            return ResponseEntity.badRequest().build();
        }

        userService.addToBlacklist(userId, blackId);
        return ResponseEntity.ok().build();
    }
    //改头像和昵称
    @PostMapping("/{user_id}")
    public ResponseEntity<Void> updateUser(
        @PathVariable("user_id") String userId,
        @RequestBody Map<String, String> requestBody
    ) {
        String newNickname = requestBody.get("nickname");
        String newAvatar = requestBody.get("avatar");
        if (newNickname == null || newNickname.isEmpty() || newAvatar == null || newAvatar.isEmpty()) {
            return ResponseEntity.badRequest().build();
        }

        userService.updateUser(userId, newNickname, newAvatar);
        return ResponseEntity.ok().build();
    }
    //删除用户
    @DeleteMapping("/{user_id}")
    public ResponseEntity<Void> deleteUser(@PathVariable("user_id") String userId) {
        userService.deleteUser(userId);
        return ResponseEntity.ok().build();
    }


    private UserId filterFields(UserId user, Set<String> fields) {
        UserId filteredUser = new UserId();
        //assert user.getId() != null : "ID cannot be null";
        filteredUser.setId(user.getId());
        //if (fields.contains("isFollowed")) filteredUser.setisFollowed("yes");
        //else filteredUser.setisFollowed("no");
        
        if (fields.contains("username")) filteredUser.setusername(user.getusername());
        else filteredUser.setusername(user.getusername());

        if (fields.contains("password")) filteredUser.setPassword(user.getPassword());
        else filteredUser.setPassword(user.getPassword());

        if (fields.contains("nickname")) filteredUser.setNickname(user.getNickname());
        else filteredUser.setNickname(user.getNickname());

        if (fields.contains("avatar")) filteredUser.setAvatar(user.getAvatar());
        else filteredUser.setAvatar(user.getAvatar());

        if (fields.contains("coins")) filteredUser.setCoins(user.getCoins());
        else filteredUser.setCoins(user.getCoins());

        if (fields.contains("coins_record")) filteredUser.setCoins_record(user.getCoins_record());
        else filteredUser.setCoins_record(user.getCoins_record());

        if (fields.contains("followed")) filteredUser.setFollowed(user.getFollowed());
        else filteredUser.setFollowed(user.getFollowed());

        if (fields.contains("follower")) filteredUser.setFollower(user.getFollower());
        else filteredUser.setFollower(user.getFollower());

        if (fields.contains("blacklist")) filteredUser.setBlacklist(user.getBlacklist());
        else filteredUser.setBlacklist(user.getBlacklist());

        if (fields.contains("questions_ask")) filteredUser.setQuestions_ask(user.getQuestions_ask());
        else filteredUser.setQuestions_ask(user.getQuestions_ask());

        if (fields.contains("questionsAsw")) filteredUser.setQuestions_asw(user.getQuestions_asw());
        else filteredUser.setQuestions_asw(user.getQuestions_asw());

        if (fields.contains("questionsCollect")) filteredUser.setQuestions_collect(user.getQuestions_collect());
        else filteredUser.setQuestions_collect(user.getQuestions_collect());

        if (fields.contains("ipaddress")) filteredUser.setIpaddress(user.getIpaddress());
        else filteredUser.setIpaddress(user.getIpaddress());

        return filteredUser;
    }
}
