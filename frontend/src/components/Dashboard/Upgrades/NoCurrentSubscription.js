import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { isMobile } from 'react-device-detect';

import { Typography, Col, Row, Card, Icon, Modal, Button } from 'antd';
import {
  Button as MobileButton,
  Card as MobileCard,
  WhiteSpace as MobileWhiteSpace
} from 'antd-mobile';

import { ReactComponent as plusIcon } from '../../../assets/plus.svg';
import v4 from 'uuid/v4';
import { getCheckoutSession } from '../../../actions/checkout';
import env from '../../../util/env';
import { chunk } from 'lodash';

import averageGradeImg from './average-grade.png';
import averageOutcomeScoreImg from './average-outcome-score.png';
import howToGetAnAImg from './how-to-get-an-a.png';
import previousGradeImg from './previous-grade.png';
import logoNavbarImg from './logo-navbar.png';

const stripe = window.Stripe(env.stripeApiKeyPub);

const benefits = [
  {
    title: 'Average Grades',
    img: averageGradeImg,
    content: "See the average grade for any class you're in."
  },
  {
    title: 'How to get an A',
    img: howToGetAnAImg,
    content:
      "For every class you haven't mastered yet, get a step-by-step list of things to do to get an A."
  },
  {
    title: 'Average Outcome Scores',
    img: averageOutcomeScoreImg,
    content: 'See an average score for every outcome in every class.'
  },
  {
    title: 'Previous Grades',
    img: previousGradeImg,
    content:
      'See how your grades have changed from your last login to now. ' +
      'Hover over any grade to see when it was from, so you can better ' +
      'understand your progression in your courses.'
  },
  {
    title: 'CanvasCBL+ Logo',
    img: logoNavbarImg,
    content: (
      <div>
        Get reminded of your awesomeness every time you log in-- the logo at the
        top left will show a little <Icon component={plusIcon} />.
      </div>
    )
  }
];

function NoCurrentSubscription(props) {
  const { dispatch, checkout, error, loading } = props;

  const [getCheckoutSessionId, setGetCheckoutSessionId] = useState();
  const sessionError = error[getCheckoutSessionId];

  const checkoutSessionIsLoading = !!loading[getCheckoutSessionId];

  useEffect(() => {
    if (sessionError) {
      Modal.error({
        title: 'Error starting checkout',
        content: `There was an error opening checkout: ${sessionError.res.data}`
      });
    }
  }, [sessionError]);

  function handleUpgradeClick() {
    // multiple clicks while loading
    if (getCheckoutSessionId) return;
    const id = v4();
    dispatch(getCheckoutSession(id, env.upgradesPurchasableProductId));
    setGetCheckoutSessionId(id);
  }

  if (!checkout.products) {
    return null;
  }

  const product = checkout.products.filter(
    p => p.id === env.upgradesPurchasableProductId
  )[0];

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

  if (isMobile) {
    return (
      <div>
        <Typography.Title level={2}>Upgrades</Typography.Title>
        <Typography.Text type="secondary">
          Take CanvasCBL to a whole new level with CanvasCBL+.
        </Typography.Text>
        <Typography.Title level={3}>Benefits</Typography.Title>
        {benefits.map(b => (
          <div key={b.title}>
            <MobileCard>
              <MobileCard.Header title={b.title} />
              <MobileCard.Body>
                <img src={b.img} alt={b.title} style={{ maxWidth: '100%' }} />
                {b.content}
              </MobileCard.Body>
            </MobileCard>
            <MobileWhiteSpace />
          </div>
        ))}
        <Typography.Text>
          Ready to take the leap? Click below to subscribe! New users get a 7
          day free trial, so there's no reason not to get started! Purchasing
          CanvasCBL+ really helps cover expenses so we can keep CanvasCBL free
          for everyone.
        </Typography.Text>
        <MobileWhiteSpace />
        Click below to purchase {product.name} for ${product.price}/month.
        <MobileWhiteSpace />
        <MobileButton
          type="primary"
          disabled={!!sessionError}
          loading={checkoutSessionIsLoading}
          onClick={() => handleUpgradeClick()}
        >
          Purchase
        </MobileButton>
      </div>
    );
  }

  // splits [1,2,3,4,5,6] into [[1,2,3],[4,5,6]]
  const chunkedBenefits = chunk(benefits, 3);

  return (
    <div>
      <Typography.Title level={2}>Upgrades</Typography.Title>
      <Typography.Text type="secondary">
        Take CanvasCBL to a whole new level with CanvasCBL+.
      </Typography.Text>
      <Typography.Title level={3}>Benefits</Typography.Title>
      {/* using v4() as keys here because there will be no updates */}
      {chunkedBenefits.map(bs => (
        <div key={v4()}>
          <Row gutter={20}>
            {bs.map(b => (
              <Col span={8} key={b.title}>
                <Card title={b.title} cover={<img src={b.img} alt={b.title} />}>
                  <Typography.Text>{b.content}</Typography.Text>
                </Card>
              </Col>
            ))}
          </Row>
          <div style={{ padding: '10px' }} />
        </div>
      ))}

      <div style={{ padding: '10px' }} />
      <Typography.Title level={3}>Get Started</Typography.Title>
      <Typography.Text>
        Ready to take the leap? Click below to subscribe! New users get a 7 day
        free trial, so there's no reason not to get started! Purchasing
        CanvasCBL+ really helps cover expenses so we can keep CanvasCBL free for
        everyone.
      </Typography.Text>
      <div style={{ padding: '10px' }} />
      <Button
        type="primary"
        disabled={!!sessionError}
        loading={checkoutSessionIsLoading}
        onClick={() => handleUpgradeClick()}
      >
        Purchase {product.name} for ${product.price}/month
      </Button>
    </div>
  );
}

const ConnectedNoCurrentSubscription = connect(state => ({
  checkout: state.checkout,
  loading: state.loading,
  error: state.error
}))(NoCurrentSubscription);

export default ConnectedNoCurrentSubscription;
