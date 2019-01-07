import React from 'react';
import {
  AsyncStorage,
  View,
  StyleSheet,
  Button,
  Text,
  Platform,
} from 'react-native';
import { Constants, Location, Permissions } from 'expo';
import Api from '../constants/Api';

export default class SignInScreen extends React.Component {
    static navigationOptions = {
      title: 'Permissions',
    };
    constructor(props) {
        super(props)
        this.state = {
            errorMessage: "",
            location: null,
        }
    }
  

    componentWillMount() {
        if (Platform.OS === 'android' && !Constants.isDevice) {
          this.setState({
            errorMessage: 'Oops, this will not work on Sketch in an Android emulator. Try it on your device!',
          });
        } else {
            this._getLocationAsync();
        }
      }
    
      _getLocationAsync = async () => {
        let { status } = await Permissions.askAsync(Permissions.LOCATION);
        if (status !== 'granted') {
            this.setState({
                errorMessage: 'Locations permissions is required to use Quirk.\nPlease enable this in your settings to continue.',
              });
          return
        }
        let location = await Location.getCurrentPositionAsync({});
        console.log(location)
        this.setState({ location });
      };


    _checkPermissions = async () => {
        let { status } = await Permissions.getAsync(Permissions.LOCATION)
        return status === 'granted'
    }

    render() {
      return (
        <View style={styles.container}>
            <Text style={styles.signInMessage}>{this.state.errorMessage}</Text>
            <View style={styles.button}>
                <Button title="Sign in" onPress={this._signIn} />
            </View>
        </View>
      );
    }

    _signIn = async () => {
        permissions = await this._checkPermissions()
        if (!permissions) {
            this.setState({
                errorMessage: 'Locations permissions is required to use Quirk.\nPlease enable this in your settings to continue.',
              });
            return
        }
        try {
            let response = await fetch(
                Api.base + '/auth/token'
            );
            if (response.status != 200) {
                this.setState({errorMessage: "Error signing in"})
                return
            }
            let responseJson = await response.json();
            await AsyncStorage.setItem('userToken', responseJson.token);
            this.props.navigation.navigate('Main');
        } catch (error) {
            this.setState({errorMessage: "Error signing in"})
            //console.error(error);
        }
    }

    
  }
  
const styles = StyleSheet.create({
    container: {
        flex: 1,
        alignItems: 'stretch',
        justifyContent: 'center',
    },
    button: {
        paddingLeft: 100,
        paddingRight: 100,
    },
    signInMessage: {
        fontSize: 17,
        color: 'rgba(96,100,109, 1)',
        textAlign: 'center',
        paddingBottom: 10,
    }
});
