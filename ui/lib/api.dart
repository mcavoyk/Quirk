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

  Post({this.id, this.user, this.title, this.score, this.voteState, this.created, this.numComments});


  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      id: json['ID'],
      user: json['User'],
      title: json['Content'],
      score: json['Score'],
      voteState: json['VoteState'],
      created: DateTime.parse(json['CreatedAt']),
      numComments: json['NumComments'],
    );
  }
}
Future<String> auth() async{
  SharedPreferences pref = await SharedPreferences.getInstance();
  String token = pref.getString(savedToken);
  if (token != null){
    return token;
  }
  final response = await http.get('$api/auth/token');
  if (response.statusCode == 200) {
    String token = jsonDecode(response.body)['token'];
    pref.setString('auth_token', token);
    return token;
  } else {
    throw Exception('Network error');
  }
}

Future<Map<String, double>> checkOrGetLocationPerms() async{
  final perms = await SimplePermissions.checkPermission(Permission.AccessFineLocation);
  if (perms == true){
      return Location().getLocation();
  } else {
    await SimplePermissions.requestPermission(Permission.AccessFineLocation);
    return checkOrGetLocationPerms();
  }
}

Future<List<Post>> getPosts() async{
  final location = await checkOrGetLocationPerms();
  if (location == null) {
    throw Exception('Location unavailable');
  }
  final double latitude = location['latitude'];
  final double longitude = location['longitude'];
  final String token = await auth();
  final response = await http.get('$api/posts?lat=$latitude&lon=$longitude',
    headers: {HttpHeaders.authorizationHeader: "Bearer $token"},
  );
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

Future<Null> vote(String postID, int voteAction) async {
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

Future<Null> createPost(String postContent) async {
  final location = await checkOrGetLocationPerms();
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

  print("Json Encoding: $body");
  final response = await http.post('$api/post',
    headers: {HttpHeaders.authorizationHeader: "Bearer $token", HttpHeaders.contentTypeHeader: "application/json"},
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
