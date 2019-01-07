import React from 'react';
import {
  Image,
  Platform,
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
  AsyncStorage,
} from 'react-native';
import {  Location, Permissions } from 'expo';

import Api from '../constants/Api';

export default class HomeScreen extends React.Component {
  static navigationOptions = {
    header: null,
  };

  constructor(props) {
    super(props)
    this.state = {
        loading: true,
        error: null,
        posts: [],
        location: null,
    }
  }

  async componentDidMount() {
    await this._getLocationAsync()
    this._fetchPostsAsync()
  }

  _fetchPostsAsync = async () => {
    try {
      console.log(this.state.location)
      let response = await fetch(
        Api.base + '/posts' + "?lat=" + this.state.location.coords.latitude + "&lon=" + this.state.location.coords.longitude, { 
          headers: {
            Authorization: 'Bearer ' + await AsyncStorage.getItem('userToken')
          }
        }
      );
      if (response.status !== 200) {
        this.setState({error: "Error performing network request"})
        return
      }
      let responseJson = await response.json();
      console.log(responseJson)
      this.setState({posts: responseJson, error: null, loading: false})
    } catch (error) {
      console.error(error);
      this.setState({error: "An error occurred"})
    }
  }

  _getLocationAsync = async () => {
    let { status } = await Permissions.askAsync(Permissions.LOCATION);
    if (status !== 'granted') {
      await AsyncStorage.clear();
      this.props.navigation.navigate('Auth');
      return
    }
    let location = await Location.getCurrentPositionAsync({});
    console.log(location)
    this.setState({ location });
  };

  render() {
    const {posts, loading, error} = this.state
    let text = ""
    if (error !== null) { text = error }
    else if (loading) { text = 'Loading...'}
    else { text = 'Found ' + posts.length + ' posts.'}
    return (
      <View style={styles.container}>
        <Text style={styles.text}>{text}</Text>
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    justifyContent: 'center',
  },
  text: {
    marginBottom: 20,
    color: 'rgba(0,0,0,0.4)',
    fontSize: 14,
    lineHeight: 19,
    textAlign: 'center',
  },
  contentContainer: {
    paddingTop: 30,
  },
  welcomeContainer: {
    alignItems: 'center',
    marginTop: 10,
    marginBottom: 20,
  },
  welcomeImage: {
    width: 100,
    height: 80,
    resizeMode: 'contain',
    marginTop: 3,
    marginLeft: -10,
  },
  getStartedContainer: {
    alignItems: 'center',
    marginHorizontal: 50,
  },
  homeScreenFilename: {
    marginVertical: 7,
  },
  codeHighlightText: {
    color: 'rgba(96,100,109, 0.8)',
  },
  codeHighlightContainer: {
    backgroundColor: 'rgba(0,0,0,0.05)',
    borderRadius: 3,
    paddingHorizontal: 4,
  },
  getStartedText: {
    fontSize: 17,
    color: 'rgba(96,100,109, 1)',
    lineHeight: 24,
    textAlign: 'center',
  },
  tabBarInfoContainer: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    ...Platform.select({
      ios: {
        shadowColor: 'black',
        shadowOffset: { height: -3 },
        shadowOpacity: 0.1,
        shadowRadius: 3,
      },
      android: {
        elevation: 20,
      },
    }),
    alignItems: 'center',
    backgroundColor: '#fbfbfb',
    paddingVertical: 20,
  },
  tabBarInfoText: {
    fontSize: 17,
    color: 'rgba(96,100,109, 1)',
    textAlign: 'center',
  },
  navigationFilename: {
    marginTop: 5,
  },
  helpContainer: {
    marginTop: 15,
    alignItems: 'center',
  },
  helpLink: {
    paddingVertical: 15,
  },
  helpLinkText: {
    fontSize: 14,
    color: '#2e78b7',
  },
});
