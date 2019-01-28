import 'package:flutter/material.dart';
import '../api.dart';
import '../pages/comment.dart';

class PostBar extends StatelessWidget  { 
  PostBar(
    this.index,
    this.post,
    this.voteAction,
  );

  final int index;
  final Post post;
  final Function voteAction;
  
  @override
  Widget build(BuildContext context) {
    return RawMaterialButton(
      onPressed: () {
        print("Press $index");
        Navigator.push(context, MaterialPageRoute(builder: (context) => CommentPage()));
      },
      child: Column(
        children: <Widget>[
          new Container(
            height: 1,
            alignment: Alignment.centerLeft,
            child: new Text(
              post.user,
              style: TextStyle(fontSize: 13, color: Colors.black.withOpacity(0.6))
            )
          ),
          new Container(
            height: 100,
            child: new Row(
              children: <Widget>[
                new Flexible(
                  fit: FlexFit.tight,
                  flex: 5,
                  child: new Text(
                    post.title,
                    overflow: TextOverflow.fade,
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.w500)
                  ),
                ),
                new Flexible(
                  flex: 1,
                  child: new Column(
                    children: <Widget>[
                      new Flexible(
                        child: new IconButton(
                          icon: Icon(Icons.keyboard_arrow_up, color: post.voteState == 1 ? Colors.amber : Colors.grey),
                          onPressed: () => voteAction(index, 1),
                          padding: EdgeInsets.only(top: 0),
                          iconSize: 46,
                       ),
                      ),
                      new Container(
                        child: new Text(
                          post.score.toString(),
                          style: TextStyle(fontSize: 18),
                        )
                      ),
                      new Flexible(
                        child: new IconButton(
                          icon: Icon(Icons.keyboard_arrow_down, color: post.voteState == -1 ? Colors.amber : Colors.grey),
                          onPressed: () => voteAction(index, -1),
                          padding: EdgeInsets.only(top: 0),
                          iconSize: 46,
                        ),
                      ),
                    ]
                  ),
                )
              ],
            )
          ),
          Container(
            child: Row(
              children: <Widget>[
                Expanded(
                  flex: 2,
                  child: Text(
                    post.createdStr,
                    style: TextStyle(fontSize: 13, color: Colors.black.withOpacity(0.6))
                  )
                ),
                Flexible(
                  flex: 3,
                  child: Text(
                    post.numComments.toString() + ' Replies',
                    style: TextStyle(fontSize: 13, color: Colors.black.withOpacity(0.6))
                  )
                )
              ]
            )
          )
        ],
      )
    );
  }
}
