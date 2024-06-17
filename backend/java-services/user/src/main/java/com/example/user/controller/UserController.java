package com.example.user.controller;

import com.example.user.model.UserId;
import com.example.user.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;

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
                Set<String> fieldSet = fields.stream().flatMap(f -> List.of(f.split(",")).stream()).collect(Collectors.toSet());
                return filterFields(user, fieldSet);
            }
            return user;
        } else {
            throw new UserNotFoundException();
        }
    }

    private UserId filterFields(UserId user, Set<String> fields) {
        UserId filteredUser = new UserId();
        if (fields.contains("id")) filteredUser.setId(user.getId());
        if (fields.contains("password")) filteredUser.setPassword(user.getPassword());
        if (fields.contains("nickname")) filteredUser.setNickname(user.getNickname());
        if (fields.contains("avatar")) filteredUser.setAvatar(user.getAvatar());
        if (fields.contains("coins")) filteredUser.setCoins(user.getCoins());
        if (fields.contains("coinsRecord")) filteredUser.setCoinsRecord(user.getCoinsRecord());
        if (fields.contains("followed")) filteredUser.setFollowed(user.getFollowed());
        if (fields.contains("follower")) filteredUser.setFollower(user.getFollower());
        if (fields.contains("blacklist")) filteredUser.setBlacklist(user.getBlacklist());
        if (fields.contains("questionsAsk")) filteredUser.setQuestionsAsk(user.getQuestionsAsk());
        if (fields.contains("questionsAsw")) filteredUser.setQuestionsAsw(user.getQuestionsAsw());
        if (fields.contains("questionsCollect")) filteredUser.setQuestionsCollect(user.getQuestionsCollect());
        if (fields.contains("ipaddress")) filteredUser.setIpaddress(user.getIpaddress());
        return filteredUser;
    }
}
