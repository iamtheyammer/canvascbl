import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import v4 from 'uuid/v4';

import { Typography, Input, Button, Row, List, Result } from 'antd';
import {
  addGiftCard,
  redeemGiftCards,
  removeGiftCard,
  updateGiftCardEntryError,
  updateGiftCardField1,
  updateGiftCardField2,
  updateGiftCardField3
} from '../../../../actions/components/redeem';
import Padding from '../../../Padding';
import moment from 'moment';
import { pageNames, trackPageView } from '../../../../util/tracking';

const giftCardRegex = /[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}/;

function Redeem(props) {
  const {
    session,
    loading,
    error,
    giftCardField1,
    giftCardField2,
    giftCardField3,
    giftCards,
    giftCardEntryError,
    redeemGiftCardsId,
    redeemGiftCardsResponse,
    dispatch
  } = props;

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
      trackPageView(pageNames.redeem);
    }
  }, [loaded]);

  if (session && session.has_valid_subscription === true) {
    return <Redirect to="/dashboard/upgrades" />;
  }

  function handleAddGiftCard() {
    const card = `${giftCardField1}-${giftCardField2}-${giftCardField3}`;

    if (!giftCardRegex.test(card)) {
      dispatch(
        updateGiftCardEntryError(
          "Your entry doesn't look like a gift card code." +
            ' A gift card code should be three sets of four uppercase letters and/or numbers.'
        )
      );
      return;
    }

    if (giftCards && giftCards.includes(card)) {
      dispatch(
        updateGiftCardEntryError('You have already added that gift card code.')
      );
      return;
    }

    dispatch(addGiftCard(card.toUpperCase()));
  }

  if (!loaded) {
    setLoaded(true);
  }

  if (redeemGiftCardsResponse && redeemGiftCardsResponse.success === true) {
    return (
      <div>
        <Row>
          <Typography.Title level={2}>Redeem</Typography.Title>
        </Row>
        <Row>
          <Result
            status="success"
            title="Successfully activated CanvasCBL+!"
            subTitle={`Your subscription is set to expire on/at ${moment
              .unix(redeemGiftCardsResponse.subscription_expires_at)
              .calendar()}.`}
            extra={[
              <Button
                type="primary"
                key="reload"
                onClick={() => window.location.reload()}
              >
                Get started!
              </Button>
            ]}
          />
        </Row>
      </div>
    );
  }

  return (
    <div>
      <Row>
        <Typography.Title level={2}>Redeem</Typography.Title>
      </Row>
      <Row>
        <Typography.Text type="secondary">
          Redeem a gift card here.
        </Typography.Text>
      </Row>
      <Padding br />
      <Row>
        <Typography.Title level={3}>Code Entry</Typography.Title>
        <Typography.Text type="secondary">
          Enter your gift card code below.
        </Typography.Text>
      </Row>
      <Row style={{ paddingTop: 10 }}>
        <Input.Group compact>
          <Input
            style={{ width: 75, textAlign: 'center' }}
            placeholder="A1CD"
            onChange={(e) =>
              dispatch(updateGiftCardField1(e.target.value.toUpperCase()))
            }
            value={giftCardField1}
            maxLength={4}
            onPressEnter={handleAddGiftCard}
          />
          <Input
            style={{
              width: 30,
              borderLeft: 0,
              borderRight: 0,
              pointerEvents: 'none',
              backgroundColor: '#fff'
            }}
            placeholder="-"
            disabled
          />
          <Input
            style={{ width: 75, textAlign: 'center', borderLeft: 0 }}
            placeholder="EF2H"
            onChange={(e) =>
              dispatch(updateGiftCardField2(e.target.value.toUpperCase()))
            }
            value={giftCardField2}
            maxLength={4}
            onPressEnter={handleAddGiftCard}
          />
          <Input
            style={{
              width: 30,
              borderLeft: 0,
              pointerEvents: 'none',
              backgroundColor: '#fff'
            }}
            placeholder="-"
            disabled
          />
          <Input
            style={{ width: 75, textAlign: 'center', borderLeft: 0 }}
            placeholder="IJK3"
            onChange={(e) =>
              dispatch(updateGiftCardField3(e.target.value.toUpperCase()))
            }
            value={giftCardField3}
            maxLength={4}
            onPressEnter={handleAddGiftCard}
          />
          <Button type="default" onClick={handleAddGiftCard}>
            Add
          </Button>
        </Input.Group>
      </Row>
      {giftCardEntryError && (
        <Row>
          <Typography.Text type="danger">{giftCardEntryError}</Typography.Text>
        </Row>
      )}
      <Row>
        <Padding br />
        <Padding bottom={10} />
        <Typography.Title level={3}>Added Cards</Typography.Title>
        <Typography.Text type="secondary">
          All codes you've added will appear below. Click Submit at the bottom
          of the page to redeem your codes.
        </Typography.Text>
        <List
          itemLayout="horizontal"
          dataSource={giftCards}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Button
                  key="remove"
                  onClick={() => dispatch(removeGiftCard(item))}
                >
                  Remove
                </Button>
              ]}
            >
              <List.Item.Meta title={item} />
            </List.Item>
          )}
        />
      </Row>
      <Row>
        <Padding br />
        <Typography.Title level={3}>Submit</Typography.Title>
        <Typography.Text type="secondary">
          Click the button below to submit your codes.{' '}
        </Typography.Text>
        <Padding br />
        <Padding br />
        <Button
          type="primary"
          disabled={!giftCards || !giftCards.length}
          loading={loading.includes(redeemGiftCardsId)}
          onClick={() => dispatch(redeemGiftCards(v4(), giftCards))}
        >
          Submit
        </Button>
      </Row>
      {error[redeemGiftCardsId] && (
        <Row>
          <Typography.Text type="danger">
            {error[redeemGiftCardsId].res.error}
          </Typography.Text>
        </Row>
      )}
    </div>
  );
}

const ConnectedRedeem = connect((state) => ({
  session: state.plus.session,
  loading: state.loading,
  error: state.error,
  giftCardField1: state.components.redeem.giftCardField1,
  giftCardField2: state.components.redeem.giftCardField2,
  giftCardField3: state.components.redeem.giftCardField3,
  giftCardEntryError: state.components.redeem.giftCardEntryError,
  giftCards: state.components.redeem.giftCards,
  redeemGiftCardsId: state.components.redeem.redeemGiftCardsId,
  redeemGiftCardsResponse: state.components.redeem.redeemGiftCardsResponse
}))(Redeem);

export default ConnectedRedeem;
