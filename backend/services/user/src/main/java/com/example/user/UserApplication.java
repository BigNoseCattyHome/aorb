package com.example.user;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class UserApplication {
    public static void main(String[] args) {
        System.out.println("Spring Boot application starting...");
        SpringApplication.run(UserApplication.class, args);
    }
}