import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { Switch, Route, Redirect } from 'react-router-dom';
import * as ReactGA from 'react-ga';
import v4 from 'uuid/v4';
import { isMobile } from 'react-device-detect';
import * as loginReturnTo from '../../util/loginReturnTo';

import {
  Layout,
  Breadcrumb,
  Popover,
  Typography,
  notification,
  Modal
} from 'antd';

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
import OAuth2Consent from './OAuth2Consent';
import Settings from './Settings';
import {
  pageNameFromPath,
  trackDashboardLoad,
  TrackingLink,
  vias
} from '../../util/tracking';
import PopoutLink from '../PopoutLink';

const { Content, Footer } = Layout;

const getBreadcrumbNameMap = (courses = []) => {
  const routes = {
    '/dashboard': 'Dashboard',
    '/dashboard/profile': 'Profile',
    '/dashboard/grades': 'Grades',
    '/dashboard/upgrades': 'Upgrades',
    '/dashboard/upgrades/redeem': 'Redeem',
    '/dashboard/authorize': 'Authorize an App',
    '/dashboard/settings': 'Settings'
  };

  courses.forEach(
    c => (routes[`/dashboard/grades/${c.id}`] = `Grade Breakdown for ${c.name}`)
  );

  return routes;
};

function Dashboard(props) {
  const [hasSentUserToGa, setHasSentUserToGa] = useState(false);
  const [hasInitializedTracking, setHasInitializedTracking] = useState(false);
  const [getInitialDataId, setGetInitialDataId] = useState();

  const {
    location,
    user,
    courses,
    activeUserId,
    session,
    loading,
    error,
    dispatch
  } = props;

  useEffect(() => {
    ReactGA.pageview(props.location.pathname + props.location.search);
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

  useEffect(() => {
    if (hasInitializedTracking || !user || !session || !activeUserId) {
      return;
    }

    trackDashboardLoad(
      user.name,
      session.email,
      session.has_valid_subscription,
      session.subscription_status,
      session.user_id,
      user.id,
      activeUserId,
      env.currentVersion,
      localStorage.prevVersion
    );

    setHasInitializedTracking(true);
  }, [session, user, activeUserId, hasInitializedTracking]);

  useEffect(() => {
    if (courses) {
      let modalShown = false;
      courses.forEach(c =>
        c.enrollments.map(e => {
          if (!modalShown && e.type === 'teacher') {
            Modal.confirm({
              title: 'Are you in the right place?',
              content:
                "You're currently at CanvasCBL for Students and Parents. Do you want to go to CanvasCBL for Teachers?",
              okText: 'Go to CanvasCBL for Teachers',
              cancelText: 'No, thanks',
              onOk: () => (window.location = env.teacherUrl)
            });
            modalShown = true;
          }
          return null;
        })
      );
    }
  });

  const pathSnippets = location.pathname.split('/').filter(i => i);
  const breadcrumbNameMap = getBreadcrumbNameMap(props.courses || []);
  const breadcrumbItems = pathSnippets.map((_, index) => {
    const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;
    return (
      <Breadcrumb.Item key={url}>
        <TrackingLink
          to={url}
          pageName={pageNameFromPath(url)}
          via={vias.breadcrumb}
        >
          {breadcrumbNameMap[url]}
        </TrackingLink>
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

  let errData;

  const err = error[getInitialDataId];
  if (err) {
    if (err.res) {
      const data = err.res.data;
      switch (data.action) {
        case 'redirect_to_oauth':
          // we'll be reauthing
          loginReturnTo.set(location);
          window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request`;
          return null;
        case 'retry':
          const id = v4();
          dispatch(getInitialData(id));
          setGetInitialDataId(id);
          return null;
        default:
          loginReturnTo.set(location);
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
    } else {
      errData = (
        <Typography.Text type="danger">
          We seem to have encountered a bit of an unexpected error. If this
          keeps happening, please{' '}
          <PopoutLink url="https://go.canvascbl.com/help/contact" addIcon>
            contact us.
          </PopoutLink>
        </Typography.Text>
      );
    }
  }

  if (session && session.status === 1) {
    return <Redirect to={'/dashboard/logout'} />;
  }

  const ready = user && !loading.includes(getInitialDataId);

  const lrt = loginReturnTo.get();
  if (lrt && ready) {
    loginReturnTo.clear();
    return <Redirect to={lrt} />;
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
      <Route exact path="/dashboard/settings" component={Settings} />
      <Route exact path="/dashboard/authorize" component={OAuth2Consent} />
      <Route exact path="/dashboard/logout" component={ConnectedLogout} />
      <Route render={() => <Redirect to="/" />} />
    </Switch>
  );

  if (isMobile) {
    function displayContent() {
      if (ready) {
        return routes;
      } else if (errData) {
        return errData;
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
            {!ready && !errData && <Loading text="CanvasCBL" />}
            {ready && routes}
            {errData && errData}
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
            Built by Sam Mendelson {new Date().getFullYear()}
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
  activeUserId: state.canvas.activeUserId,
  session: state.plus.session,
  loading: state.loading,
  error: state.error
}))(Dashboard);

export default ConnectedDashboard;
