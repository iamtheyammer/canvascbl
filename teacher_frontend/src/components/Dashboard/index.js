import React, { useEffect } from 'react';
import { Route, Switch } from 'react-router-dom';
import { connect } from 'react-redux';
import { Breadcrumb, Icon, Layout, Modal, Popover, Typography } from 'antd';
import DashboardNav from './DashboardNav';

import env from '../../util/env';
import * as loginReturnTo from '../../util/loginReturnTo';
import getUrlPrefix from '../../util/getUrlPrefix';
import {
  destinationNames,
  destinationTypes,
  pageNameFromPath,
  trackDashboardLoad,
  trackExternalLinkClickOther,
  TrackingLink,
  vias
} from '../../util/tracking';
import UpdateHandler from './UpdateHandler';
// import Courses from './Courses';
import { getCourses, getUserProfile } from '../../actions/canvas';
// import CourseOverview from './Courses/CourseOverview';
import PopoutLink from '../PopoutLink';
import ComingSoon from './ComingSoon';

const { Content, Footer } = Layout;

const getBreadcrumbNameMap = (courses = []) => {
  const routes = {
    '/dashboard': 'Dashboard',
    '/dashboard/courses': 'Courses'
  };

  // courses.forEach(
  //   c => (routes[`/dashboard/grades/${c.id}`] = `Grade Breakdown for ${c.name}`)
  // );

  courses.forEach((c) => {
    routes[
      `/dashboard/courses/${c.original_course_id}_${c.distance_learning_course_id}`
    ] = c.course_name;
    routes[
      `/dashboard/courses/${c.original_course_id}_${c.distance_learning_course_id}/overview`
    ] = 'Overview';
  });

  return routes;
};

const routes = (
  <Switch>
    <Route path="/dashboard" component={ComingSoon} />
    {/*<Route*/}
    {/*  exact*/}
    {/*  path="/dashboard"*/}
    {/*  render={() => <Redirect to="/dashboard/courses" />}*/}
    {/*/>*/}
    {/*<Route exact path="/dashboard/courses" component={Courses} />*/}
    {/*<Route*/}
    {/*  exact*/}
    {/*  path="/dashboard/courses/:courseId"*/}
    {/*  render={() => <Redirect to={'/dashboard/courses'} />}*/}
    {/*/>*/}
    {/*<Route*/}
    {/*  exact*/}
    {/*  path="/dashboard/courses/:courseId/overview"*/}
    {/*  component={CourseOverview}*/}
    {/*/>*/}
  </Switch>
);

function Dashboard(props) {
  const {
    loggedOut,
    loadingUserProfile,
    getUserProfileError,
    userProfile,
    loadingCourses,
    getCoursesError,
    courses,
    distanceLearningPairs,
    dispatch,
    location
  } = props;

  useEffect(() => {
    if (
      !loggedOut &&
      !loadingUserProfile &&
      !getUserProfileError &&
      !userProfile
    ) {
      dispatch(getUserProfile());
    }

    if (userProfile) {
      trackDashboardLoad(
        userProfile.name,
        userProfile.email,
        userProfile.id,
        userProfile.canvas_user_id,
        env.currentVersion,
        localStorage.prevVersion
      );
    }
  }, [
    loggedOut,
    loadingUserProfile,
    getUserProfileError,
    userProfile,
    dispatch
  ]);

  useEffect(() => {
    if (!loggedOut && !loadingCourses && !getCoursesError && !courses) {
      dispatch(getCourses());
    }
  }, [loggedOut, loadingCourses, getCoursesError, courses, dispatch]);

  const pathSnippets = location.pathname.split('/').filter((i) => i);
  // const breadcrumbNameMap = getBreadcrumbNameMap(courses || []);
  const breadcrumbNameMap = getBreadcrumbNameMap(distanceLearningPairs || []);
  const breadcrumbItems = pathSnippets.map((_, index) => {
    const url = `/${pathSnippets.slice(0, index + 1).join('/')}`;
    const pageName = pageNameFromPath(url);
    return (
      <Breadcrumb.Item key={url}>
        {!breadcrumbNameMap[url] && '...'}
        <TrackingLink to={url} pageName={pageName} via={vias.breadcrumb}>
          {breadcrumbNameMap[url]}
        </TrackingLink>
      </Breadcrumb.Item>
    );
  });

  // handle non-teachers
  useEffect(() => {
    if (courses) {
      let foundTeacher = false;
      let isStudent = false;
      let isOther = false;
      courses.forEach((c) =>
        c.enrollments.forEach((e) => {
          if (e.type === 'teacher' || e.role === 'TeacherEnrollment') {
            foundTeacher = true;
          } else if (e.type === 'student' || e.role === 'StudentEnrollment') {
            isStudent = true;
          } else {
            isOther = true;
          }
        })
      );

      if (!foundTeacher) {
        if (isStudent && !isOther) {
          Modal.error({
            title: 'My apologies, young grasshopper.',
            content: "Students aren't able to use CanvasCBL for Teachers.",
            onOk: () => {
              trackExternalLinkClickOther(
                env.canvascblUrl,
                destinationTypes.canvascbl,
                destinationNames.canvascblForStudentsAndParents,
                vias.notATeacherPopup
              );
              window.location = env.canvascblUrl;
            },
            okText: (
              <>
                Go to CanvasCBL for Students <Icon type="arrow-right" />
              </>
            ),
            maskClosable: false
          });
        } else {
          Modal.error({
            title: "You don't appear to be a teacher",
            content: 'Only teachers are able to use CanvasCBL for Teachers.',
            onOk: () => {
              trackExternalLinkClickOther(
                env.canvascblUrl,
                destinationTypes.canvascbl,
                destinationNames.canvascblForStudentsAndParents,
                vias.notATeacherPopup
              );
              window.location = env.canvascblUrl;
            },
            okText: (
              <>
                Go to CanvasCBL for Students and Parents{' '}
                <Icon type="arrow-right" />
              </>
            ),
            maskClosable: false
          });
        }
      }
    }
  }, [courses]);

  const err = getUserProfileError;
  if (err) {
    loginReturnTo.set(location);
    if (err.error) {
      if (err.action === 'redirect_to_oauth') {
        loginReturnTo.set(location);
        window.location.href = `${getUrlPrefix}/api/canvas/oauth2/request?dest=teacher`;
        return null;
      } else {
        window.location = env.accountUrl + '?dest=teacher';
      }
    } else {
      return (
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

  if (loggedOut) {
    window.location = env.canvascblUrl;
    return null;
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
            {routes}
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
            CanvasCBL for Teachers Beta · Built by Sam Mendelson{' '}
            {new Date().getFullYear()}
          </Popover>
        </Footer>
      </Layout>
    </div>
  );
}

const ConnectedDashboard = connect((state) => ({
  loggedOut: state.canvas.loggedOut,
  loadingUserProfile: state.canvas.loadingUserProfile,
  getUserProfileError: state.canvas.getUserProfileError,
  userProfile: state.canvas.userProfile,
  loadingCourses: state.canvas.loadingCourses,
  getCoursesError: state.canvas.getCoursesError,
  courses: state.canvas.courses,
  distanceLearningPairs: state.canvas.distanceLearningPairs
}))(Dashboard);

export default ConnectedDashboard;
