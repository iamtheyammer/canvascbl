import React, { useState } from 'react';
import { withRouter } from 'react-router-dom';
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
import {
  destinationNames,
  destinationTypes,
  pageNames,
  TrackingLink,
  vias
} from '../../../util/tracking';

const { Header } = Layout;

function DashboardNav(props) {
  const { session, observees } = props;
  const userHasActiveSubscription = session && session.has_valid_subscription;

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
          <TrackingLink
            to="/dashboard/grades"
            pageName={pageNames.grades}
            via={vias.mobileNavBarLogo}
            style={{ color: 'white' }}
          >
            CanvasCBL{userHasActiveSubscription && '+'}
          </TrackingLink>
        </MobileNavBar>
        <MobileDrawer
          className="mobile-drawer"
          style={{ minHeight: document.documentElement.clientHeight }}
          contentStyle={{ paddingTop: 12, paddingLeft: 8, paddingRight: 4 }}
          sidebar={
            <MobileList>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/profile">
                <TrackingLink
                  to="/dashboard/profile"
                  pageName={pageNames.profile}
                  via={vias.dashboardMenu}
                >
                  Profile
                </TrackingLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/grades">
                <TrackingLink
                  to="/dashboard/grades"
                  pageName={pageNames.grades}
                  via={vias.dashboardMenu}
                >
                  Grades
                </TrackingLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/upgrades">
                <TrackingLink
                  to="/dashboard/upgrades"
                  pageName={pageNames.upgrades}
                  via={vias.dashboardMenu}
                >
                  Upgrades
                </TrackingLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/settings">
                <TrackingLink
                  to="/dashboard/settings"
                  pageName={pageNames.settings}
                  via={vias.dashboardMenu}
                >
                  Settings
                </TrackingLink>
              </MobileList.Item>
              <MobileList.Item />
              {observees && observees.length > 1 && (
                <div>
                  <ObserveeHandler
                    mobileToggleMenu={toggleMenu}
                    via={vias.dashboardMenu}
                  />
                  <MobileList.Item />
                </div>
              )}
              <MobileList.Item onClick={toggleMenu} key="helpAndSupport">
                <PopoutLink
                  url="https://go.canvascbl.com/help"
                  tracking={{
                    destinationName: destinationNames.helpdesk,
                    destinationType: destinationTypes.helpdesk.home,
                    via: vias.dashboardMenu
                  }}
                  id="helpAndSupport"
                >
                  Help & Support
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="viewPrivacyPolicy">
                <PopoutLink
                  url={env.privacyPolicyUrl}
                  tracking={{
                    destinationName: destinationNames.privacyPolicy,
                    via: vias.dashboardMenu
                  }}
                  id="viewPrivacyPolicy"
                >
                  View Privacy Policy
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="viewTermsOfService">
                <PopoutLink
                  url={env.termsOfServiceUrl}
                  tracking={{
                    destinationName: destinationNames.termsOfService,
                    via: vias.dashboardMenu
                  }}
                  id="viewTermsOfService"
                >
                  View Terms of Service
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="viewSystemStatus">
                <PopoutLink
                  url="https://status.canvascbl.com"
                  tracking={{
                    destinationName: destinationNames.statusPage,
                    via: vias.dashboardMenu
                  }}
                  id="viewSystemStatus"
                >
                  View System Status
                </PopoutLink>
              </MobileList.Item>
              <MobileList.Item onClick={toggleMenu} key="/dashboard/logout">
                <TrackingLink
                  to="/dashboard/logout"
                  pageName={pageNames.logout}
                  via={vias.dashboardMenu}
                >
                  Logout
                </TrackingLink>
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
          `/${props.location.pathname.split('/').slice(1, 3).join('/')}`
        ]}
        style={{ lineHeight: '64px', float: 'left' }}
      >
        <Menu.Item key="/dashboard/profile">
          <TrackingLink
            to="/dashboard/profile"
            pageName={pageNames.profile}
            via={vias.dashboardMenu}
          >
            Profile
          </TrackingLink>
        </Menu.Item>
        <Menu.Item key="/dashboard/grades">
          <TrackingLink
            to="/dashboard/grades"
            pageName={pageNames.grades}
            via={vias.dashboardMenu}
          >
            Grades
          </TrackingLink>
        </Menu.Item>
        <Menu.Item key="/dashboard/upgrades">
          <TrackingLink
            to="/dashboard/upgrades"
            pageName={pageNames.upgrades}
            via={vias.dashboardMenu}
          >
            Upgrades
          </TrackingLink>
        </Menu.Item>
        <Menu.Item key="/dashboard/settings">
          <TrackingLink
            to="/dashboard/settings"
            pageName={pageNames.settings}
            via={vias.dashboardMenu}
          >
            Settings
          </TrackingLink>
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
        {observees && observees.length > 0 && (
          <ObserveeHandler via={vias.dashboardMenu} />
        )}
        <Menu.SubMenu key="moreActions" title="More Actions">
          <Menu.Item key="installExtension">
            <PopoutLink
              url="https://go.canvascbl.com/extension"
              tracking={{
                destinationName: destinationNames.extension,
                via: vias.moreActionsSubmenu
              }}
            >
              Install the Extension
            </PopoutLink>
          </Menu.Item>
          <Menu.Item key="helpAndSupport">
            <PopoutLink
              url={'https://go.canvascbl.com/help'}
              tracking={{
                destinationType: destinationTypes.helpdesk.home,
                destinationName: destinationNames.helpdesk,
                via: vias.moreActionsSubmenu
              }}
            >
              Help & Support
            </PopoutLink>
          </Menu.Item>
          <Menu.Item key="viewSystemStatus">
            <PopoutLink
              url="https://go.canvascbl.com/status"
              tracking={{
                destinationName: destinationNames.statusPage,
                via: vias.moreActionsSubmenu
              }}
            >
              View System Status
            </PopoutLink>
          </Menu.Item>
          <Menu.Item key="viewPrivacyPolicy">
            <PopoutLink
              url={env.privacyPolicyUrl}
              tracking={{
                destinationName: destinationNames.privacyPolicy,
                via: vias.moreActionsSubmenu
              }}
            >
              View Privacy Policy
            </PopoutLink>
          </Menu.Item>
          <Menu.Item key="viewTermsOfService">
            <PopoutLink
              url={env.termsOfServiceUrl}
              tracking={{
                destinationName: destinationNames.termsOfService,
                via: vias.moreActionsSubmenu
              }}
            >
              View Terms of Service
            </PopoutLink>
          </Menu.Item>
        </Menu.SubMenu>
        <Menu.Item key="/dashboard/logout">
          <TrackingLink
            to="/dashboard/logout"
            pageName={pageNames.logout}
            via={vias.dashboardMenu}
          >
            Logout
          </TrackingLink>
        </Menu.Item>
      </Menu>
    </Header>
  );
}

const ConnectedDashboardNav = connect((state) => ({
  session: state.plus.session,
  observees: state.canvas.observees
}))(DashboardNav);

export default withRouter(ConnectedDashboardNav);
