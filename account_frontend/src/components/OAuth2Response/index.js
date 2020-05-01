import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import destToUrl from '../../util/destToUrl';
import MainCard from '../MainCard';
import { parse } from 'qs';
import { notification, Typography } from 'antd';
import { Redirect } from 'react-router-dom';
import { setDestination, setRedirectOk } from '../../actions/home';

function OAuth2Response(props) {
  const { destination, redirectOk, location, dispatch } = props;

  const { name, dest, error } = parse(
    location.search.slice(1, location.search.length)
  );
  const destUrl = destToUrl(dest);

  useEffect(() => {
    if (!destination && dest) {
      dispatch(setDestination(destination));
    }
  });

  useEffect(() => {
    if (!error && destination && !redirectOk) {
      dispatch(setRedirectOk(true));
    }
  });

  if (!name && !dest && !error) {
    return <Redirect to="/" />;
  }

  if (error) {
    if (error === 'access_denied') {
      return <Redirect to="/" />;
    } else {
      notification.error({
        message: 'Unknown Error',
        duration: 0,
        description:
          'There was an unknown error logging you in with Canvas. Try again later.'
      });

      return <Redirect to="/" />;
    }
  }

  const firstName = name && name.split(' ')[0];
  if (!error && dest && destUrl) {
    return (
      <MainCard
        loading
        loadingText={
          <Typography.Title level={2}>
            Welcome{firstName ? `, ${firstName}` : ' to CanvasCBL!'}
          </Typography.Title>
        }
      />
    );
  }

  return (
    <MainCard>
      <Typography.Text type="danger">
        An unexpected error occured when signing you in.
      </Typography.Text>
    </MainCard>
  );
}

const ConnectedOAuth2Response = connect((state) => ({
  destination: state.home.destination,
  redirectOk: state.home.redirectOk
}))(OAuth2Response);

export default ConnectedOAuth2Response;
