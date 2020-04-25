import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { Button, Checkbox, Col, Divider, Row, Typography } from 'antd';
import {
  Button as MobileButton,
  Checkbox as MobileCheckbox
} from 'antd-mobile';
import { isMobile } from 'react-device-detect';
import { CaretLeftOutlined, CaretRightOutlined } from '@ant-design/icons';
import styled from 'styled-components';
import { parse } from 'qs';
import env from '../../util/env';
import banner from '../../assets/banner.svg';
import PopoutLink from '../PopoutLink';
import Padding from '../Padding';
import {
  setDestination,
  setRedirectOk,
  setSignInButtonAvailability
} from '../../actions/home';
import { validateDest } from '../../util/destToUrl';
import getUrlPrefix from '../../util/getUrlPrefix';
import MainCard from '../MainCard';

const CenterText = styled.div`
  text-align: center;
`;

const CenteredCheckbox = styled(Checkbox)`
  display: flex;
  align-items: center;
  justify-content: center;
`;

const CenteredButton = styled(Button)`
  left: 50%;
  transform: translate(-50%);
`;

const ButtonFloatRight = styled(Button)`
  float: right;
`;

function Home(props) {
  const {
    redirectOk,
    profile,
    loadingUserProfile,
    signInButtonAvailability,
    destination,
    location,
    dispatch
  } = props;

  const { dest } = parse(location.search.slice(1, location.search.length));
  const destIsTeacher = destination === 'teacher';

  useEffect(() => {
    if (dest && !destination) {
      dispatch(setDestination(validateDest(dest) ? dest : 'canvascbl'));
    }
  });

  useEffect(() => {
    if (!redirectOk) {
      dispatch(setRedirectOk(!!profile));
    }
  }, [dispatch, redirectOk, profile]);

  if (isMobile) {
    return (
      <div
        style={{
          textAlign: 'center',
          margin: '10px 5px 10px 5px',
          backgroundColor: '#ffffff'
        }}
      >
        <img src={banner} alt={'banner'} />
        <Typography.Title level={2}>
          Log in to CanvasCBL{destIsTeacher && ' for Teachers'}
        </Typography.Title>
        <Typography.Text>
          To use CanvasCBL, please accept the terms and sign in with Canvas.
        </Typography.Text>
        <MobileCheckbox.AgreeItem
          onChange={(e) =>
            dispatch(setSignInButtonAvailability(e.target.checked))
          }
        >
          I accept the{' '}
          <PopoutLink url={env.privacyPolicyUrl}>privacy policy</PopoutLink> and{' '}
          <PopoutLink url={env.termsOfServiceUrl}>terms of service</PopoutLink>
        </MobileCheckbox.AgreeItem>
        <MobileButton
          type="primary"
          disabled={!signInButtonAvailability}
          loading={loadingUserProfile}
          onClick={() =>
            (window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request?intent=auth&dest=${destination}`)
          }
        >
          Sign in with Canvas
        </MobileButton>
      </div>
    );
  }

  return (
    <MainCard loading={loadingUserProfile}>
      <>
        <CenterText>
          <Typography.Title level={2}>
            Log in to CanvasCBL{destIsTeacher && ' for Teachers'}
          </Typography.Title>
          {env.buildBranch !== 'master' && (
            <Typography.Text type="danger">
              CanvasCBL is running in {env.buildBranch}
              <br />
            </Typography.Text>
          )}
          <Typography.Text>
            To use CanvasCBL, please accept the terms and sign in with Canvas.
          </Typography.Text>
          <Padding all={5} />
        </CenterText>
        <div>
          <CenteredCheckbox
            onChange={(e) =>
              dispatch(setSignInButtonAvailability(e.target.checked))
            }
            className="center-checkbox"
          >
            I accept the{' '}
            <PopoutLink url={env.privacyPolicyUrl}>privacy policy</PopoutLink>{' '}
            and{' '}
            <PopoutLink url={env.termsOfServiceUrl}>
              terms of service
            </PopoutLink>
          </CenteredCheckbox>
          <Padding all={8} />
          <>
            <CenteredButton
              type="primary"
              size={!signInButtonAvailability ? 'default' : 'large'}
              loading={loadingUserProfile}
              disabled={!signInButtonAvailability}
              onClick={() =>
                (window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request?intent=auth&dest=${destination}`)
              }
            >
              Sign in with Canvas
            </CenteredButton>
          </>
        </div>
      </>
      <Divider />
      <Row>
        <Col span={8} flex>
          <Button
            type="link"
            onClick={() => (window.location = 'https://canvascbl.com')}
            style={{ float: 'left' }}
          >
            <CaretLeftOutlined /> Back home
          </Button>
        </Col>

        <Col span={16} flex>
          {destIsTeacher ? (
            <ButtonFloatRight
              type="link"
              onClick={() => dispatch(setDestination('canvascbl'))}
              style={{ float: 'right' }}
            >
              CanvasCBL for Students and Parents <CaretRightOutlined />
            </ButtonFloatRight>
          ) : (
            <ButtonFloatRight
              type="link"
              onClick={() => dispatch(setDestination('teacher'))}
              style={{ float: 'right' }}
            >
              CanvasCBL for Teachers <CaretRightOutlined />
            </ButtonFloatRight>
          )}
        </Col>
      </Row>
    </MainCard>
  );
}

const ConnectedHome = connect((state) => ({
  profile: state.canvas.profile,
  loadingUserProfile: state.canvas.loadingUserProfile,
  signInButtonAvailability: state.home.signInButtonAvailability,
  destination: state.home.destination,
  redirectOk: state.home.redirectOk
}))(Home);

export default ConnectedHome;
