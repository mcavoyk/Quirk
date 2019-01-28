import 'package:flutter/material.dart';

class CommentPage extends StatefulWidget {
  @override
  _CommentPage createState() => new _CommentPage();
}

class _CommentPage extends State<CommentPage> {

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        titleSpacing: 0.0,
        title: Text("Comments"),
      ),
      body: RefreshIndicator (
        backgroundColor: Theme.of(context).primaryColor,
        onRefresh: () {
          print("Comment page");
          return null;
        },
        child: Text("Hi"),
      )
    );
  }
}