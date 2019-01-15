import 'package:flutter/material.dart';
import '../api.dart';

class PostBar extends StatefulWidget  { 
  PostBar({Key key, this.post}) : super(key: key);

  final Post post;

  @override
  _PostBar createState() => _PostBar();
}

class _PostBar extends State<PostBar> {
  int _vote = 0;

  void _voteAction(int newVote) {
    setState(() {
      if (_vote == newVote) {
        _vote = _vote - newVote;
      } else {
        _vote = newVote;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return new Column(
      children: <Widget>[
        new Container(
          alignment: Alignment.centerLeft,
          child: new Text(
            widget.post.user, 
            style: TextStyle(fontSize: 12, color: Colors.black.withOpacity(0.5))
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
                  widget.post.title,
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
                        icon: Icon(Icons.keyboard_arrow_up, color: _vote == 1 ? Colors.amber : Colors.grey), 
                        onPressed: () => _voteAction(1),
                        padding: EdgeInsets.only(top: 0),
                        iconSize: 46,
                     ),
                    ),
                    new Container(
                      child: new Text(
                        widget.post.score.toString(), 
                        style: TextStyle(fontSize: 18),
                      )
                    ),
                    new Flexible(
                      child: new IconButton(
                        icon: Icon(Icons.keyboard_arrow_down, color: _vote == -1 ? Colors.amber : Colors.grey), 
                        onPressed: () => _voteAction(-1),
                        padding: EdgeInsets.only(top: 0),
                        iconSize: 46,
                      ),
                    ),
                  ]
                ),
              )
            ],
          )
        )
      ],
    );
  }
}
