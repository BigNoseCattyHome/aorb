package com.example.user.service;

import com.example.user.model.UserId;
import com.example.user.repository.UserIdRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Service
public class UserService {

    @Autowired
    private UserIdRepository userIdRepository;

    public Optional<UserId> findById(String id) {
        return userIdRepository.findById(id);
    }
}
