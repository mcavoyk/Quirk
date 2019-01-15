import 'package:flutter/material.dart';
import '../api.dart';
import '../widgets/posts.dart';

class PostPage extends StatefulWidget {
  @override
  _PostPage createState() => new _PostPage();
}

class _PostPage extends State<PostPage> {
  List<Post> posts = new List();
  bool loading = true;
  String message = "";


  @override
  void initState() {
    super.initState();
    _fetchPosts();
  }

  _fetchPosts() {
    setState(() {
      getPosts().then((newPosts) {
        posts = newPosts;
        message = "";
      }).catchError((e) {
        posts = new List();
        message = e.toString();
        print(message);
        
      });
      loading = false;
    });
  }

  Future<void> _refreshPosts() {
    return new Future<void>(() {_fetchPosts();});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Quirk')),
      body: RefreshIndicator(
        onRefresh: _refreshPosts,
        child: ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            itemCount: message != "" ? 1 : posts.length,
            padding: EdgeInsets.only(left:8.0, right: 4.0, top: 4.0),
            separatorBuilder: (BuildContext context, int index) => Divider(color: Colors.black, indent: 0),
            itemBuilder: (BuildContext context, int index) {
              if (message != "") {
                return Container(
                  alignment: Alignment.center,
                  padding: EdgeInsets.only(top: 16),
                  child: Text(message, style: TextStyle(fontSize: 20))
                );
              }
              //Post post = Post(user: 'Keunic', title: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.', score: 1);
              return PostBar(post: posts[index]);
            },
          )
      )
    );
  }
}
