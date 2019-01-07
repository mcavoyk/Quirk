import React from 'react';
import {
  AsyncStorage,
  ScrollView,
  StyleSheet,
  View,
  Button,
} from 'react-native';

export default class SettingsScreen extends React.Component {
  static navigationOptions = {
    title: 'Settings',
  };

  _signOut = async () => {
    await AsyncStorage.clear();
    this.props.navigation.navigate('Auth');
  };

  render() {
    return (
      <View style={styles.container}>
          <ScrollView style={styles.button}>
              <Button  title="Sign out" onPress={this._signOut} />
          </ScrollView>
      </View>
    );
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
      paddingTop: 10,
  },
});