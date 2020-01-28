import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { Switch, Route, Redirect, Link } from 'react-router-dom';
import * as ReactGA from 'react-ga';
import v4 from 'uuid/v4';
import { isMobile } from 'react-device-detect';

import { Layout, Breadcrumb, Popover, Typography, notification } from 'antd';

import DashboardNav from './DashboardNav';
import ConnectedUserProfile from './UserProfile';
import ConnectedGrades from './Grades';
import ConnectedGradeBreakdown from './Grades/GradeBreakdown';
import ConnectedUpgrades from './Upgrades';
import ConnectedLogout from './Logout';
import UpdateHandler from './UpdateHandler';
import env from '../../util/env';
import { getInitialData } from '../../actions/canvas';
import ConnectedErrorModal from './ErrorModal';
import './index.css';
import ConnectedRedeem from './Upgrades/Redeem';
import Loading from './Loading';
import getUrlPrefix from '../../util/getUrlPrefix';

const { Content, Footer } = Layout;

const getBreadcrumbNameMap = (courses = []) => {
  const routes = {
    '/dashboard': 'Dashboard',
    '/dashboard/profile': 'Profile',
    '/dashboard/grades': 'Grades',
    '/dashboard/upgrades': 'Upgrades',
    '/dashboard/upgrades/redeem': 'Redeem'
  };

  courses.forEach(
    c => (routes[`/dashboard/grades/${c.id}`] = `Grade Breakdown for ${c.name}`)
  );

  return routes;
};

function Dashboard(props) {
  const [hasSentUserToGa, setHasSentUserToGa] = useState(false);
  const [getInitialDataId, setGetInitialDataId] = useState();

  const { location, user, session, loading, error, dispatch } = props;

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

  if (!user && !getInitialDataId) {
    const id = v4();
    dispatch(getInitialData(id));
    setGetInitialDataId(id);
  }

  const err = error[getInitialDataId];
  if (err) {
    const data = err.res.data;
    switch (data.action) {
      case 'redirect_to_oauth':
        window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request`;
        return null;
      case 'retry':
        const id = v4();
        dispatch(getInitialData(id));
        setGetInitialDataId(id);
        return null;
      default:
        if (data.error.includes('no session string')) {
          return <Redirect to={'/'} />;
        } else if (data.error === 'expired session') {
          window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request?intent=reauth`;
          return null;
        } else if (data.error === 'invalid session string') {
          return <Redirect to="/" />;
        }
        return <ConnectedErrorModal error={err} />;
    }
  }

  if (session && session.status === 1) {
    return <Redirect to={'/dashboard/logout'} />;
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
      <Route
        exact
        path="/dashboard/upgrades/redeem"
        component={ConnectedRedeem}
      />
      <Route exact path="/dashboard/logout" component={ConnectedLogout} />
      <Route render={() => <Redirect to="/" />} />
    </Switch>
  );

  if (isMobile) {
    function displayContent() {
      if (user && !loading.includes(getInitialDataId)) {
        return routes;
      } else {
        return <Loading text="CanvasCBL" />;
      }
    }

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
          {displayContent()}
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
            {user && !loading.includes(getInitialDataId) ? (
              routes
            ) : (
              <Loading text="CanvasCBL" />
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
  observees: state.canvas.observees,
  session: state.plus.session,
  loading: state.loading,
  error: state.error
}))(Dashboard);

export default ConnectedDashboard;
