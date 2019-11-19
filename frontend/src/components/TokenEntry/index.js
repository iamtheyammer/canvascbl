import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import v4 from 'uuid/v4';
import { withCookies } from 'react-cookie';
import { Input, Typography, Button, Card } from 'antd';

import env from '../../util/env';
import { sendCanvasToken } from '../../actions/canvas';
import PopoutLink from '../PopoutLink';

class TokenEntry extends Component {
  constructor(props) {
    super(props);
    this.state = {
      token: this.props.token || '',
      subdomain: this.props.subdomain || '',
      sendTokenId: '',
      error: ''
    };
  }

  tokenOnChange = e => {
    e.preventDefault();

    this.setState({
      token: e.target.value
    });
  };

  onSubmit = e => {
    e.preventDefault();

    if (this.state.token.length < 64) {
      this.setState({ error: "That doesn't look like a canvas token." });
      return;
    }

    const id = v4();
    this.props.dispatch(sendCanvasToken(id, this.state.token));

    this.setState({
      sendTokenId: id,
      error: ''
    });

    // notification.success({
    //   message: 'Success!',
    //   description: 'Saved token and subdomain.'
    // });
  };

  subdomainOnChange = e => {
    e.preventDefault();

    this.setState({
      subdomain: e.target.value
    });
  };

  componentDidUpdate(prevProps, prevState, snapshot) {
    const sendErr = this.props.error[this.state.sendTokenId];
    if (sendErr) {
      const textSendErr = sendErr.res.response.data;
      if (this.state.error !== textSendErr) {
        this.setState({
          error: textSendErr
        });
      }
    }
  }

  render() {
    if (!this.props.cookies.get('session_string')) {
      return <Redirect to={'/'} />;
    }

    if (this.props.successfullySentToken) {
      return <Redirect to={'/dashboard'} />;
    }

    return (
      <div>
        <Card title="Finish Signing In">
          <Typography.Text>
            We just need you to do one more thing to use CanvasCBL. It'll only
            take a minute and you only have to do it once.
          </Typography.Text>
          <Typography.Title level={3}>
            Step 1: Generate a token
          </Typography.Title>
          <Typography.Text>
            Click{' '}
            <PopoutLink
              url={`https://${env.defaultSubdomain}.instructure.com/profile/settings`}
            >
              here
            </PopoutLink>
            , then scroll down and click on "New Access Token" towards the
            bottom of the page.
            <br />
            Enter something like CanvasCBL as the purpose and{' '}
            <Typography.Text strong>leave the expiration blank</Typography.Text>
            . If you don't leave the expiration blank, some weird stuff might
            happen and you'll have to contact support.
          </Typography.Text>
          <Typography.Title level={3}>Step 2: Copy the token</Typography.Title>
          <Typography.Text>
            Copy the entire token-- that random-looking bold thing.
          </Typography.Text>
          <Typography.Title level={3}>Step 3: Paste it below</Typography.Title>
          <Typography.Text>
            Paste your token in the field below, then click on Submit.
          </Typography.Text>
          <Input
            addonBefore="Canvas Token"
            onPressEnter={this.onSubmit}
            onChange={this.tokenOnChange}
            value={this.state.token}
          />
          {this.state.error && (
            <div>
              <br />
              <Typography.Text type="danger">
                {this.state.error}
              </Typography.Text>
            </div>
          )}

          <br />
          {/*<Input*/}
          {/*  addonBefore="Canvas Subdomain"*/}
          {/*  onPressEnter={this.onSubmit}*/}
          {/*  onChange={this.subdomainOnChange}*/}
          {/*  value={this.state.subdomain}*/}
          {/*/>*/}
          <br />
          <Button
            type="primary"
            onClick={this.onSubmit}
            loading={this.props.loading.includes(this.state.sendTokenId)}
          >
            Submit
          </Button>
          <br />
          <Typography.Title level={4}>Need a hand?</Typography.Title>
          <Typography.Text>
            It's a little tough. I get that. Shoot me an email at{' '}
            <PopoutLink url="mailto:sam@canvascbl.com?subject=Help%20To%20Finish%20Signing%20In">
              sam@canvascbl.com
            </PopoutLink>{' '}
            and I'd love to help you out.
          </Typography.Text>
        </Card>
      </div>
    );
  }
}

const ConnectedTokenEntry = connect(state => ({
  token: state.canvas.token,
  successfullySentToken: state.canvas.successfullySentToken,
  subdomain: state.canvas.subdomain,
  loading: state.loading,
  error: state.error
}))(TokenEntry);

export default withCookies(ConnectedTokenEntry);
