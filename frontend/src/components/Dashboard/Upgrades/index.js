import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import v4 from 'uuid/v4';
import { Typography, Spin } from 'antd';
import { getProducts } from '../../../actions/checkout';
import ConnectedHasValidSubscription from './HasCurrentSubscription';
import ConnectedNoCurrentSubscription from './NoCurrentSubscription';

function Upgrades(props) {
  const { dispatch, error, loading, checkout, session } = props;

  const [getProductsId, setGetProductsId] = useState();

  useEffect(() => {
    if (!checkout.products) {
      const id = v4();
      dispatch(getProducts(id));
      setGetProductsId(id);
    }
    // eslint-disable-next-line
  }, []);

  if (!session || !checkout.products || loading.includes(getProductsId)) {
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

  return (
    <div>
      {session.has_valid_subscription ? (
        <ConnectedHasValidSubscription />
      ) : (
        <ConnectedNoCurrentSubscription />
      )}
    </div>
  );
}

const ConnectedUpgrades = connect(state => ({
  loading: state.loading,
  error: state.error,
  checkout: state.checkout,
  user: state.canvas.user,
  session: state.plus.session
}))(Upgrades);

export default ConnectedUpgrades;
