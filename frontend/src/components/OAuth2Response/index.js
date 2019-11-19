import React from 'react';
import { connect } from 'react-redux';
import { parse } from 'qs';
import { Redirect } from 'react-router-dom';
import { isMobile } from 'react-device-detect';
import { notification } from 'antd';
import { Toast as MobileToast } from 'antd-mobile';
import { withCookies } from 'react-cookie';

import { gotUserOAuth } from '../../actions/canvas';
import env from '../../util/env';

function OAuth2Response(props) {
  const query = parse(
    props.location.search.slice(1, props.location.search.length)
  );

  function processCanvasResponse(token, name, refreshToken) {
    props.dispatch(gotUserOAuth(token, refreshToken, query.subdomain));

    // set the version to current since it's a new user
    localStorage.prevVersion = env.currentVersion;

    const notificationMessage = `Welcome, ${
      name.split(' ')[0]
    }! You've successfully logged in with Canvas.`;

    if (isMobile) {
      MobileToast.success(notificationMessage);
    } else {
      notification.success({
        message: 'Success!',
        description: notificationMessage
      });
    }

    return <Redirect to="/dashboard" />;
  }

  switch (query.type) {
    case 'canvas':
      if (query.error || !query.canvas_response) {
        if (query.error === 'access_denied') {
          // user said no. redirect to home.
          return <Redirect to="/" />;
        }

        // unknown error
        notification.error({
          message: 'Unknown Error',
          duration: 0,
          description:
            'There was an unknown error logging you in with Canvas. Try again later.'
        });
        return <Redirect to="/" />;
      }

      const canvasResponse = JSON.parse(query.canvas_response);
      return processCanvasResponse(
        canvasResponse.access_token,
        canvasResponse.name,
        canvasResponse.refresh_token
      );
    case 'google':
      if (query.error) {
        if (query.error_source === 'proxy') {
          notification.error({
            message: 'Error from CanvasCBL',
            description: `Error from CanvasCBL: ${query.error_text}`
          });
        } else if (query.error_source === 'google') {
          notification.error({
            message: 'Error from Google',
            description: `There was an error from Google. ${query.body}`
          });
        }
        return <Redirect to={'/'} />;
      }

      this.props.cookies.set(name, query.session_string, {
        path: '/',
        secure: true,
        sameSite: false,
        // 13 days
        expires: Date.now() + 112300
      });

      switch (query.has_token) {
        case 'true':
          return <Redirect to={'/dashboard'} />;
        case 'false':
          return <Redirect to={'/tokenEntry'} />;
        default:
          notification.error({
            message: 'Missing has_token',
            description:
              'An error occurred that would not occur during normal use.'
          });
      }
      break;
    default:
      break;
  }

  notification.error({
    message: 'Unexpected Error',
    description: 'An unexpected error occurred during the Sign in flow.'
  });

  return <Redirect to={'/'} />;
}

const ConnectedOAuth2Response = connect(state => ({}))(OAuth2Response);

export default withCookies(ConnectedOAuth2Response);
