import 'package:flutter/material.dart';

class comments_mine extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Container(
          width: 383,
          height: 62,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                '共3条评论',
                style: TextStyle(
                  color: Colors.black.withOpacity(0.5),
                  fontSize: 12,
                  fontFamily: 'Inter',
                  fontWeight: FontWeight.w600,
                  height: 0,
                ),
              ),
              const SizedBox(height: 12),
              Container(
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  mainAxisAlignment: MainAxisAlignment.start,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Container(
                      width: 35,
                      height: 35,
                      decoration: ShapeDecoration(
                        image: DecorationImage(
                          image: NetworkImage("https://via.placeholder.com/35x35"),
                          fit: BoxFit.fill,
                        ),
                        shape: OvalBorder(side: BorderSide(width: 1)),
                      ),
                    ),
                    const SizedBox(width: 15),
                    Container(
                      width: 333,
                      height: 30,
                      child: Stack(
                        children: [
                          Positioned(
                            left: 0,
                            top: 0,
                            child: Container(
                              width: 333,
                              height: 30,
                              decoration: ShapeDecoration(
                                color: Color(0x0C010000),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(20),
                                ),
                              ),
                            ),
                          ),
                          Positioned(
                            left: 13,
                            top: 6,
                            child: Text(
                              '说点什么...',
                              style: TextStyle(
                                color: Colors.black.withOpacity(0.30000001192092896),
                                fontSize: 15,
                                fontFamily: 'Inter',
                                fontWeight: FontWeight.w600,
                                height: 0,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}