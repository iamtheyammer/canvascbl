import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import v4 from 'uuid/v4';
import { isMobile } from 'react-device-detect';

import {
  Typography,
  Descriptions,
  Spin,
  Popconfirm,
  Button,
  Popover
} from 'antd';
import { Button as MobileButton } from 'antd-mobile';
import {
  cancelSubscription,
  getSubscriptions
} from '../../../actions/checkout';
import moment from 'moment';

function HasCurrentSubscription(props) {
  const { dispatch, user, checkout, loading, error } = props;

  const [getSubscriptionsId, setGetSubscriptionsId] = useState('');
  const subscriptionsAreLoading = loading.includes(getSubscriptionsId);
  const getSubscriptionsError = error[getSubscriptionsId];

  const [cancelSubscriptionId, setCancelSubscriptionId] = useState('');
  const cancelSubscriptionIsLoading = loading.includes(cancelSubscriptionId);
  const cancelSubscriptionError = error[cancelSubscriptionId];
  const hasCanceled =
    !!cancelSubscriptionId &&
    !cancelSubscriptionIsLoading &&
    !cancelSubscriptionError;

  useEffect(() => {
    if (!checkout.subscriptions && !getSubscriptionsId) {
      const id = v4();
      dispatch(getSubscriptions(id));
      setGetSubscriptionsId(id);
    }
  }, [checkout.subscriptions, getSubscriptionsId, dispatch]);

  if (getSubscriptionsError) {
    return <div>Error getting subscriptions. Please try again later.</div>;
  }

  if (subscriptionsAreLoading || !checkout.subscriptions) {
    return (
      <div align="center">
        <Spin />
        <span style={{ paddingTop: '20px' }} />
        <Typography.Title
          level={3}
        >{`Loading subscriptions...`}</Typography.Title>
      </div>
    );
  }

  if (cancelSubscriptionError) {
    return (
      <div>
        Error cancelling your subscription. Please try again or contact support.
      </div>
    );
  }

  async function handleCancel() {
    const id = v4();
    dispatch(cancelSubscription(id));
    setCancelSubscriptionId(id);
  }

  const sub = checkout.subscriptions[0];

  return (
    <div>
      <Typography.Title level={2}>Upgrades</Typography.Title>
      <Typography.Text type="subtitle">
        Thank you very much for being a CanvasCBL+ subscriber,{' '}
        {user.name.split(' ')[0]}!
      </Typography.Text>
      <Typography.Title level={3}>Subscription Info</Typography.Title>
      <Descriptions>
        <Descriptions.Item label="Status">{sub.status}</Descriptions.Item>
        <Descriptions.Item label="Current Period Started">
          {moment.unix(sub.currentPeriodStart).calendar()}
        </Descriptions.Item>
        {sub.trialEnd && (
          <Descriptions.Item label="Trial End">
            {moment.unix(sub.trialEnd).calendar()}
          </Descriptions.Item>
        )}
        <Descriptions.Item label="Next Bill Date">
          {moment.unix(sub.currentPeriodEnd).calendar()}
        </Descriptions.Item>
        <Descriptions.Item label="Subscriber Since">
          {moment.unix(sub.insertedAt).calendar()}
        </Descriptions.Item>
      </Descriptions>
      <Typography.Title level={3}>Manage Subscription</Typography.Title>
      <Popover content="Cancel your current subscription, then enter your gift card code.">
        <Typography.Text type={'link'}>
          Have a gift card or promotional code?
        </Typography.Text>
      </Popover>
      <br />
      <Typography.Text>
        Canceling your CanvasCBL+ subscription takes effect immediately.
      </Typography.Text>
      <div style={{ padding: '10px' }} />
      {!hasCanceled ? (
        <Popconfirm
          title="Are you sure you want to cancel? We'll miss you!"
          placement="top"
          onConfirm={handleCancel}
        >
          {isMobile ? (
            <MobileButton loading={cancelSubscriptionIsLoading} type="warning">
              Cancel Your Subscription
            </MobileButton>
          ) : (
            <Button loading={cancelSubscriptionIsLoading} type="danger">
              Cancel your CanvasCBL+ Subscription
            </Button>
          )}
        </Popconfirm>
      ) : (
        <Typography.Text>Successfully canceled.</Typography.Text>
      )}
    </div>
  );
}

const ConnectedHasValidSubscription = connect(state => ({
  checkout: state.checkout,
  user: state.canvas.user,
  loading: state.loading,
  error: state.error
}))(HasCurrentSubscription);

export default ConnectedHasValidSubscription;
