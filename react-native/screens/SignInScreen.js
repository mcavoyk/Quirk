import React from 'react';
import {
  AsyncStorage,
  View,
  StyleSheet,
  Button,
  Text,
} from 'react-native';

import Api from '../constants/Api';

export default class SignInScreen extends React.Component {
    static navigationOptions = {
      title: 'Permission',
    };
    constructor(props) {
        super(props)
        this.state = {
            signInMessage: "",
        }
    }
  
    render() {
      return (
        <View style={styles.container}>
            <Text style={styles.signInMessage}>{this.state.signInMessage}</Text>
            <View style={styles.button}>
                <Button title="Sign in" onPress={this._signIn} />
            </View>
        </View>
      );
    }

    _signIn = async () => {
        try {
            let response = await fetch(
                Api.base + '/auth/token'
            );
            if (response.status != 200) {
                this.setState({signInMessage: "Error signing in"})
                return
            }
            let responseJson = await response.json();
            await AsyncStorage.setItem('userToken', responseJson.token);
            this.props.navigation.navigate('Main');
        } catch (error) {
            this.setState({signInMessage: "Error signing in"})
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
