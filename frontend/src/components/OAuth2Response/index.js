import React from 'react';
import { connect } from 'react-redux';
import { parse } from 'qs';
import { Redirect } from 'react-router-dom';
import { isMobile } from 'react-device-detect';
import { notification } from 'antd';
import { Toast as MobileToast } from 'antd-mobile';
import env from '../../util/env';

function OAuth2Response(props) {
  const query = parse(
    props.location.search.slice(1, props.location.search.length)
  );

  if (query.dest === 'teacher') {
    window.location = env.teacherUrl;
    return null;
  }

  function processCanvasResponse(name, intent) {
    // set the version to current since it's a new user
    localStorage.prevVersion = env.currentVersion;

    const firstName = name.split(' ')[0];
    const notificationMessage =
      intent === 'auth'
        ? `Welcome, ${firstName}! You've successfully logged in with Canvas.`
        : `Welcome back, ${firstName}!`;

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

  if (query.type === 'canvas') {
    if (query.error) {
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
    } else {
      return processCanvasResponse(query.name, query.intent);
    }
  }

  return <Redirect to={'/'} />;
}

const ConnectedOAuth2Response = connect((state) => ({}))(OAuth2Response);

export default ConnectedOAuth2Response;
