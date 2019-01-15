import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'package:location/location.dart';
import 'package:simple_permissions/simple_permissions.dart';
import 'package:shared_preferences/shared_preferences.dart';


final String api = 'http://192.168.0.32:5005/api/v1';

Future<String> auth() async{
  SharedPreferences pref = await SharedPreferences.getInstance();
  String token = pref.getString('auth_token');
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
  } else {
    throw Exception('Network error');
  }
}

class Post {
  final String user;
  final String title;
  final int score;
  final DateTime created;
  final int numComments;

  Post({this.user, this.title, this.score, this.created, this.numComments});

  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      user: json['User'],
      title: json['Content'],
      score: json['Score'],
      created: DateTime.parse(json['CreatedAt']),
      numComments: json['NumComments'],
    );
  }
}
