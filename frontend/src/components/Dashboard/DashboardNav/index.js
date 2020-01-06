import React, { useState } from 'react';
import { Link, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import './index.css';
import logo from '../../../assets/banner-light.svg';
import logoPlus from '../../../assets/banner-light-plus.svg';
import env from '../../../util/env';
import { isMobile } from 'react-device-detect';

import { Layout, Menu, Typography } from 'antd';
import {
  NavBar as MobileNavBar,
  Drawer as MobileDrawer,
  List as MobileList,
  Icon as MobileIcon
} from 'antd-mobile';
import PopoutLink from '../../PopoutLink';
import ObserveeHandler from './ObserveeHandler';

const { Header } = Layout;

function DashboardNav(props) {
  const { session, observees } = props;
  const userHasActiveSubscription = session && session.hasValidSubscription;

  const [shouldShowMobileMenu, setShouldShowMobileMenu] = useState(false);

  if (isMobile) {
    const toggleMenu = () => setShouldShowMobileMenu(!shouldShowMobileMenu);

    return (
      <div>
        <MobileNavBar
          leftContent={<MobileIcon type="ellipsis" />}
          mode="dark"
          onLeftClick={toggleMenu}
        >
          <Link to="/dashboard" style={{ color: 'white' }}>
            CanvasCBL{userHasActiveSubscription && '+'}
          </Link>
        </MobileNavBar>
        <MobileDrawer
          className="mobile-drawer"
          style={{ minHeight: document.documentElement.clientHeight }}
          contentStyle={{ paddingTop: 12, paddingLeft: 8, paddingRight: 4 }}
          sidebar={
            <MobileList>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/profile">
                <Link to="/dashboard/profile">Profile</Link>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/grades">
                <Link to="/dashboard/grades">Grades</Link>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/upgrades">
                <Link to="/dashboard/upgrades">Upgrades</Link>
              </MobileList.Item>
              <MobileList.Item />
              {observees && observees.length > 1 && (
                <div>
                  <ObserveeHandler mobileToggleMenu={toggleMenu} />
                  <MobileList.Item />
                </div>
              )}
              <MobileList.Item onClick={toggleMenu} key="contactSupport">
                <PopoutLink
                  url={
                    'mailto:sam@canvascbl.com?subject=CanvasCBL%20Question%20or%20Comment'
                  }
                >
                  Contact Support
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="viewPrivacyPolicy">
                <PopoutLink url={env.privacyPolicyUrl}>
                  View Privacy Policy
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="viewTermsOfService">
                <PopoutLink url={env.termsOfServiceUrl}>
                  View Terms of Service
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="viewSystemStatus">
                <PopoutLink url="https://status.canvascbl.com">
                  View System Status
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/logout">
                <Link to="/dashboard/logout">Logout</Link>
              </MobileList.Item>
            </MobileList>
          }
          open={shouldShowMobileMenu}
          onOpenChange={toggleMenu}
        >
          {props.children}
        </MobileDrawer>
      </div>
    );
  }

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
        {env.buildBranch !== 'master' && (
          <Menu.Item key="nonProductionMode">
            <Typography.Text type="danger">
              CanvasCBL is running in {env.buildBranch}
            </Typography.Text>
          </Menu.Item>
        )}
        {observees && observees.length > 0 && <ObserveeHandler />}
        <Menu.SubMenu key="moreActions" title="More Actions">
          <Menu.Item key="installExtension">
            <PopoutLink url="https://chrome.google.com/webstore/detail/canvascbl-add-in-for-canv/odmbdioejfbelhcknliaihbjckggmmak">
              Install the Extension
            </PopoutLink>
          </Menu.Item>
          <Menu.Item key="contactSupport">
            <PopoutLink
              url={
                'mailto:sam@canvascbl.com?subject=CanvasCBL%20Question%20or%20Comment'
              }
            >
              Contact Support
            </PopoutLink>
          </Menu.Item>
          <Menu.Item key="viewSystemStatus">
            <PopoutLink url="https://status.canvascbl.com">
              View System Status
            </PopoutLink>
          </Menu.Item>
          <Menu.Item key="viewPrivacyPolicy">
            <PopoutLink url={env.privacyPolicyUrl}>
              View Privacy Policy
            </PopoutLink>
          </Menu.Item>
          <Menu.Item key="viewTermsOfService">
            <PopoutLink url={env.termsOfServiceUrl}>
              View Terms of Service
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
  session: state.plus.session,
  observees: state.canvas.observees
}))(DashboardNav);

export default withRouter(ConnectedDashboardNav);
