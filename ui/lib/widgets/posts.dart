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

  void _voteAction(int vote) {
    setState(() {
      if (_vote == vote) {
        _vote = _vote - vote;
      } else {
        _vote = vote;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return new Column(
      children: <Widget>[
        new Container(
          height: 100,
          child: new Row(
            children: <Widget>[
              new Flexible(
                fit: FlexFit.loose,
                flex: 5,
                child: new Text(
                  widget.post.title,
                  overflow: TextOverflow.fade,
                  style: TextStyle(fontSize: 20, fontWeight: FontWeight.w500, fontFamily: 'Lato',)
                  ),
              ),
              new Flexible(
                fit: FlexFit.tight,
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
                      child: new Text(widget.post.score.toString(), style: TextStyle(fontFamily: 'Lato', fontSize: 18)),
                      ),
                    new Flexible(
                      child: new IconButton(
                        icon: Icon(Icons.keyboard_arrow_down, color: _vote == -1 ? Colors.amber : Colors.grey), 
                        onPressed: () => _voteAction(-1),
                        padding: EdgeInsets.only(top: 0),
                        iconSize: 46,
                      ),
                    ),
                  ]),
              )
            ],
          )
        ),
        new Divider(
          color: Colors.black,
          indent: 0,
        )
      ],
    );
  }
}
