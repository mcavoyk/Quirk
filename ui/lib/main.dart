import 'package:flutter/material.dart';
import 'package:flutter_statusbarcolor/flutter_statusbarcolor.dart';
import 'pages/posts.dart';

const MaterialColor primaryColor = Colors.amber;
const Color accentColor = Colors.white;

void main() => runApp(App());

class App extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    FlutterStatusbarcolor.setStatusBarColor(primaryColor.shade600);
    return MaterialApp(
      title: 'Quirk',
      color: primaryColor,
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        primaryColor: primaryColor,
        accentColor: accentColor,
        fontFamily: 'Lato',
        iconTheme: IconThemeData(
          color: accentColor
        ),
        accentIconTheme: IconThemeData(
          color: primaryColor
        )
     ),
      home: PostPage(),
    );
  }
}
