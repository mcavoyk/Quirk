import 'package:flutter/material.dart';
import 'package:flutter_statusbarcolor/flutter_statusbarcolor.dart';
import 'pages/posts.dart';

const MaterialColor appPrimaryColor = Colors.amber;

void main() => runApp(App());

class App extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    FlutterStatusbarcolor.setStatusBarColor(appPrimaryColor);
    return MaterialApp(
      title: 'Quirk',
      color: appPrimaryColor,
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        primaryColor: appPrimaryColor,
        accentColor: Colors.white,
        fontFamily: 'Lato',
        accentIconTheme: IconThemeData(
          color: Colors.white
        )
     ),
      home: PostPage(),
    );
  }
}
