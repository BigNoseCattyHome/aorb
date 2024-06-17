package com.example.user.repository;

import com.example.user.model.UserId;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface UserIdRepository extends MongoRepository<UserId, String> {
}
