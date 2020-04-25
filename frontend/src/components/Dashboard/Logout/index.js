import React from 'react';
import { Redirect } from 'react-router-dom';
import env from '../../../util/env';

function Logout(props) {
  window.location = `${env.accountUrl}logout`;

  return <Redirect to="/" />;
}

// const ConnectedLogout = connect(state => ({}))(Logout);

export default Logout;
