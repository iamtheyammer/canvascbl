import React from 'react';
import { connect } from 'react-redux';
import { Layout, Menu, Typography } from 'antd';
import styled from 'styled-components';
import logo from '../../../assets/banner-light.svg';
import {
  destinationNames,
  destinationTypes,
  TrackingLink,
  trackLogout,
  vias,
} from '../../../util/tracking';
import PopoutLink from '../../PopoutLink';
import { logout } from '../../../actions/canvas';

const { Header } = Layout;

const StyledLogo = styled.img`
  width: 126px;
  height: auto;
  margin: 13px 18px 0 0;
  float: left;
`;

const StyledMenu = styled(Menu)`
  line-height: 64px;
  float: left;
`;

const StyledRightMenu = styled(Menu)`
  line-height: 64px;
  float: right;
`;

function DashboardNav(props) {
  const { loggedOut, loadingLogout, logoutError, dispatch } = props;

  return (
    <Header>
      <StyledLogo src={logo} alt="CanvasCBL Logo" />
      <StyledMenu
        theme="dark"
        mode="horizontal"
        defaultSelectedKeys="/dashboard/courses"
      >
        <Menu.Item key="/dashboard/courses">
          <TrackingLink to="/dashboard/courses" via={vias.dashboardMenu}>
            Courses
          </TrackingLink>
        </Menu.Item>
      </StyledMenu>
      <StyledRightMenu theme="dark" mode="horizontal" selectedKeys={[]}>
        <Menu.Item key="/dashboard/feedback">
          <PopoutLink
            url="https://go.canvascbl.com/teacher-feedback"
            tracking={{
              destinationName: destinationNames.googleForms,
              destinationType:
                destinationTypes.canvascblForTeachersFeedbackForm,
              via: vias.dashboardMenu,
            }}
            addIcon
          >
            Provide Feedback
          </PopoutLink>
        </Menu.Item>
        <Menu.Item
          key="/dashboard/logout"
          onClick={() => {
            trackLogout(vias.dashboardMenu);
            dispatch(logout());
          }}
        >
          {!loggedOut && !loadingLogout && !logoutError && 'Logout'}
          {loggedOut && 'Logged out'}
          {loadingLogout && 'Loading...'}
          {logoutError && (
            <Typography.Text type="danger">
              Error logging you out
            </Typography.Text>
          )}
        </Menu.Item>
      </StyledRightMenu>
    </Header>
  );
}

const ConnectedDashboardNav = connect((state) => ({
  loggedOut: state.canvas.loggedOut,
  loadingLogout: state.canvas.loadingLogout,
  logoutError: state.canvas.logoutError,
}))(DashboardNav);

export default ConnectedDashboardNav;
