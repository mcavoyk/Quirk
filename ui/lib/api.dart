import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'package:location/location.dart';
import 'package:simple_permissions/simple_permissions.dart';
import 'package:shared_preferences/shared_preferences.dart';


final String api = 'http://quirk.afforess.com/api/v1';
final String savedToken = 'auth_token';


class Post {
  final String id;
  final String user;
  final String title;
  int score;
  int voteState;
  final DateTime created;
  final int numComments;
  final String createdStr;

  Post({this.id, this.user, this.title, this.score, this.voteState, this.created, this.createdStr, this.numComments});


  factory Post.fromJson(Map<String, dynamic> json) {
    var createdTime = DateTime.parse(json['CreatedAt']);
    return Post(
      id: json['ID'],
      user: json['User'],
      title: json['Content'],
      score: json['Score'],
      voteState: json['VoteState'],
      created: createdTime,
      createdStr: parseTime(createdTime),
      numComments: json['NumComments'],
    );
  }

  static String parseTime(DateTime time) {
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
}

class Api {

  static Map<String, double> location;
  static DateTime locationAt;

  static Future<String> auth() async {
    final location = await _getLocation();
    if (location == null) {
      throw Exception('Location unavailable');
    }
    final double latitude = location['latitude'];
    final double longitude = location['longitude'];

    SharedPreferences pref = await SharedPreferences.getInstance();
    String token = pref.getString(savedToken);
    if (token != null) {
      return token;
    }
    final response = await http.get('$api/auth/token?lat=$latitude&lon=$longitude');
    if (response.statusCode == 200) {
      String token = jsonDecode(response.body)['token'];
      pref.setString('auth_token', token);
      return token;
    } else {
      throw Exception('Network error');
    }
  }


  static bool _hasLocation() {
    if (locationAt == null || DateTime.now().difference(locationAt) > Duration(seconds: 60)) {
      return false;
    }
    return true;
  }

  static Future<Map<String, double>> _getLocation() async {
    if (_hasLocation()) {
      return location;
    }

    final perms = await SimplePermissions.checkPermission(
        Permission.AccessFineLocation);
    if (perms == true) {
      location = await Location().getLocation();
      locationAt = DateTime.now();
      return location;
    } else {
      await SimplePermissions.requestPermission(Permission.AccessFineLocation);
      return _getLocation();
    }
  }

  static Future<List<Post>> getPosts() async {
    final location = await _getLocation();
    if (location == null) {
      throw Exception('Location unavailable');
    }
    final double latitude = location['latitude'];
    final double longitude = location['longitude'];
    final String token = await auth();
    final response = await http.get('$api/posts?lat=$latitude&lon=$longitude',
      headers: {HttpHeaders.authorizationHeader: "Bearer $token"},
    ).timeout(Duration(seconds: 5));
    if (response.statusCode == 200) {
      List<dynamic> postsJson = json.decode(response.body)['Posts'];
      List<Post> posts = new List();
      postsJson.forEach((i) => posts.add(Post.fromJson(i)));
      return posts;
    } else if (response.statusCode == 403) {
      SharedPreferences pref = await SharedPreferences.getInstance();
      pref.remove(savedToken);
      return getPosts();
    } else {
      throw Exception('Network error');
    }
  }

  static Future<Null> vote(String postID, int voteAction) async {
    final String token = await auth();
    final response = await http.post('$api/post/$postID/vote?state=$voteAction',
      headers: {HttpHeaders.authorizationHeader: "Bearer $token"},
    );
    if (response.statusCode == 200) {
      return null;
    } else if (response.statusCode == 403) {
      SharedPreferences pref = await SharedPreferences.getInstance();
      pref.remove(savedToken);
      return vote(postID, voteAction);
    } else {
      throw Exception('Network error');
    }
  }

  static Future<Null> createPost(String postContent) async {
    final location = await _getLocation();
    if (location == null) {
      throw Exception('Location unavailable');
    }
    final double latitude = location['latitude'];
    final double longitude = location['longitude'];
    final String token = await auth();
    final String body = jsonEncode({
      "lat": latitude,
      "lon": longitude,
      "accessType": "public",
      "content": postContent
    });

    final response = await http.post('$api/post',
        headers: {
          HttpHeaders.authorizationHeader: "Bearer $token",
          HttpHeaders.contentTypeHeader: "application/json"
        },
        body: body
    );
    if (response.statusCode == 200) {
      return null;
    } else if (response.statusCode == 403) {
      SharedPreferences pref = await SharedPreferences.getInstance();
      pref.remove(savedToken);
      return createPost(postContent);
    } else {
      throw Exception('Network error');
    }
  }
}
