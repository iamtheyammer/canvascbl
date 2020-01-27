import React from 'react';
import { Redirect } from 'react-router-dom';
import { connect } from 'react-redux';
import './index.css';
import banner from '../../assets/banner.svg';

import { Card, Typography, Button, Checkbox } from 'antd';
import {
  Checkbox as MobileCheckbox,
  Button as MobileButton
} from 'antd-mobile';
import { isMobile } from 'react-device-detect';
import v4 from 'uuid/v4';
import getUrlPrefix from '../../util/getUrlPrefix';
import PopoutLink from '../PopoutLink';
import env from '../../util/env';
import {
  setGetSessionId,
  setSigninButtonAvailability
} from '../../actions/components/home';
import { getSessionInformation } from '../../actions/plus';

function Home(props) {
  const {
    signInButtonAvailability,
    getSessionId,
    loading,
    error,
    session,
    dispatch
  } = props;

  if (!session && !getSessionId) {
    const id = v4();
    dispatch(getSessionInformation(id));
    dispatch(setGetSessionId(id));
  }

  if (session) {
    return <Redirect to="/dashboard" />;
  }

  const getSessionErr = error[getSessionId];
  let getSessionErrText = '';
  if (getSessionErr) {
    if (!getSessionErr.res || !getSessionErr.res.data) {
      getSessionErrText = 'There was an error checking your sign-in status.';
    } else if (getSessionErr.res.data.error === 'expired session') {
      window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request?intent=reauth`;
      return null;
    } else if (getSession.res.data.error.includes('no session string')) {
    } else {
      getSessionErrText =
        'There was a server error. Contact support or try again later.';
    }
  }

  if (isMobile) {
    return (
      <div
        style={{
          textAlign: 'center',
          margin: '10px 5px 10px 5px',
          backgroundColor: '#ffffff'
        }}
      >
        <img src={banner} alt={'banner'} />
        <Typography.Title level={2}>Welcome!</Typography.Title>
        <Typography.Text>
          This tool calculates grades based on outcomes in Canvas. To use it,
          please accept the Terms and Conditions, then log in with Canvas.
        </Typography.Text>
        <MobileCheckbox.AgreeItem
          onChange={e =>
            dispatch(setSigninButtonAvailability(e.target.checked))
          }
        >
          I accept the{' '}
          <PopoutLink url={env.privacyPolicyUrl}>privacy policy</PopoutLink> and{' '}
          <PopoutLink url={env.termsOfServiceUrl}>terms of service</PopoutLink>
        </MobileCheckbox.AgreeItem>
        {getSessionErrText.length ? (
          <Typography.Text type="danger">{getSessionErrText}</Typography.Text>
        ) : (
          <MobileButton
            type="primary"
            disabled={!signInButtonAvailability}
            loading={loading.includes(getSessionId)}
            onClick={() =>
              (window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request`)
            }
          >
            Sign in with Canvas
          </MobileButton>
        )}
      </div>
    );
  }

  return (
    <div className="background">
      <Card className="card" title={<img src={banner} alt="banner" />}>
        <div className="static-text">
          <Typography.Title level={2}>Welcome!</Typography.Title>
          {env.buildBranch !== 'master' && (
            <Typography.Text type="danger">
              CanvasCBL is running in {env.buildBranch}
              <br />
            </Typography.Text>
          )}
          <Typography.Text>
            This tool calculates grades based on outcomes in Canvas. To use it,
            please accept the Terms and Conditions, then log in with Canvas.
          </Typography.Text>
        </div>
        <div>
          <Checkbox
            onChange={e =>
              dispatch(setSigninButtonAvailability(e.target.checked))
            }
            className="center-checkbox"
          >
            I accept the{' '}
            <PopoutLink url={env.privacyPolicyUrl}>privacy policy</PopoutLink>{' '}
            and{' '}
            <PopoutLink url={env.termsOfServiceUrl}>
              terms of service
            </PopoutLink>
          </Checkbox>
          <div style={{ marginTop: '15px' }} />
          {getSessionErrText.length ? (
            <Typography.Text type="danger">{getSessionErrText}</Typography.Text>
          ) : (
            <Button
              type="primary"
              size={!signInButtonAvailability ? 'default' : 'large'}
              loading={loading.includes(getSessionId)}
              className="center button"
              disabled={!signInButtonAvailability}
              onClick={() =>
                (window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request?intent=auth`)
              }
            >
              Sign in with Canvas
            </Button>
          )}
        </div>
      </Card>
    </div>
  );
}

const ConnectedHome = connect(state => ({
  getSessionId: state.components.home.getSessionId,
  signInButtonAvailability: state.components.home.signInButtonAvailability,
  loading: state.loading,
  error: state.error,
  session: state.plus.session
}))(Home);

export default ConnectedHome;
