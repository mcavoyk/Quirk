import React from 'react';
import { createAppContainer, createSwitchNavigator, createStackNavigator } from 'react-navigation';

import MainTabNavigator from './MainTabNavigator';
import SignInScreen from '../screens/SignInScreen';
import AuthLoadingScreen from '../screens/AuthLoadingScreen';
const AuthStack = createStackNavigator({ Sign: SignInScreen})

export default createAppContainer(createSwitchNavigator({
  Main: MainTabNavigator,
  Auth: AuthStack,
  AuthLoading: AuthLoadingScreen,
},
{
  initialRouteName: 'AuthLoading',
}));