import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { logout } from '../../actions/canvas';
import MainCard from '../MainCard';
import { Typography } from 'antd';
import { Redirect } from 'react-router-dom';

function Logout(props) {
  const { loadingLogout, logoutError, loggedOut, dispatch } = props;

  useEffect(() => {
    if (!loadingLogout && !logoutError && !loggedOut) {
      dispatch(logout());
    }
  });

  return (
    <MainCard
      loading
      loadingText={
        <Typography.Title level={2}>
          {loadingLogout && 'Logging you out...'}
          {logoutError && 'Error logging you out'}
          {loggedOut && <Redirect to="/" />}
        </Typography.Title>
      }
    />
  );
}

const ConnectedLogout = connect((state) => ({
  redirectOk: state.home.redirectOk,
  loadingLogout: state.canvas.loadingLogout,
  logoutError: state.canvas.logoutError,
  loggedOut: state.canvas.loggedOut
}))(Logout);

export default ConnectedLogout;
