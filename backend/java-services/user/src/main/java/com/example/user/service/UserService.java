package com.example.user.service;

import com.example.user.exception.UserNotFoundException;
import com.example.user.model.UserId;
import com.example.user.repository.UserIdRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import com.example.user.exception.UserNotFoundException;

import java.util.Optional;
import java.util.Set;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
public class UserService {
    
    @Autowired
    private UserIdRepository userIdRepository;

    public Optional<UserId> findById(String id) {
        return userIdRepository.findById(id);
    }
    //判断用户是否关注另一个
    public boolean isUserFollowed(String user_id, String someone_user_id) {
        System.out.println("Finding user by ID: " + user_id);
        Optional<UserId> userOptional = userIdRepository.findById(user_id);
        if (userOptional.isPresent()) {
            UserId user = userOptional.get();
            System.out.println("User found: " + user);
            boolean isFollowed = user.getFollowed().contains(someone_user_id);
            System.out.println("User " + user_id + " follows " + someone_user_id + ": " + isFollowed);
            return isFollowed;
        } else {
            System.out.println("User with ID " + user_id + " not found");
        }
        return false;
    }
    //返回用户关注列表用户的昵称头像等
    
    public List<Map<String, String>> getFollowedUsersInfo(String userId, List<String> fields) {
        List<Map<String, String>> followedUsersInfo = new ArrayList<>();
        Optional<UserId> userOpt = userIdRepository.findById(userId);
        if (userOpt.isPresent()) {
            UserId user = userOpt.get();
            List<String> followedIds = user.getFollowed();
            for (String followedId : followedIds) {
                Optional<UserId> followedUserOpt = userIdRepository.findById(followedId);
                if (followedUserOpt.isPresent()) {
                    UserId followedUser = followedUserOpt.get();
                    Map<String, String> userInfo = new HashMap<>();
                    userInfo.put("nickname", followedUser.getNickname());
                    userInfo.put("avatar", followedUser.getAvatar());
                    userInfo.put("ipaddress", followedUser.getIpaddress());
                    followedUsersInfo.add(userInfo);
                } else {
                    System.out.println("Followed user not found for ID: " + followedId);
                }
            }
        }
        return followedUsersInfo;
    }
    //返回关注用户的粉丝的昵称头像
    public List<Map<String, String>> getUserFansInfo(String userId) {
        Optional<UserId> userOpt = userIdRepository.findById(userId);
        if (userOpt.isPresent()) {
            UserId user = userOpt.get();
            List<String> fansIds = user.getFollower(); // 获取关注这个用户的人的ID列表

            List<Map<String, String>> fansInfo = new ArrayList<>();
            for (String fanId : fansIds) {
                Optional<UserId> fanOpt = userIdRepository.findById(fanId);
                if (fanOpt.isPresent()) {
                    UserId fan = fanOpt.get();
                    Map<String, String> fanInfo = new HashMap<>();
                    fanInfo.put("nickname", fan.getNickname());
                    fanInfo.put("avatar", fan.getAvatar());
                    fanInfo.put("ipaddress", fan.getIpaddress());
                    fansInfo.add(fanInfo);
                }
            }
            return fansInfo;
        } else {
            throw new UserNotFoundException("User not found with ID: " + userId);
        }
    }
    //关注取关用户
        public void toggleFollow(String user_id, String follow_id) {
        Optional<UserId> userOpt = userIdRepository.findById(user_id);
        if (userOpt.isPresent()) {
            UserId user = userOpt.get();
            List<String> followed = user.getFollowed();
            if (followed.contains(follow_id)) {
                followed.remove(follow_id);
            } else {
                followed.add(follow_id);
            }
            userIdRepository.save(user);
        } else {
            throw new UserNotFoundException("User with ID " + user_id + " not found");
        }
    }
    //拉黑用户
    public void addToBlacklist(String userId, String blackId) {
        Optional<UserId> userOpt = userIdRepository.findById(userId);
        if (userOpt.isPresent()) {
            UserId user = userOpt.get();
            if (!user.getBlacklist().contains(blackId)) {
                user.getBlacklist().add(blackId);
                userIdRepository.save(user);
            }
        } else {
            throw new UserNotFoundException("User not found with ID: " + userId);
        }
    }
    //改头像昵称
    public void updateUser(String userId, String newNickname, String newAvatar) {
        Optional<UserId> userOpt = userIdRepository.findById(userId);
        if (userOpt.isPresent()) {
            UserId user = userOpt.get();
            user.setNickname(newNickname);
            user.setAvatar(newAvatar);
            userIdRepository.save(user);
        } else {
            throw new UserNotFoundException("User not found with ID: " + userId);
        }
    }
    //删除用户
    public void deleteUser(String userId) {
        Optional<UserId> userOpt = userIdRepository.findById(userId);
        if (userOpt.isPresent()) {
            userIdRepository.deleteById(userId);
        } else {
            throw new UserNotFoundException("User not found with ID: " + userId);
        }
    }
}
