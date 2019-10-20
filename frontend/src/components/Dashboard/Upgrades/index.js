import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import v4 from 'uuid/v4'
import env from '../../../util/env';
import {
  Button,
  Typography,
  Modal
} from 'antd';
import {getCheckoutSession} from "../../../actions/checkout";

const stripe = window.Stripe(env.stripeApiKeyPub);

function Upgrades(props) {
  const { dispatch, error, loading, checkout } = props;

  const [ getCheckoutSessionId, setGetCheckoutSessionId ] = useState();
  const isLoadingSession = loading.includes(getCheckoutSessionId);
  const sessionError = error[getCheckoutSessionId];

  if(sessionError) {
    Modal.error({
      title: 'Error starting checkout',
      content: `There was an error opening checkout: ${sessionError}`
    });
  }

  function handleUpgradeClick(e) {
    console.log('clicc')
    // multiple clicks while loading
    if(getCheckoutSessionId) return;
    const id = v4();
    dispatch(getCheckoutSession(id));
    setGetCheckoutSessionId(id)
  }

  if(checkout && checkout.checkoutSession) {
    stripe.redirectToCheckout({ sessionId: checkout.checkoutSession })
      .catch(e => Modal.error({
        title: 'Error redirecting to checkout',
        content: `There was an error from stripe redirecting to checkout: ${e}`
      }))
  }

  return (
    <div>
      <Typography.Title level={2} >Upgrades</Typography.Title>
      <Typography.Text type="secondary">Take CanvasCBL to a whole new level.</Typography.Text>
      <br />
      <div style={{ marginBottom: '15px' }} />
      <Button
        type="primary"
        disabled={!!sessionError}
        loading={isLoadingSession}
        onClick={handleUpgradeClick}
      >Purchase</Button>
    </div>
  )
}

const ConnectedUpgrades = connect(state => ({
  loading: state.loading,
  error: state.error,
  checkout: state.checkout
}))(Upgrades);

export default ConnectedUpgrades;