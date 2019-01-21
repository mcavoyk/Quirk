import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

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
    _newRefresh();
  }

  Future<void> _newRefresh() {
    return getPosts().then((_posts) {
      setState(() {
          posts = _posts;
          message = "";
          loading = false;   
      });
    }).catchError((e) {
      setState(() {
        posts = new List();
        message = e.toString();
        print(message);
        loading = false;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    if (message == "" && posts.length == 0) {
      message = "No posts available";
    }
    
    return Scaffold(
      appBar: AppBar(
        titleSpacing: 0.0,
        title: Text("Quirk"),
      ),
      body: RefreshIndicator(
        onRefresh: _newRefresh,
        child: ListView.separated(
            physics: const AlwaysScrollableScrollPhysics(),
            itemCount: message != "" ? 1 : posts.length,
            padding: EdgeInsets.only(left:6.0, top: 10.0, bottom: 64.0),
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
      ),
      floatingActionButton: FloatingActionButton(
        backgroundColor: Theme.of(context).primaryColor,
        child: Icon(Icons.create),
        onPressed: () {
            showModalBottomSheet<void>(context: context, builder: (BuildContext context) {
              return Column(
                children: <Widget>[
                  Padding(
                    padding: EdgeInsets.only(left: 10, right: 10, top: 4),
                    child: TextFormField(
                      autofocus: true,
                      style: TextStyle(fontSize: 18),
                      decoration: InputDecoration(
                        icon: Icon(FontAwesomeIcons.comment, size: 32),
                        border: UnderlineInputBorder(),
                        hintText: "What's happening?"
                      ),
                    )
                  )
                ]
              );
            });
          }
      ),
    drawer: Drawer(
        child: ListView(       
          children: <Widget>[
          ],
        )
      ),
    );
  }
}
