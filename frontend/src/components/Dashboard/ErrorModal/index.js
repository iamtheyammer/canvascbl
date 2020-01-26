import React, { Component } from 'react';
import { shape, object } from 'prop-types';
import { Redirect } from 'react-router-dom';
import { Modal, Icon } from 'antd';
import { logout } from '../../../actions/canvas';

class ErrorModal extends Component {
  constructor(props) {
    super(props);
    this.state = {
      redirect: false
    };
  }

  handleOk = () => {
    window.location.reload();
  };

  handleCancel = () => {
    this.props.dispatch(logout());
    this.setState({ redirect: true });
  };

  handleUnknown = () => {
    Modal.confirm({
      icon: <Icon type="exclamation-circle" style={{ color: '#D8000C' }} />,
      title: 'Unknown Error',
      content: `Do you want to reload?`,
      okText: 'Try again',
      cancelText: 'Logout',
      onCancel: this.handleCancel,
      onOk: this.handleOk
    });
  };

  componentDidMount() {
    const { res, error } = this.props;
    const result = res || error.res;
    if (!result) {
      this.handleUnknown();
      return;
    }

    const canvasStatusCode = parseInt(result.headers['x-canvas-status-code']);

    if (canvasStatusCode === 401) {
      const message = result.data.errors[0].message;

      if (message === 'Insufficient scopes on access token.') {
        Modal.info({
          title: 'Re-login required',
          content:
            "We've added some cool new features to CanvasCBL that require you to log out and log back in.",
          closable: false,
          icon: <Icon type="exclamation-circle" style={{ color: '#D8000C' }} />,
          okText: 'Logout',
          onOk: this.handleCancel
        });
        return;
      }

      Modal.info({
        title: 'Invalid Canvas Token',
        content:
          "There's an issue with your Canvas token or subdomain. Click Logout to enter a different one.",
        closable: false,
        icon: <Icon type="exclamation-circle" style={{ color: '#D8000C' }} />,
        okText: 'Logout',
        onOk: this.handleCancel
      });
      return <Redirect to="/" />;
    }
  }

  render() {
    if (this.state.redirect) {
      return <Redirect to="/" />;
    }

    return <div />;
  }
}

ErrorModal.propTypes = {
  error: shape({
    res: object
  })
};

export default ErrorModal;
