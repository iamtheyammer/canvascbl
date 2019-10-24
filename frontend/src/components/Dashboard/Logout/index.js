import React from 'react';
import { connect } from 'react-redux';
import { useCookies } from 'react-cookie';
import { Redirect } from 'react-router-dom';

import { logout } from '../../../actions/canvas';

function Logout(props) {
  const { token, subdomain } = props;
  const [, , removeCookie] = useCookies([]);

  removeCookie('session_string');
  props.dispatch(logout(token, subdomain));
  return <Redirect to="/" />;
}

const ConnectedLogout = connect(state => ({
  token: state.canvas.token,
  subdomain: state.canvas.subdomain
}))(Logout);

export default ConnectedLogout;
