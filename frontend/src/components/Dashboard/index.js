import React, { useState } from 'react';
import { connect } from 'react-redux';
import { Switch, Route, Redirect, Link } from 'react-router-dom';
import * as ReactGA from 'react-ga';
import v4 from 'uuid/v4';

import { Layout, Breadcrumb, Popover, Typography, Spin } from 'antd';

import DashboardNav from './DashboardNav';
import ConnectedUserProfile from './UserProfile';
import ConnectedGrades from './Grades';
import ConnectedGradeBreakdown from './Grades/GradeBreakdown';
import ConnectedUpgrades from './Upgrades';
import ConnectedLogout from './Logout';
import UpdateHandler from './UpdateHandler';
import env from '../../util/env';
import { getUser } from '../../actions/canvas';
import ConnectedErrorModal from './ErrorModal';

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
  const [hasSentUserToGa, setHasSentUserToGa] = useState(false);
  const [getUserId, setGetUserId] = useState();

  // if no token exists, redirect
  if (!localStorage.token) {
    return <Redirect to="/" />;
  } else if (!token) {
    // otherwise, wait for token
    return null;
  }

  ReactGA.pageview(
    props.location.pathname +
      (props.location.search.includes('~') ? '' : props.location.search)
  );

  const { location, user, subdomain, loading, error, dispatch } = props;
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

  if (!user && !getUserId) {
    const id = v4();
    dispatch(getUser(id, token, subdomain));
    setGetUserId(id);
  }

  if (error[getUserId]) {
    return <ConnectedErrorModal error={error[getUserId]} />;
  }

  return (
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
          {!loading.includes(getUserId) ? (
            <Switch>
              <Route
                exact
                path="/dashboard"
                render={() => <Redirect to="/dashboard/grades" />}
              />
              <Route
                exact
                path="/dashboard/profile"
                component={ConnectedUserProfile}
              />
              <Route
                exact
                path="/dashboard/grades"
                component={ConnectedGrades}
              />
              <Route
                exact
                path="/dashboard/grades/:courseId"
                component={ConnectedGradeBreakdown}
              />
              <Route
                exact
                path="/dashboard/upgrades"
                component={ConnectedUpgrades}
              />
              <Route
                exact
                path="/dashboard/logout"
                component={ConnectedLogout}
              />
              <Route render={() => <Redirect to="/dashboard" />} />
            </Switch>
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
  );
}

const ConnectedDashboard = connect(state => ({
  token: state.canvas.token,
  subdomain: state.canvas.subdomain,
  courses: state.canvas.courses,
  user: state.canvas.user,
  loading: state.loading,
  error: state.error
}))(Dashboard);

export default ConnectedDashboard;
