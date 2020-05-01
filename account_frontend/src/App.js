import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import {
  HashRouter as Router,
  Switch,
  Route,
  Redirect
} from 'react-router-dom';
import { isMobile } from 'react-device-detect';
import Home from './components/Home';
import { getUserProfile } from './actions/canvas';
import destToUrl from './util/destToUrl';
import styled, { keyframes } from 'styled-components';
import { Typography } from 'antd';
import OAuth2Response from './components/OAuth2Response';
import MainCard from './components/MainCard';
import Logout from './components/Logout';

const blurAnimation = keyframes`
  0% {
    filter: blur(0px);
    transform: scale(1);
  }
  
  100% {
    filter: blur(15px);
    transform: scale(1.1);
  }
`;

const HomeWrapper = styled.div`
  width: 100%;
  height: 100vh;
  overflow: hidden;
  position: absolute;
  top: 0px;
`;

const HomeBackground = styled.div`
  background-image: url('/home-background.jpeg');
  background-position: center;
  background-size: cover;
  height: 100%;
  animation: ${blurAnimation} 4s forwards;
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
`;

const routes = (
  <Router>
    <Switch>
      <Route exact path="/" component={Home} />
      <Route exact path="/oauth2response" component={OAuth2Response} />
      <Route exact path="/logout" component={Logout} />
      <Route path="/" render={() => <Redirect to="/" />} />
    </Switch>
  </Router>
);

function App(props) {
  const {
    redirectOk,
    loadingUserProfile,
    getUserProfileError,
    profile,
    destination,
    dispatch
  } = props;

  useEffect(() => {
    if (!loadingUserProfile && !getUserProfileError && !profile) {
      dispatch(getUserProfile());
    }
  });

  useEffect(() => {
    if (redirectOk === true && destination) {
      window.location = destToUrl(destination);
    }
  }, [redirectOk, destination]);

  let renderErr;

  if (getUserProfileError) {
    const err = getUserProfileError.error;
    if (
      err &&
      (err.includes('no session string') ||
        err.includes('invalid session string') ||
        err.includes('expired session'))
    ) {
      // this is OK
    } else {
      renderErr = (
        <MainCard>
          <Typography.Text type={'danger'}>
            We're in unknown lands, captain. (We've encountered an unexpected
            error.) Please try again later or contact us.
          </Typography.Text>
        </MainCard>
      );
    }
  }

  return (
    <HomeWrapper>
      {!isMobile && <HomeBackground />}
      {renderErr ? renderErr : routes}
    </HomeWrapper>
  );
}

const ConnectedApp = connect((state) => ({
  loadingUserProfile: state.canvas.loadingUserProfile,
  getUserProfileError: state.canvas.getUserProfileError,
  profile: state.canvas.profile,
  destination: state.home.destination,
  loadingLogout: state.home.loadingLogout,
  logoutError: state.home.logoutError,
  redirectOk: state.home.redirectOk
}))(App);

export default ConnectedApp;
