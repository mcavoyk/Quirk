import 'package:flutter/material.dart';
import '../api.dart';

class PostBar extends StatefulWidget  { 
  PostBar({Key key, this.post}) : super(key: key);

  final Post post;

  @override
  _PostBar createState() => _PostBar();
}

class _PostBar extends State<PostBar> {
  int voteAction = 0;
  int voteScore = 0;

  @override
  void initState() {
      super.initState();
      voteScore = widget.post.score;
    }

  void _voteAction(int newVote) {
    setState(() {
      if (voteAction == newVote) {
        voteAction = voteAction - newVote;
        voteScore = widget.post.score;
      } else {
        voteAction = newVote;
        voteScore += newVote;
      }
    });
  }

  String parseTime(DateTime time) {
    Duration diff = DateTime.now().difference(time);
    int days = diff.inDays;
    if (days != 0) {
      return days.toString() + 'd';
    }
    int hours = diff.inHours;
    if (hours != 0) {
      return hours.toString() + 'h';
    }
    int mins = diff.inMinutes;
    if (mins != 0) {
      return mins.toString() + 'm';
    }
    int secs = diff.inSeconds;
    if (secs != 0) {
      return secs.toString() + 's';
    }
    return '1s';
  }

  @override
  Widget build(BuildContext context) {
    return new Column(
      children: <Widget>[
        new Container(
          height: 1,
          alignment: Alignment.centerLeft,
          child: new Text(
            widget.post.user, 
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
                        icon: Icon(Icons.keyboard_arrow_up, color: voteAction == 1 ? Colors.amber : Colors.grey), 
                        onPressed: () => _voteAction(1),
                        padding: EdgeInsets.only(top: 0),
                        iconSize: 46,
                     ),
                    ),
                    new Container(
                      child: new Text(
                        voteScore.toString(), 
                        style: TextStyle(fontSize: 18),
                      )
                    ),
                    new Flexible(
                      child: new IconButton(
                        icon: Icon(Icons.keyboard_arrow_down, color: voteAction == -1 ? Colors.amber : Colors.grey), 
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
        ),
        Container(
          child: Row(
            children: <Widget>[
              Expanded(
                flex: 2,
                child: Text(
                  parseTime(widget.post.created),
                  style: TextStyle(fontSize: 13, color: Colors.black.withOpacity(0.6))
                )
              ),
              Flexible(
                flex: 3,
                child: Text(
                  widget.post.numComments.toString() + ' Replies',
                  style: TextStyle(fontSize: 13, color: Colors.black.withOpacity(0.6))
                )
              )
            ]
          )
        )
      ],
    );
  }
}