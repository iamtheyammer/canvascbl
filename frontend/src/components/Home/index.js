import React, { useState } from 'react';
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
import getUrlPrefix from '../../util/getUrlPrefix';
import PopoutLink from '../PopoutLink';
import env from '../../util/env';

function Home(props) {
  const [enableSignin, setEnableSignin] = useState(false);

  if (props.token) {
    return <Redirect to="/dashboard" />;
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
          onChange={e => setEnableSignin(e.target.checked)}
        >
          I accept the{' '}
          <PopoutLink url={env.privacyPolicyUrl}>privacy policy</PopoutLink> and{' '}
          <PopoutLink url={env.termsOfServiceUrl}>terms of service</PopoutLink>
        </MobileCheckbox.AgreeItem>
        <MobileButton
          type="primary"
          disabled={!enableSignin}
          onClick={() =>
            (window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request`)
          }
        >
          Sign in with Canvas
        </MobileButton>
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
            onChange={e => setEnableSignin(e.target.checked)}
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
          <Button
            type="primary"
            size={!enableSignin ? 'default' : 'large'}
            className="center button"
            disabled={!enableSignin}
            onClick={() =>
              (window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request`)
            }
          >
            Sign in with Canvas
          </Button>
        </div>
      </Card>
    </div>
  );
}

const ConnectedHome = connect(state => ({
  token: state.canvas.token
}))(Home);

export default ConnectedHome;
