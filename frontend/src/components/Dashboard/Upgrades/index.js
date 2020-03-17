import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { Typography } from 'antd';
import { pageNames, trackPageView } from '../../../util/tracking';

function Upgrades(props) {
  // const { dispatch, error, loading, checkout, session } = props;

  // const [getProductsId, setGetProductsId] = useState();
  //
  // useEffect(() => {
  //   if (!checkout.products) {
  //     const id = v4();
  //     dispatch(getProducts(id));
  //     setGetProductsId(id);
  //   }
  //   // eslint-disable-next-line
  // }, []);

  const [loaded, setLoaded] = useState(false);
  useEffect(() => {
    /*
    This system is to prevent sending tons of Page View events to Mixpanel.
    Those tons of events are sent because, every time the state changes,
    this component is rerendered. The most common state change is when
    a grade average loads in for plus users.

    It works with two hooks: state and effect.

    There's a loaded state hook set to false just above.

    The effect hook is used here to run whenever loaded changes--
    if it's true, we'll track a page view. If not, whatever.

    The reason that this works is because state is reset on unmount.
    So we only get one page view per actual page view.
     */

    if (loaded) {
      trackPageView(pageNames.upgrades);
    }
  }, [loaded]);

  if (!loaded) {
    setLoaded(true);
  }

  return (
    <>
      <Typography.Title level={2}>Upgrades</Typography.Title>
      <Typography.Text strong>
        CanvasCBL+ is on us while we're doing distance learning.
      </Typography.Text>
      <br />
      <br />
      <Typography.Text>Stay safe out there!</Typography.Text>
    </>
  );

  // if (!session || !checkout.products || loading.includes(getProductsId)) {
  //   return (
  //     <div align="center">
  //       <Spin />
  //       <span style={{ paddingTop: '20px' }} />
  //       <Typography.Title level={3}>{`Loading products...`}</Typography.Title>
  //     </div>
  //   );
  // }
  //
  // if (error[getProductsId]) {
  //   return (
  //     <div align="center">
  //       <Typography.Title level={3}>
  //         There was an error loading products.
  //       </Typography.Title>
  //       <Typography.Text>Please try again later.</Typography.Text>
  //     </div>
  //   );
  // }
  //
  // if (!loaded) {
  //   setLoaded(true);
  // }
  //
  // return (
  //   <div>
  //     {session.has_valid_subscription ? (
  //       <ConnectedHasValidSubscription />
  //     ) : (
  //       <ConnectedNoCurrentSubscription />
  //     )}
  //   </div>
  // );
}

const ConnectedUpgrades = connect(state => ({
  loading: state.loading,
  error: state.error,
  checkout: state.checkout,
  user: state.canvas.user,
  session: state.plus.session
}))(Upgrades);

export default ConnectedUpgrades;
