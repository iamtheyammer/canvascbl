import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import v4 from 'uuid/v4';
import env from '../../../util/env';
import { Button, Typography, Modal, Spin } from 'antd';
import { getCheckoutSession, getProducts } from '../../../actions/checkout';

const stripe = window.Stripe(env.stripeApiKeyPub);

function Upgrades(props) {
  const { dispatch, error, loading, checkout, user } = props;

  const [getProductsId, setGetProductsId] = useState();

  const [getCheckoutSessionId, setGetCheckoutSessionId] = useState();
  const isLoadingSession = loading.includes(getCheckoutSessionId);
  const sessionError = error[getCheckoutSessionId];

  useEffect(() => {
    if (!checkout.products) {
      const id = v4();
      dispatch(getProducts(id));
      setGetProductsId(id);
    }
    // eslint-disable-next-line
  }, []);

  useEffect(() => {
    if (sessionError) {
      Modal.error({
        title: 'Error starting checkout',
        content: `There was an error opening checkout: ${sessionError.res.data}`
      });
    }
  }, [sessionError]);
  if (!checkout.products || loading.includes(getProductsId)) {
    return (
      <div align="center">
        <Spin />
        <span style={{ paddingTop: '20px' }} />
        <Typography.Title level={3}>{`Loading products...`}</Typography.Title>
      </div>
    );
  }

  if (error[getProductsId]) {
    return (
      <div align="center">
        <Typography.Title level={3}>
          There was an error loading products.
        </Typography.Title>
        <Typography.Text>Please try again later.</Typography.Text>
      </div>
    );
  }

  function handleUpgradeClick(productId) {
    // multiple clicks while loading
    if (getCheckoutSessionId) return;
    const id = v4();
    dispatch(getCheckoutSession(id, productId, user.primary_email));
    setGetCheckoutSessionId(id);
  }

  if (checkout && checkout.session) {
    setGetCheckoutSessionId('');
    stripe
      .redirectToCheckout({ sessionId: checkout.session.session })
      .catch(e =>
        Modal.error({
          title: 'Error redirecting to checkout',
          content: `There was an error from stripe redirecting to checkout: ${e}`
        })
      );
  }

  return (
    <div>
      <Typography.Title level={2}>Upgrades</Typography.Title>
      <Typography.Text type="secondary">
        Take CanvasCBL to a whole new level.
      </Typography.Text>
      <br />
      <div style={{ marginBottom: '15px' }} />
      {checkout.products.map(p => (
        <div key={p.id}>
          <div style={{ marginBottom: '15px' }} />
          <Button
            type="primary"
            disabled={!!sessionError}
            loading={isLoadingSession}
            onClick={() => handleUpgradeClick(p.id)}
          >
            Purchase {p.name} for ${p.price}
          </Button>
        </div>
      ))}
    </div>
  );
}

const ConnectedUpgrades = connect(state => ({
  loading: state.loading,
  error: state.error,
  checkout: state.checkout,
  user: state.canvas.user
}))(Upgrades);

export default ConnectedUpgrades;
