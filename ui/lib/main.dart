import 'package:flutter/material.dart';
import 'pages/posts.dart';

void main() {
  runApp(new MaterialApp(
    home: PostPage(), 
    theme: ThemeData(
      primarySwatch: Colors.amber,
      fontFamily: 'Lato',
    ),
  ));
} 
