import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { Button, Card, List, Skeleton, Typography } from 'antd';
import { Redirect } from 'react-router-dom';
import { parse } from 'qs';
import { isMobile } from 'react-device-detect';
import v4 from 'uuid/v4';
import './index.css';
import {
  setConsentCode,
  getConsentInfo,
  sendConsent
} from '../../../actions/oauth2';
import {
  setGetConsentInfoId,
  setSendConsentId
} from '../../../actions/components/oauth2consent';
import Padding from '../../Padding';
import {
  pageNames,
  trackOAuth2Decision,
  trackPageView
} from '../../../util/tracking';

function OAuth2Consent(props) {
  const {
    location,
    loading,
    getConsentInfoId,
    consentCode,
    consentInfo,
    getConsentInfoError,
    sendConsentId,
    consentRedirectUrl,
    sendConsentError,
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
      trackPageView(pageNames.settings);
    }
  }, [loaded]);

  if (consentRedirectUrl) {
    window.location = consentRedirectUrl;
    return (
      <>
        <Typography.Title level={3}>Redirecting...</Typography.Title>
        <Typography.Text>
          We're redirecting you back to the app that requested permission to
          your account. One sec...
        </Typography.Text>
      </>
    );
  }

  const queryConsentCode = parse(location.search.substr(1)).consent_code;
  if (queryConsentCode && !consentCode) {
    dispatch(setConsentCode(queryConsentCode));
    // removes the query string from the url
    return <Redirect to="/dashboard/authorize" />;
  } else if ((!queryConsentCode && !consentCode) || getConsentInfoError) {
    return (
      <>
        <Typography.Title level={3}>Unable to authorize</Typography.Title>
        <Typography.Text>
          Something's wrong with the app that sent you here. Head back there and
          try authorizing again.
        </Typography.Text>
      </>
    );
  }

  if (
    consentCode &&
    !getConsentInfoId &&
    !consentInfo &&
    !getConsentInfoError
  ) {
    const id = v4();
    dispatch(getConsentInfo(id, consentCode));
    dispatch(setGetConsentInfoId(id));
    return null;
  }

  function handleConsent(allow) {
    const id = v4();
    dispatch(sendConsent(id, consentCode, allow ? 'authorize' : 'deny'));
    dispatch(setSendConsentId(id));
    trackOAuth2Decision(
      consentInfo.credential_name,
      allow,
      consentCode,
      consentInfo.scopes.reduce((acc, val) => acc.concat([val.short_name]), [])
    );
  }

  const cardIsLoading = !consentInfo || loading.includes(getConsentInfoId);
  const buttonsState = {
    loading: loading.includes(sendConsentId),
    disabled: !!sendConsentError || consentRedirectUrl
  };
  const consentError = sendConsentError ? (
    <>
      <Typography.Text type="danger">
        {sendConsentError.error &&
          (() => {
            switch (sendConsentError.error) {
              case 'expired code, restart oauth2 flow':
                return 'This request has timed out. Return to the app that sent you here, then have them send you to CanvasCBL again.';
              case 'invalid consent_code as query param':
                return 'It appears that either CanvasCBL had an unknown error, or you modified the consent_code query param. Retry authorization.';
              default:
                return 'There was an unknown error processing the authorization request. Please return to the app that sent you here.';
            }
          })()}
      </Typography.Text>
      <Padding all={10} />
    </>
  ) : null;

  if (!cardIsLoading && !loaded) {
    setLoaded(true);
  }

  return (
    <>
      <Card
        title={
          cardIsLoading ? (
            <Skeleton paragraph={false} />
          ) : (
            `Authorize ${consentInfo.credential_name}`
          )
        }
        className={!isMobile && 'oauth2consentcard'}
      >
        {cardIsLoading ? (
          <Skeleton active />
        ) : (
          <>
            <Typography.Text>
              <Typography.Text strong>
                {consentInfo.credential_name}
              </Typography.Text>{' '}
              would like to access to following data:
            </Typography.Text>
            <Padding all={10} />
            <List
              itemLayout="horizontal"
              bordered
              dataSource={consentInfo.scopes}
              renderItem={scope => (
                <List.Item key={scope.short_name}>
                  <List.Item.Meta
                    title={scope.name}
                    description={scope.description}
                  />
                </List.Item>
              )}
            />
            <Padding all={10} />
            <Typography.Text>
              <Typography.Text strong>
                {consentInfo.credential_name}
              </Typography.Text>{' '}
              will never have direct access to your Canvas account.
              <br />
              You can revoke{' '}
              <Typography.Text strong>
                {consentInfo.credential_name}
              </Typography.Text>
              's access at any time on the Profile page of CanvasCBL.
              <br />
            </Typography.Text>
            <Padding all={15} />
            {consentError}
            <div style={{ float: 'left' }}>
              <Button onClick={() => handleConsent(false)} {...buttonsState}>
                Cancel
              </Button>
            </div>
            <div style={{ float: 'right' }}>
              <Button
                type="primary"
                onClick={() => handleConsent(true)}
                {...buttonsState}
              >
                Authorize {consentInfo.credential_name}
              </Button>
            </div>
          </>
        )}
      </Card>
    </>
  );
}

export default connect(state => ({
  loading: state.loading,
  consentCode: state.oauth2.consentCode,
  getConsentInfoId: state.components.oauth2consent.getConsentInfoId,
  consentInfo: state.oauth2.consentInfo,
  getConsentInfoError: state.oauth2.getConsentInfoError,
  sendConsentId: state.components.oauth2consent.sendConsentId,
  consentRedirectUrl: state.oauth2.consentRedirectUrl,
  sendConsentError: state.oauth2.sendConsentError
}))(OAuth2Consent);
