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
  if (perms == PermissionStatus.authorized){
      return Location().getLocation();
  } else if (perms == PermissionStatus.deniedNeverAsk) {
    return null;
  } else {
    await SimplePermissions.requestPermission(Permission.AccessFineLocation);
    return checkOrGetLocationPerms();
  }
}

Future<Post> getPosts() async{
  final location = await checkOrGetLocationPerms();
  if (location == null) {
    throw Exception('Network error');
  }
  final response = await http.get('$api/posts?lat=$location["latitude"]&lon=$location["longitude"]',
    headers: {HttpHeaders.authorizationHeader: "Bearer "},
  );
  if (response.statusCode == 200) {
    return Post.fromJson(json.decode(response.body));
  } else {
    throw Exception('Network error');
  }
}

class Post {
  final int userId;
  final String title;
  final int score;

  Post({this.userId, this.title, this.score});

  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      userId: json['User'],
      title: json['Content']['Title'],
      score: json['Score'],
    );
  }
}
