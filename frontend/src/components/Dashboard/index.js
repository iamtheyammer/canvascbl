import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { Switch, Route, Redirect, Link } from 'react-router-dom';
import * as ReactGA from 'react-ga';
import v4 from 'uuid/v4';
import { useCookies } from 'react-cookie';
import { isMobile } from 'react-device-detect';

import {
  Layout,
  Breadcrumb,
  Popover,
  Typography,
  Spin,
  Modal,
  Icon,
  notification
} from 'antd';

import DashboardNav from './DashboardNav';
import ConnectedUserProfile from './UserProfile';
import ConnectedGrades from './Grades';
import ConnectedGradeBreakdown from './Grades/GradeBreakdown';
import ConnectedUpgrades from './Upgrades';
import ConnectedLogout from './Logout';
import UpdateHandler from './UpdateHandler';
import env from '../../util/env';
import { getCanvasToken, getUser, logout } from '../../actions/canvas';
import ConnectedErrorModal from './ErrorModal';
import { getSessionInformation } from '../../actions/plus';
import './index.css';

const { Content, Footer } = Layout;

const getBreadcrumbNameMap = (courses = []) => {
  const routes = {
    '/dashboard': 'Dashboard',
    '/dashboard/profile': 'Profile',
    '/dashboard/grades': 'Grades',
    '/dashboard/upgrades': 'Upgrades'
  };

  courses.forEach(
    c => (routes[`/dashboard/grades/${c.id}`] = `Grade Breakdown for ${c.name}`)
  );

  return routes;
};

function Dashboard(props) {
  const { token } = props;
  const [cookies] = useCookies(['session_string']);

  const [hasSentUserToGa, setHasSentUserToGa] = useState(false);
  const [getUserId, setGetUserId] = useState();
  const [getSessionId, setGetSessionId] = useState();

  const {
    location,
    user,
    subdomain,
    session,
    loading,
    error,
    dispatch
  } = props;

  useEffect(() => {
    ReactGA.pageview(
      props.location.pathname +
        (props.location.search.includes('~') ? '' : props.location.search)
    );
  }, [props.location]);

  useEffect(() => {
    if (session && session.status === 1) {
      notification.error({
        message: 'CanvasCBL is disabled for this user.',
        description:
          'CanvasCBL is disabled for this user. Please contact your school for more information.'
      });
    }
  }, [session]);

  // if no token exists, redirect
  if (!localStorage.token && !cookies.session_string) {
    return <Redirect to="/" />;
  } else if (localStorage.token && !token) {
    // otherwise, wait for token
    return null;
  }

  const pathSnippets = location.pathname.split('/').filter(i => i);
  const breadcrumbNameMap = getBreadcrumbNameMap(props.courses || []);
  const breadcrumbItems = pathSnippets.map((_, index) => {
    const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;
    return (
      <Breadcrumb.Item key={url}>
        <Link to={url}>{breadcrumbNameMap[url]}</Link>
      </Breadcrumb.Item>
    );
  });

  if (hasSentUserToGa === false && user) {
    ReactGA.set({ userId: user.id });
    setHasSentUserToGa(true);
  }

  if (token && !user && !getUserId) {
    const id = v4();
    dispatch(getUser(id, !cookies['session_string'], token, subdomain));
    setGetUserId(id);
  }

  if (token && user && !session && !getSessionId) {
    const id = v4();
    dispatch(getSessionInformation(id));
    setGetSessionId(id);
  }

  if (error[getUserId] || error[getSessionId]) {
    return <ConnectedErrorModal error={error[getUserId]} />;
  }

  if (session && session.status === 1) {
    return <Redirect to={'/dashboard/logout'} />;
  }

  if (user && session && user.primary_email !== session.email) {
    Modal.confirm({
      icon: <Icon type="exclamation-circle" style={{ color: '#D8000C' }} />,
      title: 'Session mismatch',
      content:
        "There's a deep issue with your sessions or you're trying to use someone else's subscription. " +
        'Either way, contact support to try again.',
      okText: 'Logout',
      cancelText: 'Logout',
      onCancel: () => dispatch(logout(token, subdomain)),
      onOk: () => dispatch(logout(token, subdomain))
    });
    return null;
  }

  const routes = (
    <Switch>
      <Route
        exact
        path="/dashboard"
        render={() => <Redirect to="/dashboard/grades" />}
      />
      <Route exact path="/dashboard/profile" component={ConnectedUserProfile} />
      <Route exact path="/dashboard/grades" component={ConnectedGrades} />
      <Route
        exact
        path="/dashboard/grades/:courseId"
        component={ConnectedGradeBreakdown}
      />
      <Route exact path="/dashboard/upgrades" component={ConnectedUpgrades} />
      <Route exact path="/dashboard/logout" component={ConnectedLogout} />
      <Route render={() => <Redirect to="/" />} />
    </Switch>
  );

  if (isMobile) {
    return (
      <DashboardNav>
        <div
          style={{
            background: '#ffffff',
            padding: '8px 8px 12px 8px',
            marginRight: '8px',
            height: 'auto'
          }}
        >
          {subdomain &&
          token &&
          !loading.includes(getUserId) &&
          !loading.includes(getSessionId)
            ? routes
            : loading}
        </div>
      </DashboardNav>
    );
  }

  return (
    <div className="dashboard">
      <Layout className="layout">
        <DashboardNav />
        <Content style={{ padding: '0 50px' }}>
          <Breadcrumb style={{ marginTop: 12 }}>{breadcrumbItems}</Breadcrumb>
          <div
            style={{
              background: '#fff',
              padding: 24,
              marginTop: 12,
              minHeight: 280
            }}
          >
            {token &&
            (!loading.includes(getUserId) ||
              !loading.includes(getSessionId)) ? (
              routes
            ) : (
              <div align="center">
                <Spin />
                <span style={{ paddingTop: '20px' }} />
                <Typography.Title
                  level={3}
                >{`Loading your user...`}</Typography.Title>
              </div>
            )}
          </div>
        </Content>
        <UpdateHandler />
        <Footer style={{ textAlign: 'center' }}>
          <Popover
            trigger="click"
            content={
              <Typography.Text>
                Version {env.currentVersion}
                {env.nodeEnv === 'development' && '-DEV'}
              </Typography.Text>
            }
          >
            Built by iamtheyammer 2019
          </Popover>
        </Footer>
      </Layout>
    </div>
  );
}

const ConnectedDashboard = connect(state => ({
  token: state.canvas.token,
  subdomain: state.canvas.subdomain,
  courses: state.canvas.courses,
  user: state.canvas.user,
  session: state.plus.session,
  loading: state.loading,
  error: state.error
}))(Dashboard);

export default ConnectedDashboard;
