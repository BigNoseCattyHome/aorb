package com.example.user.repository;

import com.example.user.model.UserId;

import java.util.List;

import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.data.mongodb.repository.Query;
import org.springframework.stereotype.Repository;

@Repository
public interface UserIdRepository extends MongoRepository<UserId, String> {
    //@Query("{ 'followerId' : ?0, 'followeeId' : ?1 }")
    //boolean isfollowed(String followerId, String followeeId);
    //boolean existsByFollowerAnyInAndFollowedAnyIn(List<String> followerIds, List<String> followedIds);
}
