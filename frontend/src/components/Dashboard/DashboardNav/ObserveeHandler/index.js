import React from 'react';
import { connect } from 'react-redux';
import { isMobile } from 'react-device-detect';

import { Icon, Menu, Dropdown, Button } from 'antd';
import { Accordion as MobileAccordion, List as MobileList } from 'antd-mobile';
import { changeActiveUser } from '../../../../actions/canvas';

function Index(props) {
  const { mobileToggleMenu, observees, activeUserId, users, dispatch } = props;

  if (!users || !activeUserId || !observees || observees.length < 1) {
    return null;
  }

  if (isMobile) {
    return (
      <MobileAccordion defaultActiveKey="observees">
        <MobileAccordion.Panel
          header={<b>{users[activeUserId].name}</b>}
          key="observees"
        >
          <MobileList>
            {observees.map(
              o =>
                o.id !== activeUserId && (
                  <MobileList.Item
                    type="link"
                    key={o.id}
                    onClick={() => {
                      dispatch(changeActiveUser(o.id));
                      mobileToggleMenu && mobileToggleMenu();
                    }}
                  >
                    {o.name}
                  </MobileList.Item>
                )
            )}
          </MobileList>
        </MobileAccordion.Panel>
      </MobileAccordion>
    );
  }

  const menu = (
    <Menu>
      {observees.map(
        o =>
          o.id !== activeUserId && (
            <Menu.Item
              key={o.id}
              onClick={() => dispatch(changeActiveUser(o.id))}
            >
              {o.name}
            </Menu.Item>
          )
      )}
    </Menu>
  );

  return (
    <Dropdown overlay={menu}>
      <Button type="link">
        {users[activeUserId].name} <Icon type="down" />
      </Button>
    </Dropdown>
  );
}

const ConnectedObserveeHandler = connect(state => ({
  observees: state.canvas.observees,
  activeUserId: state.canvas.activeUserId,
  users: state.canvas.users
}))(Index);

export default ConnectedObserveeHandler;
