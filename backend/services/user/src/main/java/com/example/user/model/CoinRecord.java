package com.example.user.model;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

@Document(collection = "coinRecords")
public class CoinRecord {
    @Id
    private String id; // ObjectId generated by MongoDB
    private String userId;
    private int consume;
    private String questionId;

    // Constructors, getters, and setters omitted for brevity
}
