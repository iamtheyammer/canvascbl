import React from 'react';
import { Link, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import './index.css';
import logo from '../../../assets/banner-light.svg';
import logoPlus from '../../../assets/banner-light-plus.svg';

import { Layout, Menu } from 'antd';
import PopoutLink from '../../PopoutLink';

const { Header } = Layout;

function DashboardNav(props) {
  const { session } = props;
  const userHasActiveSubscription = session && session.hasValidSubscription;

  return (
    <Header>
      <img
        src={userHasActiveSubscription ? logoPlus : logo}
        className={userHasActiveSubscription ? 'logo-plus' : 'logo'}
        alt="canvas-grade-calculator light banner"
      />
      <Menu
        theme="dark"
        mode="horizontal"
        defaultSelectedKeys={['profile']}
        selectedKeys={[
          `/${props.location.pathname
            .split('/')
            .slice(1, 3)
            .join('/')}`
        ]}
        style={{ lineHeight: '64px', float: 'left' }}
      >
        <Menu.Item key="/dashboard/profile">
          <Link to="/dashboard/profile">Profile</Link>
        </Menu.Item>
        <Menu.Item key="/dashboard/grades">
          <Link to="/dashboard/grades">Grades</Link>
        </Menu.Item>
        <Menu.Item key="/dashboard/upgrades">
          <Link to="/dashboard/upgrades">Upgrades</Link>
        </Menu.Item>
      </Menu>
      <Menu
        theme="dark"
        mode="horizontal"
        style={{ lineHeight: '64px', float: 'right' }}
        selectable={false}
      >
        <Menu.SubMenu key="moreActions" title="More Actions">
          <Menu.Item key="contactSupport">
            <PopoutLink
              url={
                'mailto:sam@canvascbl.com?subject=CanvasCBL%20Question%20or%20Comment'
              }
            >
              Contact Support
            </PopoutLink>
          </Menu.Item>
        </Menu.SubMenu>
        <Menu.Item key="/dashboard/logout">
          <Link to={'/dashboard/logout'}>Logout</Link>
        </Menu.Item>
      </Menu>
    </Header>
  );
}

const ConnectedDashboardNav = connect(state => ({
  session: state.plus.session
}))(DashboardNav);

export default withRouter(ConnectedDashboardNav);
